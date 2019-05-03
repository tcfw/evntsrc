package storer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/spf13/viper"

	nats "github.com/nats-io/go-nats"
	"github.com/tcfw/evntsrc/pkg/event"
	"github.com/tcfw/evntsrc/pkg/websocks"
	"github.com/tcfw/go-queue"
)

var pgdb *sql.DB

//Start inits required processes
func Start(nats string) {
	dbURL, ok := os.LookupEnv("PGDB_HOST")
	if !ok {
		log.Fatal("Failed to connect to DB")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	pgdb = db

	if err := createUpdateTables(pgdb); err != nil {
		log.Fatal(err)
	}

	if viper.GetBool("migrate") {
		os.Exit(0)
		return
	}

	StartMonitor(nats)
}

//StartMonitor subscripts to all user channels
func StartMonitor(nats string) {
	connectNats(nats)
	defer natsConn.Close()

	monitorUserStreams()
	monitorReplayRequests()

	//Wait forever
	select {}
}

type eventProcessor struct{}

func (ep *eventProcessor) Handle(job interface{}) {
	usrEvent := job.(*event.Event)

	if isReplay, ok := usrEvent.Metadata["replay"]; ok && isReplay == "true" {
		return
	}

	if isNonPersistent, ok := usrEvent.Metadata["non-persistent"]; ok && isNonPersistent == "true" {
		return
	}

	metadataJSON, err := json.Marshal(usrEvent.Metadata)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := pgdb.Begin()
	if err != nil {
		panic(err.Error)
	}

	if _, err := tx.Exec(
		`INSERT INTO event_store.events VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		usrEvent.ID,
		usrEvent.Stream,
		usrEvent.Time.Time,
		usrEvent.Type,
		usrEvent.TypeVersion,
		usrEvent.CEVersion,
		usrEvent.Source,
		usrEvent.Subject,
		usrEvent.Acknowledged.Time,
		string(metadataJSON),
		usrEvent.ContentType,
		usrEvent.Data,
	); err != nil {
		log.Fatal(err)
		panic(err.Error)
	}

	tx.Commit()

}

func monitorUserStreams() {
	log.Println("Watching for user streams...")

	dispatcher := queue.NewDispatcher(&eventProcessor{})
	if c := runtime.NumCPU(); c < 4 {
		dispatcher.MaxWorkers = 4
	}
	dispatcher.Run()

	natsConn.QueueSubscribe("_USER.>", "storers", func(m *nats.Msg) {
		event := &event.Event{}
		err := json.Unmarshal(m.Data, event)
		if err != nil {
			log.Printf("%s\n", err.Error())
			return
		}

		dispatcher.Queue(event)

	})
}

func monitorReplayRequests() {
	log.Println("Watching for replay requests...")
	natsConn.QueueSubscribe("replay.broadcast", "replayers", func(m *nats.Msg) {
		command := &websocks.ReplayCommand{}
		json.Unmarshal(m.Data, command)

		reply := make(chan []byte)
		errCh := make(chan error)

		go doReplay(command, reply, errCh)

		select {
		case msg := <-reply:
			natsConn.Publish(m.Reply, msg)
		case err := <-errCh:
			natsConn.Publish(m.Reply, []byte(fmt.Sprintf("failed: %s", err.Error())))
			log.Print(err.Error())
		}
	})
}

func doReplay(command *websocks.ReplayCommand, reply chan []byte, errCh chan error) {
	defer func(reply chan []byte, errCh chan error) {
		close(reply)
		close(errCh)
	}(reply, errCh)

	if command.JustMe && command.Dest == "" {
		errCh <- fmt.Errorf("no socket dest set")
		return
	}

	qSelector, params := buildBaseQuery(command)

	count, err := countEvents(qSelector, params)
	if err != nil {
		errCh <- err
		return
	}
	if count == 0 {
		reply <- []byte("no events")
		return
	}

	rD, err := pgdb.Query(`SELECT * `+qSelector, params...)
	if err != nil {
		errCh <- fmt.Errorf("sqld: %s", err.Error())
		return
	}
	defer func() {
		rD.Close()
	}()

	reply <- []byte("OK")

	for rD.Next() {
		event := &event.Event{
			Metadata: map[string]string{},
		}
		var mdString []byte

		err := rD.Scan(&event.ID,
			&event.Stream,
			&event.Time.Time,
			&event.Type,
			&event.TypeVersion,
			&event.CEVersion,
			&event.Source,
			&event.Subject,
			&event.Acknowledged.Time,
			&mdString,
			&event.ContentType,
			&event.Data,
		)
		if err != nil {
			errCh <- fmt.Errorf("sqld: %s", err.Error())
			return
		}

		if command.Query.EndID != "" && event.ID == command.Query.EndID {
			break
		}

		err = json.Unmarshal(mdString, &event.Metadata)
		if err != nil {
			errCh <- fmt.Errorf("sqld md: %s", err.Error())
			return
		}
		if event.Metadata == nil {
			event.Metadata = map[string]string{}
		}

		event.Metadata["replay"] = "true"

		jsonBytes, err := json.Marshal(event)
		if err != nil {
			errCh <- fmt.Errorf("sqld md: %s", err.Error())
			return
		}
		if command.JustMe {
			natsConn.Publish(fmt.Sprintf("_CONN.%s", command.Dest), jsonBytes)
		} else {
			natsConn.Publish(fmt.Sprintf("_USER.%d.%s", command.Stream, command.Subject), jsonBytes)
		}
	}
}

func buildBaseQuery(command *websocks.ReplayCommand) (string, []interface{}) {
	qSelector := `FROM event_store.events WHERE stream = $1`

	params := []interface{}{
		command.Stream,
	}

	if command.Query.StartTime != nil && command.Query.EndTime == nil {
		qSelector += ` AND time >= $2`
		params = append(params, command.Query.StartTime)
	} else if command.Query.StartTime != nil && command.Query.EndTime != nil {
		qSelector += ` AND time >= $2 AND time <= $3`
		params = append(params, command.Query.StartTime)
		params = append(params, command.Query.EndTime)
	}

	// qSelector += ` GROUP BY id, time ORDER BY time;`
	return qSelector, params
}

func countEvents(qSelector string, params []interface{}) (int, error) {
	count := 0

	rC, err := pgdb.Query(`SELECT COUNT(id) `+qSelector, params...)
	if err != nil {
		return count, fmt.Errorf("sqlc: %s", err.Error())
	}

	rC.Next()

	if err := rC.Scan(&count); err != nil {
		return count, fmt.Errorf("sqlc: %s", err.Error())
	}
	rC.Close()

	return count, nil
}

package storer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gogo/protobuf/proto"
	nats "github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	pbEvent "github.com/tcfw/evntsrc/internal/event/protos"
	"github.com/tcfw/evntsrc/internal/websocks"
)

func monitorReplayRequests() {
	log.Println("Watching for replay requests...")
	natsConn.QueueSubscribe("replay.broadcast", "replayers", func(m *nats.Msg) {
		timer := prometheus.NewTimer(replayDurations)
		defer timer.ObserveDuration()

		command := &websocks.ReplayCommand{}
		json.Unmarshal(m.Data, command)

		reply := make(chan []byte, 10)
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

	//Ack replay request
	reply <- []byte("OK")

	for rD.Next() {
		event, err := scanEvent(rD)
		if err != nil {
			errCh <- err
		}
		if command.Query.EndID != "" && event.ID == command.Query.EndID {
			break
		}

		bytes, err := proto.Marshal(event)
		if err != nil {
			errCh <- fmt.Errorf("failed replay proto marshal: %s", err.Error())
			return
		}
		if command.JustMe {
			dest := fmt.Sprintf("_CONN.%s", command.Dest)
			err = natsConn.Publish(dest, bytes)
			if err != nil {
				errCh <- fmt.Errorf("natspub: %s", err.Error())
				return
			}
		} else {
			dest := fmt.Sprintf("_USER.%d.%s", command.Stream, command.Subject)
			err = natsConn.Publish(dest, bytes)
			if err != nil {
				errCh <- fmt.Errorf("natspub: %s", err.Error())
				return
			}
		}
		replayEventCount.With(prometheus.Labels{"stream": fmt.Sprintf("%d", command.Stream)}).Inc()
	}
}

func scanEvent(rD *sql.Rows) (*pbEvent.Event, error) {
	event := &pbEvent.Event{
		Metadata: map[string]string{},
	}
	var mdString []byte

	err := rD.Scan(&event.ID,
		&event.Stream,
		&event.Time,
		&event.Type,
		&event.TypeVersion,
		&event.CEVersion,
		&event.Source,
		&event.Subject,
		&event.Acknowledged,
		&mdString,
		&event.ContentType,
		&event.Data,
	)
	if err != nil {
		return nil, fmt.Errorf("sqld: %s", err.Error())
	}

	err = json.Unmarshal(mdString, &event.Metadata)
	if err != nil {
		return nil, fmt.Errorf("sqld md: %s", err.Error())
	}
	if event.Metadata == nil {
		event.Metadata = map[string]string{}
	}

	event.Metadata["replay"] = "true"

	return event, nil
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

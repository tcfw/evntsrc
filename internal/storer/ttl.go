package storer

import (
	"database/sql"
	"fmt"
	"time"

	pb "github.com/tcfw/evntsrc/internal/storer/protos"
)

func (s *server) handleTTLQuery(req *pb.QueryRequest, stream pb.StorerService_QueryServer) error {

	ttlQuery, ok := req.Query.(*pb.QueryRequest_Ttl)
	if !ok {
		return fmt.Errorf("Not TTL query")
	}
	ttl := ttlQuery.Ttl
	uStream := req.Stream
	limit := req.Limit

	if limit == 0 {
		limit = 1000
	}
	if limit > 1000 {
		return fmt.Errorf("Limit too large")
	}

	query := `SELECT * FROM event_store.events WHERE stream = $1 AND metadata->>'ttl' <= $2 AND metadata->>'ttl' >= $3 AND acknowledged IS NULL ORDER BY time DESC LIMIT $4`

	maxTTL := time.Now().Add(-1 * 7 * 24 * time.Hour).Format(time.RFC3339)

	rD, err := pgdb.Query(query, uStream, ttl, maxTTL, limit)
	if err != nil {
		return err
	}

	for rD.Next() {
		event, err := scanEvent(rD)
		if err != nil {
			return err
		}

		if err = stream.Send(event); err != nil {
			return err
		}
	}

	return nil
}

func extendTTL(req *pb.ExtendTTLRequest) error {
	if req.CurrentTTL == nil {
		return fmt.Errorf("Current TTL required")
	}

	if req.TTLTime == nil {
		return fmt.Errorf("New TTL time required")
	}

	tx, err := pgdb.Begin()
	if err != nil {
		return err
	}

	query := `WHERE stream = $1 AND id = $2 LIMIT 1`

	rD, err := tx.Query(`SELECT * FROM event_store.events `+query, req.GetStream(), req.GetEventID())
	if err != nil {
		tx.Rollback()
		return err
	}

	rD.Next()
	event, err := scanEvent(rD)
	if err != nil {
		return err
	}
	rD.Close()

	eventMD := event.GetMetadata()

	if event.Acknowledged != nil && !event.Acknowledged.IsZero() {
		tx.Rollback()
		return fmt.Errorf("event already acknowledged")
	}

	ttl, ok := eventMD["ttl"]
	ttlTime, _ := time.Parse(time.RFC3339, ttl)

	fmt.Printf("%s", ttlTime.String())

	if ok && ttlTime != *req.CurrentTTL {
		tx.Rollback()
		return fmt.Errorf("incorrect matching TTL")
	}

	err = setTTL(tx, req.GetEventID(), *req.GetTTLTime())
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update TTL: %s", err)
	}

	return tx.Commit()
}

func setTTL(tx *sql.Tx, id string, ttl time.Time) error {
	_, err := tx.Exec(`
		UPDATE 
			event_store.events 
		SET 
			metadata = jsonb_set(IFNULL(to_jsonb(metadata), '{}'::jsonb), '{ttl}'::string[], to_jsonb($1::timestamp), TRUE) 
		WHERE 
			id = $2
		LIMIT 1`,
		ttl, id)
	return err
}

package storer

import (
	"fmt"
	"strconv"
	"time"

	pb "github.com/tcfw/evntsrc/internal/storer/protos"
)

const (
	//MaxTTLDiff The maximum lifetime of an event
	MaxTTLDiff = -1 * 168 * time.Hour //1 Week
)

func (s *server) handleTTLQuery(req *pb.QueryRequest, stream pb.StorerService_QueryServer) error {

	ttlQuery, ok := req.Query.(*pb.QueryRequest_Ttl)
	if !ok {
		return fmt.Errorf("Not TTL query")
	}
	ttl := ttlQuery.Ttl.GetTime()
	uStream := req.Stream
	limit := req.Limit

	if limit == 0 {
		limit = 1000
	}
	if limit > 1000 {
		return fmt.Errorf("Limit too large")
	}

	if ttl == nil {
		return fmt.Errorf("Last TTL Check time cannot be nil")
	}

	query := `
		SELECT 
			*
		FROM 
			event_store.events 
		WHERE 
			stream = $1 AND 
			metadata->>'ttl' <= $2 AND 
			(metadata->>'retries' <= 5 OR metadata->>'retries' IS NULL) AND 
			time >= $3 AND 
			acknowledged IS NULL 
		ORDER BY time DESC LIMIT $4`

	maxTTL := time.Now().Add(MaxTTLDiff).Format(time.RFC3339)

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
	if ok && ttlTime != *req.CurrentTTL {
		tx.Rollback()
		return fmt.Errorf("incorrect matching TTL")
	}

	err = setMD(tx, req.EventID, "ttl", req.GetTTLTime().Format(time.RFC3339))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update TTL: %s", err)
	}

	retries, ok := eventMD["retries"]
	if !ok {
		err = setMD(tx, req.EventID, "retries", strconv.Itoa(1))
	} else {
		retInt, _ := strconv.Atoi(retries)
		err = setMD(tx, req.EventID, "retries", strconv.Itoa(retInt+1))
	}
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update Retries: %s", err)
	}

	return tx.Commit()
}

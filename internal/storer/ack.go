package storer

import (
	"fmt"
	"time"
)

func ackEvent(stream int32, eventID string) (*time.Time, error) {
	tx, err := pgdb.Begin()
	if err != nil {
		return nil, err
	}

	query := `WHERE stream = $1 AND id = $2 LIMIT 1`

	rD, err := tx.Query(`SELECT * FROM event_store.events `+query, stream, eventID)
	if err != nil {
		tx.Rollback()
		return nil, nil
	}

	rD.Next()
	event, err := scanEvent(rD)
	if err != nil {
		return nil, err
	}
	rD.Close()

	if event.Acknowledged != nil && !event.Acknowledged.IsZero() {
		tx.Rollback()
		return nil, fmt.Errorf("event already acknowledged")
	}

	ackTime := time.Now()

	_, err = tx.Exec(`UPDATE event_store.events SET acknowledged = $3 `+query, stream, eventID, ackTime)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &ackTime, tx.Commit()
}

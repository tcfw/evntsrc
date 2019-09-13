package storer

import (
	"database/sql"
	"encoding/json"
	"fmt"

	pbEvent "github.com/tcfw/evntsrc/internal/event/protos"

	"github.com/prometheus/client_golang/prometheus"
)

func storeEvent(event *pbEvent.Event, db *sql.DB) error {
	if isReplay, ok := event.Metadata["replay"]; ok && isReplay == "true" {
		return fmt.Errorf("Cannot store replay event")
	}

	if _, ok := event.Metadata["forwarded"]; ok {
		return fmt.Errorf("Cannot store forwarded event")
	}

	if isNonPersistent(event) {
		return fmt.Errorf("Refusing to store non-persistent event")
	}

	metadataJSON, err := json.Marshal(event.Metadata)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	md := string(metadataJSON)
	if len(metadataJSON) == 0 || len(event.Metadata) == 0 {
		md = "{}"
	}

	if _, err := tx.Exec(
		`INSERT INTO event_store.events VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		event.ID,
		event.Stream,
		event.Time,
		event.Type,
		event.TypeVersion,
		event.CEVersion,
		event.Source,
		event.Subject,
		event.Acknowledged,
		md,
		event.ContentType,
		event.Data,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("ev store: %s", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ev commit: %s", err)
	}

	storeCount.With(prometheus.Labels{"stream": fmt.Sprintf("%d", event.Stream)}).Inc()
	return nil
}

func isNonPersistent(event *pbEvent.Event) bool {
	if isNonPersistent, ok := event.Metadata["non-persistent"]; ok && isNonPersistent == "true" {
		return true
	}

	if isNonPersistent, ok := event.Metadata["persistent"]; ok && isNonPersistent == "false" {
		return true
	}

	return false
}

func setMD(tx *sql.Tx, id string, mdField string, data string) error {
	_, err := tx.Exec(`
		UPDATE
			event_store.events 
		SET 
			metadata = jsonb_set(IFNULL(to_jsonb(metadata), '{}'::jsonb), $1::string[], to_jsonb($2::string), TRUE) 
		WHERE 
			id = $3
		LIMIT 1`,
		fmt.Sprintf("{%s}", mdField), data, id)
	return err
}

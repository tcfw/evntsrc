package storer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

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

	if isNonPersistent, ok := event.Metadata["non-persistent"]; ok && isNonPersistent == "true" {
		return fmt.Errorf("Refusing to store non-persistent event")
	}

	metadataJSON, err := json.Marshal(event.Metadata)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err.Error)
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
		log.Fatal(err)
		panic(err.Error)
	}

	tx.Commit()

	storeCount.With(prometheus.Labels{"stream": fmt.Sprintf("%d", event.Stream)}).Inc()
	return nil
}

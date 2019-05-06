package storer

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/spf13/viper"
)

type migration struct {
	file       string
	migratedAt time.Time
}

func createUpdateTables(db *sql.DB) error {
	statements := []string{
		`CREATE DATABASE IF NOT EXISTS event_store`,
		`CREATE TABLE IF NOT EXISTS event_store.migrations (
			file STRING,
			migrated_at TIMESTAMP,
			PRIMARY KEY (file)
		)`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	if !viper.GetBool("migrate") && flag.Lookup("test.v") == nil {
		return nil
	}

	log.Println("Starting migrations...")

	files, err := ioutil.ReadDir("./migrations")
	if err != nil {
		log.Fatal(err)
	}

	migratedRows, err := db.Query("select file from event_store.migrations;")
	defer migratedRows.Close()

	migrations := map[string]bool{}
	if migratedRows != nil {
		for migratedRows.Next() {
			var file string
			if err := migratedRows.Scan(&file); err != nil {
				// Check for a scan error.
				// Query rows will be closed with defer.
				log.Fatal(err)
			}
			migrations[file] = true
		}
	}

	for _, f := range files {
		if _, ok := migrations[f.Name()]; ok {
			continue
		}

		log.Printf("Running %s...\n", f.Name())

		migBytes, err := ioutil.ReadFile("./migrations/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}

		if _, err := db.Exec(string(migBytes)); err != nil {
			log.Fatal(err)
		}

		if _, err := db.Exec(`insert into event_store.migrations VALUES ($1, NOW())`, f.Name()); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

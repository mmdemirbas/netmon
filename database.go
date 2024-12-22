package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/showwin/speedtest-go/speedtest"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var dbConnection *sql.DB // Global database connection

type Metrics struct {
	Timestamp   time.Time
	NetworkName string
	Online      bool
	Data        *speedtest.Server
}

func initDatabase(absoluteDbFilePath string) error {
	// Create parent directories if they don't exist
	dataDirName := filepath.Dir(absoluteDbFilePath)
	err := os.MkdirAll(dataDirName, 0755)
	if err != nil {
		return fmt.Errorf("Error creating data directory: %v", err)
	}

	// Open the database
	dbConnection, err = sql.Open("sqlite3", absoluteDbFilePath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Create the table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS metrics (
		timestamp INTEGER,
		network_name TEXT,
		online INTEGER,
		metrics TEXT
	);
	`
	_, err = dbConnection.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

func saveMetric(metrics *Metrics) error {

	// Marshal the metrics to JSON
	metricsJson, err := json.Marshal(metrics.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %v", err)
	}

	// Begin a transaction
	tx, err := dbConnection.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback() // Rollback the transaction if it wasn't committed

	// Prepare the SQL statement
	stmt, err := tx.Prepare(`
	INSERT INTO metrics(timestamp, network_name, online, metrics)
	VALUES(?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(
		metrics.Timestamp.UnixMilli(),
		metrics.NetworkName,
		metrics.Online,
		metricsJson,
	)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	logger.Infof("Saved metrics")
	return nil
}

func getAllMetrics() ([]Metrics, error) {
	rows, err := dbConnection.Query(`
		SELECT timestamp, network_name, online, metrics
		FROM metrics
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %v", err)
	}
	defer rows.Close()

	var allMetrics []Metrics
	for rows.Next() {
		var timestampMillis int64
		var networkName string
		var onlineInt int
		var metricsJson string

		if err = rows.Scan(&timestampMillis, &networkName, &onlineInt, &metricsJson); err != nil {
			logger.Errorf("failed to scan row: %v", err)
		}

		var metrics = Metrics{
			Timestamp:   time.UnixMilli(timestampMillis),
			NetworkName: networkName,
			Online:      onlineInt != 0,
		}

		err = json.Unmarshal([]byte(metricsJson), &metrics.Data)
		if err != nil {
			logger.Errorf("failed to unmarshal metrics: %v", err)
		} else {
			allMetrics = append(allMetrics, metrics)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return allMetrics, nil
}

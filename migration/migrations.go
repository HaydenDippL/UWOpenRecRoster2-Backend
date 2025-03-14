package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Session struct {
	SessionID     string
	IP            string
	NumQueries    string
	DateOfQueries string
	TimeOfQueries string
	Device        string
	Browser       string
}

type Query struct {
	ID          string
	DateQueried string
	Gym         string
	GymFacility string
	SessionID   string
}

func readCSVFiles() ([]Session, []Query, error) {
	// Get current directory
	migrationDir, err := os.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Read sessions CSV
	sessionsFile, err := os.Open(filepath.Join(migrationDir, "test.sessions.0000000010000.csv"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open sessions file: %w", err)
	}
	defer sessionsFile.Close()

	// Read queries CSV
	queriesFile, err := os.Open(filepath.Join(migrationDir, "test.queries.0000000010000.csv"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open queries file: %w", err)
	}
	defer queriesFile.Close()

	// Parse sessions CSV
	sessionsReader := csv.NewReader(sessionsFile)
	sessionsRecords, err := sessionsReader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read sessions CSV: %w", err)
	}

	// Parse queries CSV
	queriesReader := csv.NewReader(queriesFile)
	queriesRecords, err := queriesReader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read queries CSV: %w", err)
	}

	// Convert records to structs (skip header row)
	var sessions []Session
	for _, record := range sessionsRecords[1:] {
		session := Session{
			SessionID:     record[0],
			IP:            record[1],
			NumQueries:    record[2],
			DateOfQueries: record[3],
			TimeOfQueries: record[4],
			Device:        record[5],
			Browser:       record[6],
		}
		sessions = append(sessions, session)
	}

	var queries []Query
	for _, record := range queriesRecords[1:] {
		query := Query{
			ID:          record[0],
			DateQueried: record[1],
			Gym:         record[2],
			GymFacility: record[3],
			SessionID:   record[4],
		}
		queries = append(queries, query)
	}

	log.Printf("Successfully loaded %d sessions", len(sessions))
	log.Printf("Successfully loaded %d queries", len(queries))

	return sessions, queries, nil
}

func main() {
	sessions, queries, err := readCSVFiles()
	if err != nil {
		log.Fatalf("Error reading CSV files: %v", err)
	}

	// Preview the data
	log.Println("\nSessions Preview:")
	for i := 0; i < 5 && i < len(sessions); i++ {
		log.Printf("%+v", sessions[i])
	}

	log.Println("\nQueries Preview:")
	for i := 0; i < 5 && i < len(queries); i++ {
		log.Printf("%+v", queries[i])
	}
}

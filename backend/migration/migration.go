package main

import (
	"UWOpenRecRoster2-Backend/models"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LegacySession struct {
	SessionId string
	IP        string
	Date      string
	Time      string
}

type LegacyQueries struct {
	Id        int
	Date      string
	Gym       string
	SessionId string
	Facility  string
}

func main() {
	// legacySessionsCSVPath := "test.sessions.0000000010000.csv"
	// legacyQueriesCSVPath := "test.queries.0000000010000.csv"
	legacySessionsCSVPath := "sessions.dummy.csv"
	legacyQueriesCSVPath := "queries.dummy.csv"
	legacySessions, legacyQueries, err := readCSVFiles(legacySessionsCSVPath, legacyQueriesCSVPath)
	if err != nil {
		log.Fatalf("Error reading CSV files: %v", err)
	}

	users, sessions, queries := convertData(legacySessions, legacyQueries)

	fmt.Println("Postgres Insertion...")

	err = insertData(users, sessions, queries)
	if err != nil {
		log.Printf("Error migrating data: %v", err)
	}
}

func readCSVFiles(legacySessionsCSVPath string, legacyQueriesCSVPath string) ([]LegacySession, []LegacyQueries, error) {
	// Get current directory
	migrationDir, err := os.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Read sessions CSV
	sessionsFile, err := os.Open(filepath.Join(migrationDir, legacySessionsCSVPath))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open sessions file: %w", err)
	}
	defer sessionsFile.Close()

	// Read queries CSV
	queriesFile, err := os.Open(filepath.Join(migrationDir, legacyQueriesCSVPath))
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
	var sessions []LegacySession
	for _, record := range sessionsRecords[1:] {
		session := LegacySession{
			SessionId: record[0],
			IP:        record[1],
			Date:      record[3],
			Time:      record[4],
		}
		sessions = append(sessions, session)
	}

	var queries []LegacyQueries
	for _, record := range queriesRecords[1:] {
		query := LegacyQueries{
			Date:      record[1],
			Gym:       record[2],
			Facility:  record[3],
			SessionId: record[4],
		}
		queries = append(queries, query)
	}

	log.Printf("Successfully loaded %d legacy sessions", len(sessions))
	log.Printf("Successfully loaded %d legacy queries", len(queries))

	return sessions, queries, nil
}

func convertData(legacySessions []LegacySession, legacyQueries []LegacyQueries) ([]models.User, []models.Session, []models.Query) {
	ipToUserId := make(map[string]string)
	sessionIdToUserId := make(map[string]string)
	queryIdToSessionId := make(map[int]string)

	legacySessionIdToLegacySession := make(map[string]LegacySession)
	userIdToUser := make(map[string]models.User)

	sessionIdToQueries := make(map[string][]models.Query)
	userIdToSessions := make(map[string][]models.Session)

	queries := make([]models.Query, 0)
	sessions := make([]models.Session, 0)
	users := make([]models.User, 0)

	// populate legacySessionIdToLegacySession
	for _, legacySession := range legacySessions {
		var sessionId string = legacySession.SessionId
		legacySessionIdToLegacySession[sessionId] = legacySession
	}

	// populate ipToUserId
	for _, legacySession := range legacySessions {
		var ip string = legacySession.IP
		if _, exists := ipToUserId[ip]; !exists {
			var userId string = uuid.New().String()
			ipToUserId[ip] = userId
		}
	}

	// populate sessionIdToUserId
	for _, legacySession := range legacySessions {
		// get sessionId and make sure it isn't a duplicate
		var sessionId string = legacySession.SessionId
		if _, exists := sessionIdToUserId[sessionId]; exists {
			log.Fatalf("Two entries with the same sessionId: %s\n", sessionId)
		}

		// get ip, translate to userId, insert into sessionIdToUserId
		var ip string = legacySession.IP
		if userId, exists := ipToUserId[ip]; !exists {
			log.Fatalf("Could not find userId for IP: %s\n", ip)
		} else {
			sessionIdToUserId[sessionId] = userId
		}
	}

	// put LegacyQueries[i].Id from [1:...] and populate queryIdToSessionId
	for i, legacyQuery := range legacyQueries {
		var id int = i + 1
		var sessionId string = legacyQuery.SessionId

		legacyQueries[i].Id = id
		queryIdToSessionId[id] = sessionId
	}

	// Construct queries, populate sessionIdToQueries
	for _, legacyQuery := range legacyQueries {
		// lowercase gym and ensure it is "bakke" or "nick"
		var gym string = strings.ToLower(legacyQuery.Gym)
		if gym != "bakke" && gym != "nick" {
			log.Fatalf("Gym must be \"bakke\" or \"nick\", found \"%s\"\n", gym)
		}

		// retrieve the date-time of the session creation and use as query time
		var sessionId string = queryIdToSessionId[legacyQuery.Id]
		var session LegacySession = legacySessionIdToLegacySession[sessionId]
		queriedTime, err := time.Parse("2006-01-02 15:04:05", session.Date+" "+session.Time)
		if err != nil {
			log.Fatalf("Failed to parse date (%s) and time (%s): %v", session.Date, session.Time, err)
		}

		// parse the queried court date
		scheduleDate, _ := time.Parse("2006-01-02", legacyQuery.Date)

		query := models.Query{
			Id:           legacyQuery.Id,
			SessionID:    legacyQuery.SessionId,
			ScheduleDate: scheduleDate,
			QueriedTime:  queriedTime,
		}

		queries = append(queries, query)
		sessionIdToQueries[query.SessionID] = append(sessionIdToQueries[query.SessionID], query)
	}

	// created sessions and append to userIdToSessions
	for _, legacySession := range legacySessions {
		created, err := time.Parse("2006-01-02 15:04:05", legacySession.Date+" "+legacySession.Time)
		if err != nil {
			log.Fatalf("Failed to parse date (%s) and time (%s): %v", legacySession.Date, legacySession.Time, err)
		}

		var sessionId string = legacySession.SessionId

		userId, exists := sessionIdToUserId[sessionId]
		if !exists {
			log.Fatalf("UserId does not exist for sessionId: %s\n", sessionId)
		}

		sessionQueries, exists := sessionIdToQueries[sessionId]
		if !exists {
			log.Fatalf("Quereis do not exist for sessionId: %s\n", sessionId)
		}

		session := models.Session{
			SessionID: legacySession.SessionId,
			UserID:    userId,
			Created:   created,
			Queries:   sessionQueries,
		}

		sessions = append(sessions, session)
		userIdToSessions[session.UserID] = append(userIdToSessions[session.UserID], session)
	}

	// parially create users and get the earliest occurence of user
	for _, legacySession := range legacySessions {
		created, err := time.Parse("2006-01-02 15:04:05", legacySession.Date+" "+legacySession.Time)
		if err != nil {
			log.Fatalf("Failed to parse date (%s) and time (%s): %v", legacySession.Date, legacySession.Time, err)
		}

		var ip string = legacySession.IP
		userId, exists := ipToUserId[ip]
		if !exists {
			log.Fatalf("user does not exist for IP: %s\n", ip)
		}

		user, exists := userIdToUser[userId]
		if !exists || created.Before(user.Created) {
			userIdToUser[userId] = models.User{
				UserID:  userId,
				Created: created,
			}
		}
	}

	// complete the users
	for _, user := range userIdToUser {
		user.Sessions = userIdToSessions[user.UserID]
		users = append(users, user)
	}

	return users, sessions, queries
}

func insertData(users []models.User, sessions []models.Session, queries []models.Query) error {
	// Connect to the database
	dsn := getDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the models
	if err := db.AutoMigrate(&models.User{}, &models.Session{}, &models.Query{}, &models.Schedule{}); err != nil {
		return fmt.Errorf("failed to auto migrate models: %w", err)
	}

	// Start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Defer rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete all existing data in reverse order of dependencies
	log.Println("Deleting existing data...")
	if err := tx.Exec("DELETE FROM queries").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete queries: %w", err)
	}
	if err := tx.Exec("DELETE FROM sessions").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete sessions: %w", err)
	}
	if err := tx.Exec("DELETE FROM users").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete users: %w", err)
	}

	const batchSize = 100

	// Insert users in batches
	log.Printf("Inserting %d users in batches of %d...", len(users), batchSize)
	for i := 0; i < len(users); i += batchSize {
		end := min(i+batchSize, len(users))
		if err := tx.Select("user_id", "created").Create(users[i:end]).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert users batch: %w", err)
		}
	}

	// Insert sessions in batches
	log.Printf("Inserting %d sessions in batches of %d...", len(sessions), batchSize)
	for i := 0; i < len(sessions); i += batchSize {
		end := min(i+batchSize, len(sessions))
		if err := tx.Select("session_id", "user_id", "created").Create(sessions[i:end]).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert sessions batch: %w", err)
		}
	}

	// Insert queries in batches
	log.Printf("Inserting %d queries in batches of %d...", len(queries), batchSize)
	for i := 0; i < len(queries); i += batchSize {
		end := min(i+batchSize, len(queries))
		if err := tx.Create(queries[i:end]).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert queries batch: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully inserted all data")
	return nil
}

func getDSN() string {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf(".env file not found")
	}

	// Verfiy .env file has expected variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		log.Fatalf(
			"One or more of the environment variables was empty or doesn't exist\n"+
				"\tDB_HOST: \"%s\"\n"+
				"\tDB_USER: \"%s\"\n"+
				"\tDB_PASSWORD: \"%s\"\n"+
				"\tDB_NAME: \"%s\"\n"+
				"\tDB_PORT: \"%s\"\n",
			host, user, password, dbname, port,
		)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Chicago",
		host, user, password, dbname, port,
	)

	return dsn
}

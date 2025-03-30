package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"UWOpenRecRoster2-Backend/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	initDB()

	r := gin.Default()

	r.Use(middleware)

	r.GET("/", hello_world)
	r.GET("/schedule", schedule)

	r.Run(":8000")
}

func initDB() {
	err := godotenv.Load()
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

	// Format DSN and get DB object
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Chicago",
		host, user, password, dbname, port,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}

	// Get the underlying SQL DB object
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database: %v", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto migrate your models
	DB.AutoMigrate(&models.User{}, &models.Session{}, &models.Query{}, &models.Schedule{})
}

func middleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Writer.Header().Set("Content-Type", "application/json")

	c.Next()
}

func hello_world(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
		"message": "Hello World!",
	})
}

func schedule(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date parameter is required"})
		return
	}

	dateFormat := "2006-01-02"
    parsedDate, err := time.Parse(dateFormat, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date parameter must be of the form yyyy-MM-dd"})
		return
	}

	gyms := c.Query("gym")
	if gyms == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gym parameter is required"})
		return
	}

	gymList := strings.Split(gyms, ",")
	if len(gymList) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gym parameter must be a non-empty comma separated list of gyms: \"bakke\", \"nick\""})
		return
	}

	for _, gym := range gymList {
		if gym != "bakke" && gym != "nick" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "gym parameter must be a non-empty comma separated list of gyms: \"bakke\", \"nick\""})
			return
		}
	}

    schedules := make(map[string]models.ScheduleJSON)

    for _, gym := range gymList {
        fmt.Printf("getting\n")
        schedule, err := getSchedule(date, gym)
        if err != nil {
            fmt.Printf("%v\n", err)
            fmt.Printf("fetching\n")
            schedule, err = fetchSchedule(date, gym)
            if err != nil {
                fmt.Printf("err on fetch\n")
                c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong on our end"})
                return
            } else if err = memoSchedule(schedule, parsedDate, gym); err != nil {
                fmt.Printf("err on memorize\n")
            }
        }
        
        schedules[gym] = schedule
    }

    c.JSON(http.StatusOK, schedules)
}

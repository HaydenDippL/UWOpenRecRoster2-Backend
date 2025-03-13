package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(middleware)

	r.GET("/", hello_world)
	r.GET("/schedule", schedule)

	r.Run(":8000")
}

func middleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
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
	_, err := time.Parse(dateFormat, date)
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

	for _, gym := range gymList {
		if gym == "bakke" {
			fetchSchedule(bakke, date)
		} else {
			fetchSchedule(nick, date)
		}
	}
}

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"event-registration-api/internal/db"
	"event-registration-api/internal/event"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	r := gin.Default()

	event.RegisterRoutes(r)

	r.Run(":8080")
}
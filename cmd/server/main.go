package main

import (
	"log"

	"event-registration-api/internal/db"
	"event-registration-api/internal/event"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	r := gin.Default()
	event.RegisterRoutes(r)
	r.Run(":8080")
}

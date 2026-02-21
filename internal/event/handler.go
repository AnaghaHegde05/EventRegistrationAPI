package event

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/events", createEvent)
	r.GET("/events", listEvents)
	r.POST("/events/:id/register", registerEvent)
}

func createEvent(c *gin.Context) {
	var req struct {
		Name  string `json:"name"`
		Seats int    `json:"seats"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := CreateNewEvent(req.Name, req.Seats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "event created"})
}

func listEvents(c *gin.Context) {
	events, err := GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

func registerEvent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := Register(id); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "event full"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}
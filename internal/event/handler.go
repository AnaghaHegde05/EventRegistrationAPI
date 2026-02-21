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
	r.GET("/events/:id/registrations", getRegistrations) // new
	r.DELETE("/registrations/:id", cancelRegistration)   // new
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	var req struct {
		UserEmail string `json:"user_email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "valid user_email required"})
		return
	}

	if err := Register(id, req.UserEmail); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

// GET /events/:id/registrations
func getRegistrations(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	registrations, err := GetRegistrationsByEventID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, registrations)
}

// DELETE /registrations/:id
func cancelRegistration(c *gin.Context) {
	regID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid registration id"})
		return
	}

	if err := CancelRegistration(regID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registration cancelled"})
}

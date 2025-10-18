package handler

import (
	"log"
	"net/http"

	"github.com/Met-String/AgentSquare/internal/extension_observer/storage"
	"github.com/gin-gonic/gin"
)

// extensionEventPayload captures the event fields sent by an extension observer.
type extensionEventPayload struct {
	Event    string `json:"event"`
	ClientID string `json:"clientId"`
	Time     string `json:"time"`
}

// ExtensionEventHandler reads event data from the request body and logs the key fields.
func ExtensionEventHandler(c *gin.Context) {
	var payload extensionEventPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("extension observer: failed to parse request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	log.Printf("extension observer event=%s clientId=%s time=%s", payload.Event, payload.ClientID, payload.Time)

	if err := storage.SaveExtensionEvent(c.Request.Context(), storage.ExtensionEventDocument{
		Event:    payload.Event,
		ClientID: payload.ClientID,
		Time:     payload.Time,
	}); err != nil {
		log.Printf("extension observer: failed to persist event: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist event"})
		return
	}

	c.Status(http.StatusOK)
}

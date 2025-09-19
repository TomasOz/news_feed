package feed

import (
	"net/http"
	"github.com/gin-gonic/gin"

)

type FeedHandler struct {
	service FeedService
}

func NewFeedHandler(service FeedService) *FeedHandler {
	return &FeedHandler{service: service}
}

func (h *FeedHandler) GetFeed(c *gin.Context) {
	currentUserID, exists := c.Get("ID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := currentUserID.(uint)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid User ID"})
		return
	}

	posts, err := h.service.GetFeed(userID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
} 
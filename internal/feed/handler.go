package feed

import (
	"net/http"
	"github.com/gin-gonic/gin"

	"strconv"
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

	limit := 10

	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	
	cursor := c.Query("cursor")

	posts, nextCursor, err := h.service.GetFeed(userID, limit, cursor)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"next_cursor": nextCursor,
	})
} 
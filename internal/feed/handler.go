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


// GetFeed godoc
// @Summary     Get user feed
// @Description Returns the authenticated user's feed with optional limit and offset pagination.
// @Tags        feed
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       limit   query     int    false  "Max posts to return" minimum(1) maximum(100) default(10)
// @Param       offset  query     int    false  "Number of posts to skip" minimum(0) default(0)
// @Success     200     {object}  FeedResponse
// @Failure     400     {object}  ErrorResponse
// @Failure     401     {object}  ErrorResponse
// @Failure     500     {object}  ErrorResponse
// @Router      /feed [get]
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
	
	offset := 0
	if o := c.Query("offset"); o != "" {
		if n, err := strconv.Atoi(o); err == nil && n >= 0 {
			offset = n
		}
	}

	posts, err := h.service.GetFeed(userID, limit, offset)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hasMore := len(posts) == limit

	c.JSON(http.StatusOK, gin.H{"posts": posts, "has_more": hasMore})
} 
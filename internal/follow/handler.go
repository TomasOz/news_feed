package follow

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)

type FollowHandler struct {
	service FollowService
}

func NewFollowHandler(service FollowService) *FollowHandler {
	return &FollowHandler{service: service}
}

func (h *FollowHandler) Follow(c *gin.Context) {
	user_id := c.Param("id")

    idUint64, err := strconv.ParseUint(user_id, 10, 64)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid ID"})
        return
    }

    followee_id := uint(idUint64)

	userIDValue, exists := c.Get("ID")
	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	follower_id, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user id in token"})
        return
	}

	error := h.service.Follow(follower_id, followee_id)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"You start following:": user_id})
}

func (h *FollowHandler) UnFollow(c *gin.Context) {
	user_id := c.Param("id")

    idUint64, err := strconv.ParseUint(user_id, 10, 64)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid ID"})
        return
    }

    followee_id := uint(idUint64)

	userIDValue, exists := c.Get("ID")
	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	follower_id, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user id in token"})
        return
	}

	error := h.service.UnFollow(follower_id, followee_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"You unfollowed:": user_id})
}
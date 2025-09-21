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
	userId := c.Param("id")

    idUint64, err := strconv.ParseUint(userId, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    followeeId := uint(idUint64)

	userIDValue, exists := c.Get("ID")
	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	followerId, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user id in token"})
        return
	}

	err = h.service.Follow(followerId, followeeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"You start following:": userId})
}

func (h *FollowHandler) UnFollow(c *gin.Context) {
	userId := c.Param("id")

    idUint64, err := strconv.ParseUint(userId, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    followee_id := uint(idUint64)

	userIDValue, exists := c.Get("ID")
	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	follower_id, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user id in token"})
        return
	}

	err = h.service.UnFollow(follower_id, followee_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"You unfollowed:": userId})
}
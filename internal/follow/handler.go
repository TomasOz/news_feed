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

// Follow godoc
// @Summary     Follow another user
// @Description Authenticated user follows onether user
// @Tags        follow
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id   path      int  true  "ID of the user to follow"
// @Success     200     {json} error {"You start following:":"123"}  
// @Failure     400     {json} error {"error":"Invalid ID"}
// @Failure     401     {json} error {"error":"Unauthorized"}
// @Failure     500     {json} error {"error":"Invalid user id in token"}
// @Router      /follow [get]
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
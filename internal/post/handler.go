package post

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"

	"news-feed/internal/background_jobs"
)

type PostHandler struct {
	service 		PostService
	fanoutWorker 	*background_jobs.FanoutWorker
}

func NewPostHandler(service PostService, fanoutWorker *background_jobs.FanoutWorker) *PostHandler {
	return &PostHandler{service: service, fanoutWorker: fanoutWorker }
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.service.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Post was not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := h.service.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}


func (h *PostHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
		
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post"})
		return
	}

	userIDValue, exists := c.Get("ID")
	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user id in token"})
        return
	}

	post, err := h.service.Create(userID, req.Title, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func() {
		if err := h.fanoutWorker.FanoutToAllFollowers(post.ID, userID); err != nil {
			println("Failed to fanout post:", err.Error())
		}
	}()

	c.JSON(http.StatusOK, post)
}


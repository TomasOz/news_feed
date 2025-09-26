package user

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"news-feed/internal/auth"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterUser godoc
// @Summary      Register a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  user.UserInput  true  "User credentials"
// @Success      201    {object}  user.UserResponse
// @Failure      400    {object}  map[string]string
// @Router       /register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var input UserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapUserToResponse(user))
}

// LoginUser godoc
// @Summary     Login
// @Description Authenticate a user with username/password and receive a JWT.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       input  body  UserInput  true  "Credentials"
// @Success     200    {object}  TokenResponse
// @Failure     400    {object}  ErrorResponse
// @Failure     401    {object}  ErrorResponse
// @Failure     500    {object}  ErrorResponse
// @Router      /login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var input UserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Login(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.GenerateJWT(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

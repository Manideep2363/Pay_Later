package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"paylater/internal/service"
)

// UserHandler handles all HTTP requests related to users.
type UserHandler struct {
	// service contains the business logic.
	service *service.UserService
}

// NewUserHandler creates and returns a new UserHandler.
func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

// CreateUserRequest represents the JSON request body
// expected from the client when creating a new user.
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUser handles POST /users requests.
func (h *UserHandler) CreateUser(c *gin.Context) {

	// Variable to store the incoming JSON data.
	var req CreateUserRequest

	// Read and bind the JSON request body into req.
	if err := c.ShouldBindJSON(&req); err != nil {

		// Return 400 Bad Request if the JSON is invalid.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Call the service layer to create the user.
	// The handler only receives the request and sends the response.
	err := h.service.CreateUser(
		c.Request.Context(),
		req.Name,
		req.Email,
	)

	// If the service returns an error, send it to the client.
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return success response.
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

// ListUsers handles GET /users requests.
func (h *UserHandler) ListUsers(c *gin.Context) {

	// Ask the service layer to fetch all users.
	users, err := h.service.ListUsers(c.Request.Context())

	// If something goes wrong, return 500 Internal Server Error.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the list of users as JSON.
	c.JSON(http.StatusOK, users)
}
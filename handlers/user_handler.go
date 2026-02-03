package handlers

import (
	"context"
	"follooow-be/models"
	"follooow-be/repositories"
	"follooow-be/responses"
	"follooow-be/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(c echo.Context) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.CreateUserModel
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"error": err.Error()},
		})
	}

	// Validate input
	if user.Username == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"error": "username and password are required"},
		})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"error": "failed to hash password"},
		})
	}

	// Set hashed password
	user.Password = hashedPassword

	// Create user
	newUser, err := repositories.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"error": err.Error()},
		})
	}

	if newUser == nil {
		return c.JSON(http.StatusConflict, responses.GlobalResponse{
			Status:  http.StatusConflict,
			Message: "error",
			Data:    &echo.Map{"error": "username already exists"},
		})
	}

	// Prepare response (without password)
	userResponse := models.UserResponse{
		ID:        newUser.ID,
		Username:  newUser.Username,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, responses.GlobalResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &echo.Map{"user": userResponse},
	})
}

func GetUserByID(c echo.Context) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.Param("user_id")
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"error": "invalid user ID"},
		})
	}

	user, err := repositories.FindUserByID(objId)
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.GlobalResponse{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    &echo.Map{"error": "user not found"},
		})
	}


	// Prepare response (without password)
	userResponse := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.JSON(http.StatusOK, responses.GlobalResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &echo.Map{"user": userResponse},
	})
}

func LoginUser(c echo.Context) error {
	var loginReq models.LoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"error": err.Error()},
		})
	}

	// Validate input
	if loginReq.Username == "" || loginReq.Password == "" {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"error": "username and password are required"},
		})
	}

	// Find user by username
	user, err := repositories.FindUserByUsername(loginReq.Username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, responses.GlobalResponse{
			Status:  http.StatusUnauthorized,
			Message: "error",
			Data:    &echo.Map{"error": "invalid username or password"},
		})
	}

	// Check password
	if !utils.CheckPasswordHash(loginReq.Password, user.Password) {
		return c.JSON(http.StatusUnauthorized, responses.GlobalResponse{
			Status:  http.StatusUnauthorized,
			Message: "error",
			Data:    &echo.Map{"error": "invalid username or password"},
		})
	}

	// Prepare login response
	loginResponse := models.LoginResponse{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		Message:  "Login successful",
	}

	return c.JSON(http.StatusOK, responses.GlobalResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &echo.Map{"login": loginResponse},
	})
}

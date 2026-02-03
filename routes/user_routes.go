package routes

import (
	"follooow-be/handlers"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	// all routes relates to users comes here
	e.POST("/api/users", handlers.CreateUser)
	e.POST("/api/users/login", handlers.LoginUser)
	e.GET("/api/users/:user_id", handlers.GetUserByID)
}

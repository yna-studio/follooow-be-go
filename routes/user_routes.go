package routes

import (
	"follooow-be/handlers"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	// all routes relates to users comes here
	e.POST("/users", handlers.CreateUser)
	e.POST("/users/login", handlers.LoginUser)
	e.GET("/users/:user_id", handlers.GetUserByID)
}

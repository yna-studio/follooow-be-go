package routes

import (
	"follooow-be/handlers"

	"github.com/labstack/echo/v4"
)

// MediaRoute defines all media-related routes
func MediaRoute(e *echo.Echo) {
	// Media upload route
	e.POST("/api/media/upload", handlers.UploadMedia)
}

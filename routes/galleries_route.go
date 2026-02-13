package routes

import (
	"follooow-be/handlers"

	"github.com/labstack/echo/v4"
)

func GalleriesRoute(e *echo.Echo) {
	// all routes relates to influencers comes here
	e.GET("/galleries", handlers.ListGalleries)
	e.GET("/galleries/:gallery_id", handlers.DetailGallery)
	e.POST("/galleries", handlers.CreateGallery)
	e.POST("/galleries/upload", handlers.CreateGalleryWithUpload)
}

package routes

import (
    "follooow-be/handlers"
    "github.com/labstack/echo/v4"
)

// SyncRoute registers the label synchronization endpoint.
func SyncRoute(e *echo.Echo) {
    e.GET("/sync/labels", handlers.SyncLabels)
}

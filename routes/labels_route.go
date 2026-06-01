package routes

import (
    "follooow-be/handlers"
    "github.com/labstack/echo/v4"
)

func LabelsRoute(e *echo.Echo) {
    e.GET("/labels", handlers.ListLabels)
}

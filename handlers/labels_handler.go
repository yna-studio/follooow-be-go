package handlers

import (
    "context"
    "time"
    "strconv"
    "follooow-be/repositories"
    "follooow-be/responses"
    "net/http"
    "github.com/labstack/echo/v4"
)

// ListLabels handles GET /labels and returns all labels with their counts.
func ListLabels(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Parse optional "limit" query parameter
    limitStr := c.QueryParam("limit")
    var limit int64 = 0
    if limitStr != "" {
        if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil && l > 0 {
            limit = l
        }
    }

    // Fetch labels with optional limit, sorted by total descending
    labels, err := repositories.ListLabels(ctx, limit)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
    }
    return c.JSON(http.StatusOK, responses.GlobalResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"labels": labels}})
}

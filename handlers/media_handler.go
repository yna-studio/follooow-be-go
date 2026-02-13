package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"follooow-be/configs"
	"follooow-be/responses"
	"follooow-be/utils"

	"github.com/labstack/echo/v4"
)

// MediaUploadPayload represents the request payload for media upload
type MediaUploadPayload struct {
	File      string `json:"file" validate:"required"`
	Directory string `json:"directory" validate:"required"`
}

// UploadMedia handles single file upload from base64 data
func UploadMedia(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var payload MediaUploadPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "Error parsing request body",
			Data:    &echo.Map{"error": err.Error()},
		})
	}

	// Validate required fields
	if payload.File == "" {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "File is required",
			Data:    nil,
		})
	}

	if payload.Directory == "" {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "Directory is required",
			Data:    nil,
		})
	}

	// Initialize Cloudinary client if not already done
	if configs.CloudinaryClient == nil {
		configs.InitCloudinary()
	}

	// Construct the full directory path
	fullDirectory := configs.EnvCloudinaryDir() + "/" + strings.TrimPrefix(payload.Directory, "/")

	// Generate unique filename
	filename := fmt.Sprintf("media_%d", time.Now().Unix())

	// Upload the base64 image to Cloudinary
	result, err := utils.UploadImageFromBase64(ctx, payload.File, fullDirectory, filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error uploading file to Cloudinary",
			Data:    &echo.Map{"error": err.Error()},
		})
	}

	// Return success response with CDN URL
	return c.JSON(http.StatusOK, responses.GlobalResponse{
		Status:  http.StatusOK,
		Message: "File uploaded successfully",
		Data: &echo.Map{
			"url":       result.SecureURL,
			"public_id": result.PublicID,
			"format":    result.Format,
			"size":      result.Bytes,
			"directory": fullDirectory,
		},
	})
}

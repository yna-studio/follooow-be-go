package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"follooow-be/configs"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadImageFromForm uploads an image from form data to Cloudinary
func UploadImageFromForm(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (*uploader.UploadResult, error) {
	if fileHeader == nil {
		return nil, fmt.Errorf("no file provided")
	}

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Generate unique filename
	filename := generateUniqueFilename(fileHeader.Filename)

	// Set folder if provided
	if folder == "" {
		folder = configs.EnvCloudinaryDir()
	}

	uploadParams := uploader.UploadParams{
		Folder:       folder,
		PublicID:     filename,
		ResourceType: "image",
	}

	result, err := configs.CloudinaryClient.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	return result, nil
}

// UploadImageFromBase64 uploads a base64 encoded image to Cloudinary
func UploadImageFromBase64(ctx context.Context, base64Data string, folder string, filename string) (*uploader.UploadResult, error) {
	if base64Data == "" {
		return nil, fmt.Errorf("no base64 data provided")
	}

	// Remove data URL prefix if present
	if strings.HasPrefix(base64Data, "data:") {
		parts := strings.SplitN(base64Data, ",", 2)
		if len(parts) == 2 {
			base64Data = parts[1]
		}
	}

	// Decode base64
	imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Set folder if provided
	if folder == "" {
		folder = configs.EnvCloudinaryDir()
	}

	// Create a reader from the byte slice for Cloudinary
	reader := bytes.NewReader(imageBytes)

	uploadParams := uploader.UploadParams{
		Folder:       folder,
		PublicID:     filename,
		ResourceType: "image",
	}

	result, err := configs.CloudinaryClient.Upload.Upload(ctx, reader, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload base64 image: %w", err)
	}

	return result, nil
}

// UploadImageFromURL uploads an image from URL to Cloudinary
func UploadImageFromURL(ctx context.Context, imageURL string, filename string) (*uploader.UploadResult, error) {
	if imageURL == "" {
		return nil, fmt.Errorf("no image URL provided")
	}

	uploadParams := uploader.UploadParams{
		Folder:       configs.EnvCloudinaryDir(),
		PublicID:     filename,
		ResourceType: "image",
	}

	result, err := configs.CloudinaryClient.Upload.Upload(ctx, imageURL, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image from URL: %w", err)
	}

	return result, nil
}

// DeleteImageFromCloudinary deletes an image from Cloudinary
func DeleteImageFromCloudinary(ctx context.Context, publicID string) (*uploader.DestroyResult, error) {
	if publicID == "" {
		return nil, fmt.Errorf("no public ID provided")
	}

	destroyParams := uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
	}

	result, err := configs.CloudinaryClient.Upload.Destroy(ctx, destroyParams)
	if err != nil {
		return nil, fmt.Errorf("failed to delete image: %w", err)
	}

	return result, nil
}

// generateUniqueFilename generates a unique filename for upload
func generateUniqueFilename(originalFilename string) string {
	// Get file extension
	parts := strings.Split(originalFilename, ".")
	extension := ""
	if len(parts) > 1 {
		extension = parts[len(parts)-1]
	}

	// Generate timestamp for uniqueness
	timestamp := time.Now().Unix()

	// Remove special characters from original filename
	baseName := strings.ReplaceAll(originalFilename, ".", "_")
	baseName = strings.ReplaceAll(baseName, " ", "_")

	// Create unique filename
	if extension != "" {
		return fmt.Sprintf("%s_%d.%s", baseName, timestamp, extension)
	}
	return fmt.Sprintf("%s_%d", baseName, timestamp)
}

// GetPublicIDFromURL extracts public ID from Cloudinary URL
func GetPublicIDFromURL(imageURL string) string {
	// Example URL: https://res.cloudinary.com/cloud_name/image/upload/v1234567890/folder/image_name.jpg
	parts := strings.Split(imageURL, "/")

	// Find folder and image parts
	for i, part := range parts {
		if part == "upload" && i+2 < len(parts) {
			// Skip version part (v1234567890)
			publicID := strings.Join(parts[i+2:], "/")
			// Remove file extension for deletion
			if dotIndex := strings.LastIndex(publicID, "."); dotIndex != -1 {
				publicID = publicID[:dotIndex]
			}
			return publicID
		}
	}

	return ""
}

package configs

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Cloudinary client instance
var CloudinaryClient *cloudinary.Cloudinary

// InitCloudinary initializes the Cloudinary client
func InitCloudinary() {
	cld, err := cloudinary.NewFromURL(
		fmt.Sprintf("cloudinary://%s:%s@%s",
			EnvCloudinaryAPIKey(),
			EnvCloudinaryAPISecret(),
			EnvCloudinaryCloudName(),
		))
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary: ", err)
	}

	CloudinaryClient = cld
	fmt.Println("Connected to Cloudinary")
}

// UploadImage uploads an image to Cloudinary and returns the upload result
func UploadImage(ctx context.Context, file interface{}, filename string) (*uploader.UploadResult, error) {
	if CloudinaryClient == nil {
		InitCloudinary()
	}

	uploadParams := uploader.UploadParams{
		Folder:   EnvCloudinaryDir(),
		PublicID: filename,
	}

	result, err := CloudinaryClient.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	return result, nil
}

// UploadImageWithURL uploads an image from URL to Cloudinary
func UploadImageWithURL(ctx context.Context, imageURL string, filename string) (*uploader.UploadResult, error) {
	if CloudinaryClient == nil {
		InitCloudinary()
	}

	uploadParams := uploader.UploadParams{
		Folder:   EnvCloudinaryDir(),
		PublicID: filename,
	}

	result, err := CloudinaryClient.Upload.Upload(ctx, imageURL, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image from URL: %w", err)
	}

	return result, nil
}

// DeleteImage deletes an image from Cloudinary
func DeleteImage(ctx context.Context, publicID string) (*uploader.DestroyResult, error) {
	if CloudinaryClient == nil {
		InitCloudinary()
	}

	destroyParams := uploader.DestroyParams{
		PublicID: publicID,
	}

	result, err := CloudinaryClient.Upload.Destroy(ctx, destroyParams)
	if err != nil {
		return nil, fmt.Errorf("failed to delete image: %w", err)
	}

	return result, nil
}

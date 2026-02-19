package handlers

import (
	"context"
	"encoding/json"
	"follooow-be/configs"
	"follooow-be/models"
	"follooow-be/responses"
	"follooow-be/utils"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UpdateGallery handles updating an existing gallery
func UpdateGallery(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	galleryID := c.Param("gallery_id")
	if galleryID == "" {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "Gallery ID is required",
			Data:    nil,
		})
	}

	// Convert gallery ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(galleryID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid gallery ID",
			Data:    &echo.Map{"error": err.Error()},
		})
	}

	// Parse request body
	var payload models.PayloadGallery
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "Error parsing request body",
			Data:    &echo.Map{"error": err.Error()},
		})
	}

	// Get existing gallery
	var existingGallery models.GalleryModel
	err = galleryCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&existingGallery)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, responses.GlobalResponse{
				Status:  http.StatusNotFound,
				Message: "Gallery not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error fetching gallery",
			Data:    &echo.Map{"error": err.Error()},
		})
	}

	// Prepare update data
	updateData := bson.M{
		"updated_on": time.Now().UnixNano() / int64(time.Millisecond),
	}

	// Update fields if provided
	if payload.Title != "" {
		updateData["title"] = payload.Title
		// Update slug if title changed
		slug := strings.Replace(payload.Title, " ", "-", -1)
		slug = strings.ToLower(slug)
		updateData["slug"] = slug
	}

	if payload.Description != "" {
		updateData["description"] = payload.Description
	}

	if payload.Lang != "" {
		updateData["lang"] = payload.Lang
	}

	if payload.Images != nil {
		updateData["images"] = payload.Images
	}

	if payload.Influencers != nil {
		updateData["influencers"] = payload.Influencers
	}

	if payload.Tags != nil {
		updateData["tags"] = payload.Tags
	}

	// Update gallery in database
	_, err = galleryCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error updating gallery",
			Data:    &echo.Map{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, responses.GlobalResponse{
		Status:  http.StatusOK,
		Message: "Gallery updated successfully",
		Data:    &echo.Map{"gallery_id": objID},
	})
}

// UpdateGalleryWithUpload handles updating gallery with image uploads
func UpdateGalleryWithUpload(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	galleryID := c.Param("gallery_id")
	if galleryID == "" {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "Gallery ID is required",
			Data:    nil,
		})
	}

	// Convert gallery ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(galleryID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid gallery ID",
			Data:    &echo.Map{"error": err.Error()},
		})
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{
			Status:  http.StatusBadRequest,
			Message: "Error parsing multipart form",
			Data:    &echo.Map{"error": err.Error()},
		})
	}

	// Get existing gallery
	var existingGallery models.GalleryModel
	err = galleryCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&existingGallery)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, responses.GlobalResponse{
				Status:  http.StatusNotFound,
				Message: "Gallery not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error fetching gallery",
			Data:    &echo.Map{"error": err.Error()},
		})
	}

	// Get form fields
	title := c.FormValue("title")
	description := c.FormValue("description")
	lang := c.FormValue("lang")
	influencersStr := c.FormValue("influencers")
	tagsStr := c.FormValue("tags")

	// Prepare update data
	updateData := bson.M{
		"updated_on": time.Now().UnixNano() / int64(time.Millisecond),
	}

	// Update fields if provided
	if title != "" {
		updateData["title"] = title
		// Update slug if title changed
		slug := strings.Replace(title, " ", "-", -1)
		slug = strings.ToLower(slug)
		updateData["slug"] = slug
	}

	if description != "" {
		updateData["description"] = description
	}

	if lang != "" {
		updateData["lang"] = lang
	}

	// Parse influencers if provided
	if influencersStr != "" {
		influencers := strings.Split(influencersStr, ",")
		for i, inf := range influencers {
			influencers[i] = strings.TrimSpace(inf)
		}
		updateData["influencers"] = influencers
	}

	// Parse tags if provided
	if tagsStr != "" {
		tags := strings.Split(tagsStr, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		updateData["tags"] = tags
	}

	// Handle image uploads if provided
	files := form.File["images"]
	if len(files) > 0 {
		// Initialize Cloudinary if not already done
		if configs.CloudinaryClient == nil {
			configs.InitCloudinary()
		}

		var images []models.ImageModel
		for i, file := range files {
			result, err := utils.UploadImageFromForm(ctx, file, "galleries")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{
					Status:  http.StatusInternalServerError,
					Message: "Error uploading image",
					Data:    &echo.Map{"error": err.Error()},
				})
			}

			imageModel := models.ImageModel{
				IsCover:   i == 0,
				Url:       result.SecureURL,
				Caption:   file.Filename,
				CreatedOn: int(time.Now().Unix()),
				UpdatedOn: int(time.Now().Unix()),
			}

			images = append(images, imageModel)
		}

		updateData["images"] = images
	}

	// Update gallery in database
	_, err = galleryCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error updating gallery",
			Data:    &echo.Map{"error": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, responses.GlobalResponse{
		Status:  http.StatusOK,
		Message: "Gallery updated successfully with images",
		Data:    &echo.Map{"gallery_id": objID},
	})
}

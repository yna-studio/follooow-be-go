package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"follooow-be/configs"
	"follooow-be/models"
	"follooow-be/repositories"
	"follooow-be/responses"
	"follooow-be/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var galleryCollection *mongo.Collection = configs.GetCollection(configs.DB, "galleries")
var galleryUsersCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var galleryInfluencersCollection *mongo.Collection = configs.GetCollection(configs.DB, "influencers")

// handler of GET /influencers
func ListGalleries(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var galleries []models.GalleryModel

	filterListData := bson.M{}

	var limit int64
	var page int64

	// handling limit, by default 6
	if c.QueryParam("limit") != "" {
		i, err := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
		}
		limit = i
	} else {
		limit = int64(6)
	}

	// handling page, by default 1
	if c.QueryParam("page") != "" {
		i, err := strconv.ParseInt(c.QueryParam("page"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
		}
		page = (i - 1) * limit
	} else {
		page = int64(0)
	}

	optsListData := options.Find().SetLimit(limit).SetSkip(page)

	// handling filter by language
	if c.QueryParam("lang") != "" {
		filterListData["lang"] = c.QueryParam("lang")
	}

	// handling filter by influencer id keyword [DONE]
	if c.QueryParam("influencer_ids") != "" {
		idsArr := strings.Split(c.QueryParam("influencer_ids"), ",")
		filterListData["influencers"] = bson.M{"$in": idsArr}
	}

	// by default sortby last update [DONE]
	if c.QueryParam("order_by") == "created_on" { //oldest created
		optsListData = optsListData.SetSort(bson.D{{"created_on", 1}})
	} else if c.QueryParam("order_by") == "created_on_new" { // latest created
		optsListData = optsListData.SetSort(bson.D{{"created_on", -1}})
	} else if c.QueryParam("order_by") == "popular" {
		optsListData = optsListData.SetSort(bson.D{{"views", -1}})
	} else {
		optsListData = optsListData.SetSort(bson.D{{"updated_on", -1}})
	}

	// get data from database
	results, err := galleryCollection.Find(ctx, filterListData, optsListData)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// get count data from database
	count, err := galleryCollection.CountDocuments(ctx, filterListData)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// reading data from db in an optimal way
	// defer use to delay execution
	defer results.Close(ctx)

	// normalize db results
	for results.Next(ctx) {
		var singleGallery models.GalleryModel
		var influencers []models.InfluencerSmallDataModel

		if err = results.Decode(&singleGallery); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
		}

		// convert influencers to data
		if len(singleGallery.Influencers) > 0 {
			// get all influencers on post
			// max result is 20
			optsListDataInfluencers := options.Find().SetLimit(20)

			// filter generator
			filterListDataInfluencers := bson.D{}

			idsArr := singleGallery.Influencers
			var idsObjId []primitive.ObjectID

			// normalize ids
			for key := range idsArr {
				objId, _ := primitive.ObjectIDFromHex(idsArr[key])
				idsObjId = append(idsObjId, objId)
			}

			filterListDataInfluencers = bson.D{{"_id", bson.M{"$in": idsObjId}}}

			// get data from database
			resultsInfluencers, err := galleryInfluencersCollection.Find(ctx, filterListDataInfluencers, optsListDataInfluencers)
			defer resultsInfluencers.Close(ctx)
			// normalize db results
			for resultsInfluencers.Next(ctx) {
				var singleInfluencer models.InfluencerSmallDataModel
				if err = resultsInfluencers.Decode(&singleInfluencer); err != nil {
					return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
				}

				influencers = append(influencers, singleInfluencer)
			}

			singleGallery.Influencers = nil
			singleGallery.InfluencersData = influencers
		}

		// get author information if author_id exists
		if singleGallery.AuthorID != "" {
			fmt.Printf("DEBUG: Found AuthorID: %s\n", singleGallery.AuthorID)
			authorObjID, err := primitive.ObjectIDFromHex(singleGallery.AuthorID)
			if err == nil {
				var author models.UserModel
				err := galleryUsersCollection.FindOne(ctx, bson.M{"_id": authorObjID}).Decode(&author)
				if err == nil {
					singleGallery.Author = &models.AuthorModel{
						ID:       author.ID.Hex(),
						Username: author.Username,
					}
					fmt.Printf("DEBUG: Author found: %s\n", author.Username)
				} else {
					fmt.Printf("DEBUG: Error finding author: %v\n", err)
				}
			} else {
				fmt.Printf("DEBUG: Error parsing AuthorID: %v\n", err)
			}
		} else {
			fmt.Printf("DEBUG: No AuthorID found for gallery: %s\n", singleGallery.Id.Hex())
		}

		galleries = append(galleries, singleGallery)
	}

	// check is no data available
	if len(galleries) < 1 {
		return c.JSON(http.StatusOK, responses.GlobalResponse{Status: http.StatusNoContent, Message: "Gallery not available", Data: &echo.Map{"galleries": galleries, "total": count}})
	} else {
		return c.JSON(http.StatusOK, responses.GlobalResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"galleries": galleries, "total": count}})
	}

}

// handle of GET /galleries/<id>
func DetailGallery(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get influencer_id
	galleryId := c.Param("gallery_id")
	var gallery models.GalleryModel
	var influencers []models.InfluencerSmallDataModel

	objId, _ := primitive.ObjectIDFromHex(galleryId)

	filterListData := bson.M{}

	filterListData["_id"] = objId

	// handling filter by language
	if c.QueryParam("lang") != "" {
		filterListData["lang"] = c.QueryParam("lang")
	}

	err := galleryCollection.FindOne(ctx, filterListData).Decode(&gallery)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// update views + 1
	_, err = galleryCollection.UpdateOne(ctx, bson.D{{"_id", objId}}, bson.D{{"$set", bson.D{{"views", gallery.Views + 1}}}})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
	}

	// get influencers data

	// convert influencers to data
	if len(gallery.Influencers) > 0 {
		// get all influencers on post
		// max result is 20
		optsListDataInfluencers := options.Find().SetLimit(20)

		// filter generator
		filterListDataInfluencers := bson.D{}

		idsArr := gallery.Influencers
		var idsObjId []primitive.ObjectID

		// normalize ids
		for key := range idsArr {
			objId, _ := primitive.ObjectIDFromHex(idsArr[key])
			idsObjId = append(idsObjId, objId)
		}

		// filter generator
		filterListDataInfluencers = bson.D{{"_id", bson.M{"$in": idsObjId}}}

		// get data from database
		resultsInfluencers, err := galleryInfluencersCollection.Find(ctx, filterListDataInfluencers, optsListDataInfluencers)
		defer resultsInfluencers.Close(ctx)
		// normalize db results
		for resultsInfluencers.Next(ctx) {
			var singleInfluencer models.InfluencerSmallDataModel
			if err = resultsInfluencers.Decode(&singleInfluencer); err != nil {
				return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"error": err.Error()}})
			}

			influencers = append(influencers, singleInfluencer)
		}

		gallery.Influencers = nil
		gallery.InfluencersData = influencers
	}
	// end of get all influencers on post

	// get author information if author_id exists
	if gallery.AuthorID != "" {
		fmt.Printf("DEBUG: DetailGallery - Found AuthorID: %s\n", gallery.AuthorID)
		authorObjID, err := primitive.ObjectIDFromHex(gallery.AuthorID)
		if err == nil {
			var author models.UserModel
			err := galleryUsersCollection.FindOne(ctx, bson.M{"_id": authorObjID}).Decode(&author)
			if err == nil {
				gallery.Author = &models.AuthorModel{
					ID:       author.ID.Hex(),
					Username: author.Username,
				}
				fmt.Printf("DEBUG: DetailGallery - Author found: %s\n", author.Username)
			} else {
				fmt.Printf("DEBUG: DetailGallery - Error finding author: %v\n", err)
			}
		} else {
			fmt.Printf("DEBUG: DetailGallery - Error parsing AuthorID: %v\n", err)
		}
	} else {
		fmt.Printf("DEBUG: DetailGallery - No AuthorID found for gallery: %s\n", gallery.Id.Hex())
	}

	return c.JSON(http.StatusOK, responses.GlobalResponse{Status: http.StatusOK, Message: "OK", Data: &echo.Map{"gallery": gallery}})
}

// handle of POST /galleries
func CreateGallery(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var payload models.PayloadGallery
	err := json.NewDecoder(c.Request().Body).Decode(&payload)

	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{Status: http.StatusBadRequest, Message: "Error parsing json", Data: nil})
	} else {

		// ref: https://stackoverflow.com/a/8689281/2780875
		slug := strings.Replace(payload.Title, " ", "-", -1)
		slug = strings.ToLower(slug)

		// insert data to db
		result, errInsertGallery := repositories.CreateGallery(ctx, repositories.CreateGalleryParams{
			Title:       payload.Title,
			Description: payload.Description,
			Images:      payload.Images,
			Influencers: payload.Influencers,
			Lang:        payload.Lang,
			Slug:        slug,
		})

		if errInsertGallery != nil {
			return c.JSON(http.StatusBadRequest, responses.GlobalResponse{Status: http.StatusBadRequest, Message: "Error insert data", Data: nil})
		} else {
			// post gallery to telegram channel

			chatMessage := "New Gallery:\n" + payload.Title +
				"\nhttps://follooow.com/" + payload.Lang + "/gallery/" + slug + "-" + result.InsertedID.(primitive.ObjectID).Hex()
			repositories.TelegramSendMessage(chatMessage)
			// end of gallery news to telegram channel
			return c.JSON(http.StatusCreated, responses.GlobalResponse{Status: http.StatusCreated, Message: "Success create gallery", Data: nil})
		}
	}
}

// handle of POST /galleries/upload - for creating gallery with image uploads
func CreateGalleryWithUpload(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{Status: http.StatusBadRequest, Message: "Error parsing multipart form", Data: &echo.Map{"error": err.Error()}})
	}

	// Get form fields
	title := c.FormValue("title")
	description := c.FormValue("description")
	lang := c.FormValue("lang")
	influencersStr := c.FormValue("influencers")
	authorID := c.FormValue("author_id")

	// Validate required fields
	if title == "" {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{Status: http.StatusBadRequest, Message: "Title is required", Data: nil})
	}
	if lang == "" {
		lang = "ID" // default language
	}

	// Parse influencers
	var influencers []string
	if influencersStr != "" {
		influencers = strings.Split(influencersStr, ",")
		// Trim whitespace from each influencer ID
		for i, inf := range influencers {
			influencers[i] = strings.TrimSpace(inf)
		}
	}

	// Get uploaded files
	files := form.File["images"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{Status: http.StatusBadRequest, Message: "At least one image is required", Data: nil})
	}

	// Initialize Cloudinary if not already done
	if configs.CloudinaryClient == nil {
		configs.InitCloudinary()
	}

	// Upload images to Cloudinary
	var images []models.ImageModel
	for i, file := range files {
		// Upload image
		result, err := utils.UploadImageFromForm(ctx, file, "galleries")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "Error uploading image", Data: &echo.Map{"error": err.Error()}})
		}

		// Create image model
		imageModel := models.ImageModel{
			IsCover:   i == 0, // First image is cover
			Url:       result.SecureURL,
			Caption:   file.Filename,
			CreatedOn: int(time.Now().Unix()),
			UpdatedOn: int(time.Now().Unix()),
		}

		images = append(images, imageModel)
	}

	// Generate slug
	slug := strings.Replace(title, " ", "-", -1)
	slug = strings.ToLower(slug)

	// Insert gallery to database
	result, err := repositories.CreateGallery(ctx, repositories.CreateGalleryParams{
		Title:       title,
		Description: description,
		Images:      images,
		Influencers: influencers,
		Lang:        lang,
		Slug:        slug,
		AuthorID:    authorID,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.GlobalResponse{Status: http.StatusBadRequest, Message: "Error creating gallery", Data: &echo.Map{"error": err.Error()}})
	}

	// Post gallery to telegram channel
	chatMessage := "New Gallery:\n" + title +
		"\nhttps://follooow.com/" + lang + "/gallery/" + slug + "-" + result.InsertedID.(primitive.ObjectID).Hex()
	repositories.TelegramSendMessage(chatMessage)

	return c.JSON(http.StatusCreated, responses.GlobalResponse{Status: http.StatusCreated, Message: "Success create gallery with images", Data: &echo.Map{"gallery_id": result.InsertedID}})
}

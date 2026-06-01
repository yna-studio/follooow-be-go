package handlers

import (
    "context"
    "time"
    "follooow-be/configs"
    "follooow-be/models"
    // "follooow-be/repositories"
    "follooow-be/responses"
    "net/http"

    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// SyncLabels aggregates label usage across influencers, news, and galleries and updates the labels collection.
// Endpoint: GET /sync/labels
func SyncLabels(c echo.Context) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Prepare collections
    influencersCol := configs.GetCollection(configs.DB, "influencers")
    newsCol := configs.GetCollection(configs.DB, "news")
    galleriesCol := configs.GetCollection(configs.DB, "galleries")

    // Map to hold aggregated counts
    labelCounts := make(map[string]int)

    // Helper to process a cursor and extract string slice from a field
    processCursor := func(cur *mongo.Cursor, field string) error {
        defer cur.Close(ctx)
        for cur.Next(ctx) {
            var doc bson.M
            if err := cur.Decode(&doc); err != nil {
                return err
            }
            if val, ok := doc[field]; ok {
                // Mongo driver returns arrays as bson.A (alias of []interface{})
                if arr, ok := val.(bson.A); ok {
                    for _, v := range arr {
                        if s, ok := v.(string); ok && s != "" {
                            labelCounts[s]++
                        }
                    }
                } else if arr, ok := val.([]interface{}); ok { // fallback
                    for _, v := range arr {
                        if s, ok := v.(string); ok && s != "" {
                            labelCounts[s]++
                        }
                    }
                }
            }
        }
        return cur.Err()
    }

    // Influencers: field "label"
    curInf, err := influencersCol.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"label": 1}))
    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error fetching influencers", Data: &echo.Map{"error": err.Error()}})
    }
    if err := processCursor(curInf, "label"); err != nil {
        return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error processing influencers", Data: &echo.Map{"error": err.Error()}})
    }

    // News: field "tags"
    curNews, err := newsCol.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"tags": 1}))
    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error fetching news", Data: &echo.Map{"error": err.Error()}})
    }
    if err := processCursor(curNews, "tags"); err != nil {
        return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error processing news", Data: &echo.Map{"error": err.Error()}})
    }

    // Galleries: field "tags"
    curGal, err := galleriesCol.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"tags": 1}))
    if err != nil {
        return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error fetching galleries", Data: &echo.Map{"error": err.Error()}})
    }
    if err := processCursor(curGal, "tags"); err != nil {
        return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error processing galleries", Data: &echo.Map{"error": err.Error()}})
    }

    // Replace the entire labels collection with the newly computed totals
    labelCol := configs.GetCollection(configs.DB, "labels")
    // Remove all existing label documents
    if _, err := labelCol.DeleteMany(ctx, bson.M{}); err != nil {
        return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error clearing labels", Data: &echo.Map{"error": err.Error()}})
    }
    // Prepare documents for insertion
    var docs []interface{}
    for lbl, cnt := range labelCounts {
        docs = append(docs, models.LabelModel{Label: lbl, Total: int64(cnt)})
    }
    if len(docs) > 0 {
        if _, err := labelCol.InsertMany(ctx, docs); err != nil {
            return c.JSON(http.StatusInternalServerError, responses.GlobalResponse{Status: http.StatusInternalServerError, Message: "error inserting labels", Data: &echo.Map{"error": err.Error()}})
        }
    }
    return c.JSON(http.StatusOK, responses.GlobalResponse{Status: http.StatusOK, Message: "labels synchronized", Data: nil})
}

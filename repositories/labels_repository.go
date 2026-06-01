package repositories

import (
    "context"
    "follooow-be/configs"
    "follooow-be/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var labelCollection *mongo.Collection = configs.GetCollection(configs.DB, "labels")

// ListLabels returns all label documents, sorted by total count descending, limited by the provided amount.
func ListLabels(ctx context.Context, limit int64) ([]models.LabelModel, error) {
    opts := options.Find().SetSort(bson.D{{"total", -1}})
    if limit > 0 {
        opts.SetLimit(limit)
    }
    cursor, err := labelCollection.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    var results []models.LabelModel
    for cursor.Next(ctx) {
        var lbl models.LabelModel
        if err = cursor.Decode(&lbl); err != nil {
            return nil, err
        }
        results = append(results, lbl)
    }
    return results, nil
}

// UpdateLabelCount increments (inc can be negative) the total count for a label. If the label does not exist and inc > 0, it will be created (upsert).
func UpdateLabelCount(ctx context.Context, label string, inc int) error {
    filter := bson.M{"label": label}
    update := bson.M{"$inc": bson.M{"total": inc}}
    opts := options.Update().SetUpsert(true)
    _, err := labelCollection.UpdateOne(ctx, filter, update, opts)
    return err
}


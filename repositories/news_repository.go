package repositories

import (
	"context"
	"follooow-be/configs"
	"follooow-be/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// struct of GetDetailNews() params
type DetailNewsParams struct {
	NewsId string
	Lang   string
}

var NewsCollections *mongo.Collection = configs.GetCollection(configs.DB, "news")
var UsersCollections *mongo.Collection = configs.GetCollection(configs.DB, "users")
var NewsInfluencersCollections *mongo.Collection = configs.GetCollection(configs.DB, "influencers")

// function to to get detail by news_id
// auto increase visits + 1 if data found on DB
func GetDetailNews(ctx context.Context, params DetailNewsParams) (error, models.NewsModel) {
	var news models.NewsModel
	var influencers []models.InfluencerSmallDataModel

	objId, _ := primitive.ObjectIDFromHex(params.NewsId)
	filterListData := bson.M{
		"_id": objId,
	}

	// Only add lang filter if it exists in the document
	// Some news documents might not have lang field
	if params.Lang != "" {
		filterListData["lang"] = params.Lang
	}

	err := NewsCollections.FindOne(ctx, filterListData).Decode(&news)

	if err != nil {
		// Return the error (could be "no documents in result" or other errors)
		return err, news
	}

	// Document found successfully, continue processing
	// increase visits
	_, updateErr := NewsCollections.UpdateOne(ctx, bson.D{{"_id", objId}}, bson.D{{"$set", bson.D{{"views", news.Views + 1}}}})

	if updateErr != nil {
		return updateErr, news
	}

	// get influencer data if news has influencers
	if len(news.Influencers) > 0 {
		var idsArr []string = news.Influencers
		var idsObjId []primitive.ObjectID

		// normalize ids
		for key := range idsArr {
			objId, _ := primitive.ObjectIDFromHex(idsArr[key])
			idsObjId = append(idsObjId, objId)
		}
		optsListDataInfluencers := options.Find().SetLimit(20)

		// filter generator
		filterLastDataInfluencers := bson.D{{"_id", bson.M{"$in": idsObjId}}}

		// get influencer data from database
		resultsInfluencers, _ := NewsInfluencersCollections.Find(ctx, filterLastDataInfluencers, optsListDataInfluencers)
		defer resultsInfluencers.Close(ctx)

		// normalize db results
		for resultsInfluencers.Next(ctx) {
			var singleInfluencer models.InfluencerSmallDataModel
			if err = resultsInfluencers.Decode(&singleInfluencer); err != nil {
				return err, news
			}

			influencers = append(influencers, singleInfluencer)
		}

		news.Influencers = nil
		news.InfluencersData = influencers
	}

	// get author information if author_id exists
	if news.AuthorID != "" {
		authorObjID, err := primitive.ObjectIDFromHex(news.AuthorID)
		if err == nil {
			var author models.UserModel
			err := UsersCollections.FindOne(ctx, bson.M{"_id": authorObjID}).Decode(&author)
			if err == nil {
				news.Author = &models.AuthorModel{
					ID:       author.ID.Hex(),
					Username: author.Username,
				}
			}
		}
	}

	return nil, news

}

package model

import (
	"context"

	"github.com/kartikayshankh/food-delivery/utils"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DBNAME = "food-delivery"
)

func Register(c echo.Context, request interface{}, collection *mongo.Collection) *utils.ErrorHandler {
	ctx := c.Request().Context()
	_, err := collection.InsertOne(ctx, &request)
	if err != nil {
		// Check if the error is a duplicate key error
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				if e.Code == 11000 {
					return &utils.ErrorHandler{
						Message:    "already exists",
						DevMessage: e.Message,
						Code:       404,
					}
				}
			}
		}
		return &utils.ErrorHandler{
			Message:    utils.SOMETHING_WENT_WRONG,
			DevMessage: err.Error(),
			Code:       400,
		}
	}
	return nil
}

func UpdateRiderLocation(c echo.Context, location Location, collection *mongo.Collection, filter primitive.M) *utils.ErrorHandler {
	update := bson.M{"$set": bson.M{"location": location}}

	result := collection.FindOneAndUpdate(c.Request().Context(), filter, update)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return &utils.ErrorHandler{
				Message:    utils.DATA_NOT_FOUND,
				DevMessage: result.Err().Error(),
				Code:       404,
			}
		}
		return &utils.ErrorHandler{
			Message:    utils.SOMETHING_WENT_WRONG,
			DevMessage: result.Err().Error(),
			Code:       400,
		}
	}
	return nil
}

func GetOrderHistory(c echo.Context, filter primitive.M, collection *mongo.Collection) (*[]Order, *utils.ErrorHandler) {

	Order := new([]Order)
	cursor, err := collection.Find(c.Request().Context(), filter)
	if err != nil {
		return nil, &utils.ErrorHandler{
			Message:    utils.SOMETHING_WENT_WRONG,
			DevMessage: err.Error(),
			Code:       400}
	}
	defer cursor.Close(c.Request().Context())
	errCursor := cursor.All(c.Request().Context(), Order)
	if errCursor != nil {
		return nil, &utils.ErrorHandler{
			Message:    utils.SOMETHING_WENT_WRONG,
			DevMessage: errCursor.Error(),
			Code:       400}
	}

	if len(*Order) == 0 {
		return nil, &utils.ErrorHandler{
			Message: utils.DATA_NOT_FOUND,
			Code:    404,
		}
	}
	return Order, nil
}

func SuggestRestaurant(c echo.Context, preferences UserPreferences, collection *mongo.Collection) ([]Restaurant, *utils.ErrorHandler) {
	var restaurants []Restaurant
	// need to pagination
	filter := bson.M{
		"cuisine": bson.M{"$in": preferences.Cuisines},
		"rating":  bson.M{"$gte": preferences.MinRating},
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{preferences.Location.Longitude, preferences.Location.Latitude},
				},
				"$maxDistance": preferences.MaxDistance,
			},
		},
	}
	cursor, err := collection.Find(c.Request().Context(), filter, options.Find().SetLimit(50))
	if err != nil {
		return nil, &utils.ErrorHandler{
			Message:    utils.SOMETHING_WENT_WRONG,
			DevMessage: err.Error(),
			Code:       400,
		}
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(c.Request().Context(), &restaurants); err != nil {
		return nil, &utils.ErrorHandler{
			Message:    utils.SOMETHING_WENT_WRONG,
			DevMessage: err.Error(),
			Code:       400,
		}
	}
	return restaurants, nil
}

func GetRestaurant(c echo.Context, restaurantId string, collection *mongo.Collection) (*Restaurant, *utils.ErrorHandler) {
	filter := bson.M{"_id": restaurantId}
	var restaurant Restaurant
	err := collection.FindOne(c.Request().Context(), filter).Decode(&restaurant)
	if err != nil {
		return nil, &utils.ErrorHandler{
			Message:    utils.SOMETHING_WENT_WRONG,
			DevMessage: err.Error(),
			Code:       400,
		}
	}
	return &restaurant, nil
}

func GetUser(c echo.Context, filter primitive.M, collection *mongo.Collection) (*User, *utils.ErrorHandler) {
	user := new(User)
	err := collection.FindOne(c.Request().Context(), filter).Decode(user)
	if err != nil {
		return nil, &utils.ErrorHandler{
			Message:    utils.SOMETHING_WENT_WRONG,
			DevMessage: err.Error(),
			Code:       400,
		}
	}
	return user, nil
}

func AcceptOrder(c echo.Context, orderDetails Order, collection *mongo.Collection) *utils.ErrorHandler {
	_, err := collection.InsertOne(c.Request().Context(), orderDetails)
	if err != nil {
		return &utils.ErrorHandler{
			Message:    utils.SOMETHING_WENT_WRONG,
			DevMessage: err.Error(),
			Code:       400,
		}
	}
	return nil
}

func GetUserOrderHistory(c echo.Context, filter primitive.M, collection *mongo.Collection) ([]Order, *utils.ErrorHandler) {
	result := new([]Order)
	cursor, err := collection.Find(c.Request().Context(), filter, options.Find().SetLimit(10))
	if err != nil {
		return nil, &utils.ErrorHandler{Message: utils.SOMETHING_WENT_WRONG, DevMessage: err.Error(), Code: 400}
	}
	defer cursor.Close(c.Request().Context())
	errCursor := cursor.All(c.Request().Context(), result)
	if errCursor != nil {
		return nil, &utils.ErrorHandler{
			Message:    utils.SOMETHING_WENT_WRONG,
			DevMessage: errCursor.Error(),
			Code:       400,
		}
	}
	if len(*result) == 0 {
		return nil, &utils.ErrorHandler{Message: utils.DATA_NOT_FOUND,
			Code: 404}
	}
	return *result, nil
}

func GetRiderLocation(c echo.Context, location *Location, collection *mongo.Collection) ([]Rider, *utils.ErrorHandler) {
	var rider []Rider
	// need to pagination
	pipeline := mongo.Pipeline{
		{{Key: "$geoNear", Value: bson.D{
			{Key: "near", Value: bson.D{
				{Key: "type", Value: "Point"},
				{Key: "coordinates", Value: bson.A{location.Latitude, location.Longitude}},
			}},
			{Key: "distanceField", Value: "dist.calculated"},
			{Key: "maxDistance", Value: 20000},
			{Key: "spherical", Value: true},
		}}},
	}
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, &utils.ErrorHandler{
			Message:    utils.SOMETHING_WENT_WRONG,
			DevMessage: err.Error(),
			Code:       400,
		}
	}
	defer cursor.Close(context.Background())

	// Check if cursor has any documents
	if !cursor.Next(context.Background()) {
		return nil, &utils.ErrorHandler{
			Message:    utils.DATA_NOT_FOUND,
			DevMessage: utils.NO_RIDER_FOUND_NEAR_LOCATION,
			Code:       404,
		}
	}

	// Decode results into the provided interface (slice)
	if err := cursor.All(c.Request().Context(), &rider); err != nil {
		return nil, &utils.ErrorHandler{
			Message:    utils.SOMETHING_WENT_WRONG,
			DevMessage: err.Error(),
			Code:       400,
		}
	}
	return rider, nil
}

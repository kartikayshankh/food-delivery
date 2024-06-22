package model

import (
	"assignment/utils"
	"context"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(c echo.Context, request interface{}, Client mongo.Client) *utils.ErrorHandler {

	collection := Client.Database("test").Collection("users")
	_, err := collection.InsertOne(context.Background(), &request)
	if err != nil {
		return &utils.ErrorHandler{Message: utils.SOMETHING_WENT_WRONG, DevMessage: err.Error()}
	}

	return nil

}

func UpdateRiderLocation(c echo.Context, location Location, Client mongo.Client, filter primitive.M) *utils.ErrorHandler {

	collection := Client.Database("test").Collection("riders")

	update := bson.M{"$set": bson.M{"location": location}}

	result := collection.FindOneAndUpdate(c.Request().Context(), filter, update)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return &utils.ErrorHandler{Message: utils.DATA_NOT_FOUND, DevMessage: result.Err().Error()}
		}
		return &utils.ErrorHandler{Message: utils.SOMETHING_WENT_WRONG, DevMessage: result.Err().Error()}
	}

	return nil
}

func GetRiderOrderHistory(c echo.Context, filter primitive.M, Client mongo.Client) (*[]Rider, *utils.ErrorHandler) {
	result := new([]Rider)
	collection := Client.Database("food-delivery").Collection("riders")
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, &utils.ErrorHandler{Message: utils.SOMETHING_WENT_WRONG, DevMessage: err.Error(), Code: 400}
	}
	defer cursor.Close(c.Request().Context())

	errCursor := cursor.All(context.Background(), result)
	if errCursor != nil {
		return nil, &utils.ErrorHandler{Message: utils.SOMETHING_WENT_WRONG, DevMessage: errCursor.Error(), Code: 400}
	}

	if len(*result) == 0 {
		return nil, &utils.ErrorHandler{Message: utils.DATA_NOT_FOUND, Code: 404}
	}

	return result, nil

}

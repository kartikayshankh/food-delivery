package model

import (
	"assignment/utils"
	"context"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(c echo.Context, request interface{}, Client mongo.Client, role string) *utils.ErrorHandler {
	ctx := c.Request().Context()
	collectionName := ""
	switch role {
	case string(utils.User):
		collectionName = string(utils.User)
	case string(utils.Rider):
		collectionName = string(utils.Rider)
	}
	collection := Client.Database("user").Collection(collectionName)
	_, err := collection.InsertOne(ctx, &request)
	if err != nil {
		// Check if the error is a duplicate key error
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				if e.Code == 11000 {
					return &utils.ErrorHandler{Message: collectionName + " already exists", DevMessage: e.Message}
				}
			}
		}
		return &utils.ErrorHandler{Message: "Something went wrong", DevMessage: err.Error()}
	}
	return nil
}

func UpdateRiderLocation(c echo.Context, location Location, Client mongo.Client, filter primitive.M) *utils.ErrorHandler {

	collection := Client.Database("user").Collection("riders")

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

package user

import (
	"assignment/model"
	"assignment/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	Register(c echo.Context, request *model.User) *utils.ErrorHandler
	GetUserOrderHistory(c echo.Context, riderId string) (*[]model.Order, *utils.ErrorHandler)
}

type userService struct {
	config      *viper.Viper
	mongoClient *mongo.Client
}

func NewUserService(config *viper.Viper, mongoClient *mongo.Client) UserService {
	return &userService{config: config, mongoClient: mongoClient}
}

func (r *userService) getCollection() *mongo.Collection {
	dbName := r.config.GetString("mongodb.database")
	collectionName := r.config.GetString("mongodb.collection.user")
	collection := r.mongoClient.Database(dbName).Collection(collectionName)
	return collection
}

func (u *userService) Register(c echo.Context, request *model.User) *utils.ErrorHandler {
	//we will verify user's email with otp (assuming it is otp already there)
	// dbName := u.config.GetString("mongodb.database")
	// collectionName := u.config.GetString("mongodb.collection.user")
	request.ID = uuid.NewString()
	collection := u.getCollection()
	err := model.Register(c, request, collection)
	if err != nil {
		return &utils.ErrorHandler{
			Message:    err.Message,
			DevMessage: err.DevMessage,
			Code:       err.Code,
		}
	}
	return nil
}

func (r *userService) GetUserOrderHistory(c echo.Context, userId string) (*[]model.Order, *utils.ErrorHandler) {
	dbName := r.config.GetString("mongodb.database")
	collectionName := r.config.GetString("mongodb.collection.orders")
	collection := r.mongoClient.Database(dbName).Collection(collectionName)
	// collection := r.getCollection()
	filter := bson.M{"user_id": userId}
	response, err := model.GetOrderHistory(c, filter, collection)
	if err != nil {
		return nil, &utils.ErrorHandler{
			Message:    err.Message,
			DevMessage: err.DevMessage,
			Code:       err.Code}
	}
	return response, nil
}

package rider

import (
	"assignment/model"
	"assignment/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RiderService interface {
	ResgisterRider(c echo.Context, request *model.Rider) *utils.ErrorHandler
	UpdateRiderLocation(c echo.Context, Location model.Location, riderId string) *utils.ErrorHandler
	GetRiderOrderHistory(c echo.Context, riderId string) (*[]model.Order, *utils.ErrorHandler)
}

type riderService struct {
	config      *viper.Viper
	mongoClient *mongo.Client
}

func NewRiderService(config *viper.Viper, mongoClient *mongo.Client) RiderService {
	return &riderService{config: config, mongoClient: mongoClient}
}

func (r *riderService) getCollection() *mongo.Collection {
	dbName := r.config.GetString("mongodb.database")
	collectionName := r.config.GetString("mongodb.collection.rider")
	collection := r.mongoClient.Database(dbName).Collection(collectionName)
	return collection
}

func (r *riderService) ResgisterRider(c echo.Context, request *model.Rider) *utils.ErrorHandler {
	// dbName := r.config.GetString("mongodb.database")
	// collectionName := r.config.GetString("mongodb.collection.rider")
	request.ID = uuid.NewString()
	collection := r.getCollection()
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

func (r *riderService) UpdateRiderLocation(c echo.Context, Location model.Location, riderId string) *utils.ErrorHandler {
	// dbName := r.config.GetString("mongodb.database")
	// collectionName := r.config.GetString("mongodb.collection.rider")
	collection := r.getCollection()

	filter := bson.M{"_id": riderId}
	err := model.UpdateRiderLocation(c, Location, collection, filter)

	if err != nil {
		return &utils.ErrorHandler{
			Message:    err.Message,
			DevMessage: err.DevMessage,
			Code:       err.Code,
		}
	}

	return nil
}

func (r *riderService) GetRiderOrderHistory(c echo.Context, riderId string) (*[]model.Order, *utils.ErrorHandler) {
	dbName := r.config.GetString("mongodb.database")
	collectionName := r.config.GetString("mongodb.collection.orders")
	collection := r.mongoClient.Database(dbName).Collection(collectionName)
	// collection := r.getCollection()
	filter := bson.M{"rider_id": riderId}
	response, err := model.GetOrderHistory(c, filter, collection)
	if err != nil {
		return nil, &utils.ErrorHandler{
			Message:    err.Message,
			DevMessage: err.DevMessage,
			Code:       err.Code,
		}
	}
	return response, nil
}

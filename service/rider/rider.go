package rider

import (
	"assignment/model"
	"assignment/utils"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RiderService interface {
	ResgisterRider(c echo.Context, request *model.Rider) *utils.ErrorHandler
	UpdateRiderLocation(c echo.Context, Location model.Location, riderId string) *utils.ErrorHandler
	GetRiderOrderHistory(c echo.Context, riderId string) ([]model.Rider, *utils.ErrorHandler)
}

type riderService struct {
	config      *viper.Viper
	mongoClient *mongo.Client
}

func NewRiderService(config *viper.Viper, mongoClient *mongo.Client) RiderService {
	return &riderService{config: config, mongoClient: mongoClient}
}

func (r *riderService) ResgisterRider(c echo.Context, request *model.Rider) *utils.ErrorHandler {
	uriParts := strings.Split(c.Request().RequestURI, "/")
	role := ""
	if len(uriParts) > 1 {
		role = uriParts[1]
	}
	err := model.Register(c, request, *r.mongoClient, role)
	if err != nil {
		return &utils.ErrorHandler{Message: err.Message, DevMessage: err.DevMessage}
	}
	return nil
}

func (r *riderService) UpdateRiderLocation(c echo.Context, Location model.Location, riderId string) *utils.ErrorHandler {
	filter := bson.M{"_id": riderId}

	err := model.UpdateRiderLocation(c, Location, *r.mongoClient, filter)

	if err != nil {
		return &utils.ErrorHandler{Message: err.Message, DevMessage: err.DevMessage}
	}

	return nil
}

func (r *riderService) GetRiderOrderHistory(c echo.Context, riderId string) ([]model.Rider, *utils.ErrorHandler) {
	filter := bson.M{"_id": riderId}

	response, err := model.GetRiderOrderHistory(c, filter, *r.mongoClient)
	if err != nil {
		return nil, &utils.ErrorHandler{Message: err.Message, DevMessage: err.DevMessage, Code: err.Code}
	}

	return *response, nil
}

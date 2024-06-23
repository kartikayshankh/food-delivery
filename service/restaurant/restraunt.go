package restaurant

import (
	"assignment/model"
	"assignment/utils"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type RestaurantService interface {
	Resgister(c echo.Context, request *model.Restaurant) *utils.ErrorHandler
	SuggestRestaurant(c echo.Context, request model.UserPreferences) ([]model.Restaurant, *utils.ErrorHandler)
	GetRestaurantMenu(c echo.Context, restaurantId string) (*model.Restaurant, *utils.ErrorHandler)
	AcceptOrder(c echo.Context, orderDetails model.Order) *utils.ErrorHandler
	GetRiderLocation(c echo.Context, location *model.Location) (*[]model.Rider, *utils.ErrorHandler)
	NearestRider(c echo.Context, id string) (*model.Rider, *utils.ErrorHandler)
}

type restaurantService struct {
	config      *viper.Viper
	mongoClient *mongo.Client
}

func NewRestaurantService(config *viper.Viper, mongoClient *mongo.Client) RestaurantService {
	return &restaurantService{config: config, mongoClient: mongoClient}
}

func (u *restaurantService) getCollection() *mongo.Collection {
	dbName := u.config.GetString("mongodb.database")
	collectionName := u.config.GetString("mongodb.collection.restaurant")
	collection := u.mongoClient.Database(dbName).Collection(collectionName)
	return collection
}

func (u *restaurantService) Resgister(c echo.Context, request *model.Restaurant) *utils.ErrorHandler {
	// dbName := u.config.GetString("mongodb.database")
	// collectionName := u.config.GetString("mongodb.collection.restaurant")
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

func (u *restaurantService) SuggestRestaurant(c echo.Context, request model.UserPreferences) ([]model.Restaurant, *utils.ErrorHandler) {
	collection := u.getCollection()
	response, err := model.SuggestRestaurant(c, request, collection)
	if err != nil {
		return nil, &utils.ErrorHandler{Message: err.Message, DevMessage: err.DevMessage}
	}

	return response, nil

}

func (u *restaurantService) GetRestaurantMenu(c echo.Context, restaurantId string) (*model.Restaurant, *utils.ErrorHandler) {
	collection := u.getCollection()
	response, err := model.GetRestaurant(c, restaurantId, collection)
	if err != nil {
		return nil, &utils.ErrorHandler{Message: err.Message, DevMessage: err.DevMessage}
	}
	return response, nil
}

func (u *restaurantService) AcceptOrder(c echo.Context, orderDetails model.Order) *utils.ErrorHandler {
	dbName := u.config.GetString("mongodb.database")
	userCollectionName := u.config.GetString("mongodb.collection.user")
	userCollection := u.mongoClient.Database(dbName).Collection(userCollectionName)
	filter := bson.M{"_id": orderDetails.UserID}
	user, err := model.GetUser(c, filter, userCollection)
	if err != nil {
		return &utils.ErrorHandler{
			Message:    err.Message,
			DevMessage: err.DevMessage,
			Code:       err.Code,
		}
	}

	//restaurant
	collection := u.getCollection()
	restaurant, err := model.GetRestaurant(c, orderDetails.RestaurantID, collection)
	if err != nil {
		return &utils.ErrorHandler{
			Message:    err.Message,
			DevMessage: err.DevMessage,
			Code:       err.Code,
		}
	}
	prices := make(map[string]float64)
	for _, v := range restaurant.Menu {
		prices[v.ID] = v.Price
	}
	order := new(model.Order)
	totalPrices := 0.00
	for _, value := range orderDetails.Items {
		totalPrices += prices[value.MenuItemID]
	}
	order.ID = uuid.NewString()
	order.Items = orderDetails.Items
	order.UserID = user.ID

	//need to update status once the order is completed
	order.Status = "Pending"
	order.TotalPrice = totalPrices
	order.Createdat = time.Now()
	order.Updatedat = time.Now()
	order.RestaurantID = orderDetails.RestaurantID
	OrderCollectionName := u.config.GetString("mongodb.collection.orders")
	orderCollection := u.mongoClient.Database(dbName).Collection(OrderCollectionName)

	err = model.AcceptOrder(c, *order, orderCollection)
	if err != nil {
		return &utils.ErrorHandler{
			Message:    err.Message,
			DevMessage: err.DevMessage,
			Code:       err.Code,
		}
	}
	return nil
}

func (u *restaurantService) GetRiderLocation(c echo.Context, location *model.Location) (*[]model.Rider, *utils.ErrorHandler) {

	dbName := u.config.GetString("mongodb.database")
	ridercollectionName := u.config.GetString("mongodb.collection.rider")
	riderCollection := u.mongoClient.Database(dbName).Collection(ridercollectionName)
	rider, err := model.GetRiderLocation(c, location, riderCollection)
	if err != nil {
		return nil, &utils.ErrorHandler{
			Message:    err.Message,
			DevMessage: err.DevMessage,
			Code:       err.Code,
		}
	}
	return &rider, nil
}

func (u *restaurantService) NearestRider(c echo.Context, id string) (*model.Rider, *utils.ErrorHandler) {
	collection := u.getCollection()
	restaurant, err := model.GetRestaurant(c, id, collection)
	if err != nil {
		return nil, &utils.ErrorHandler{
			Message:    err.Message,
			DevMessage: err.DevMessage,
		}
	}
	location := new(model.Location)
	dbName := u.config.GetString("mongodb.database")
	ridercollectionName := u.config.GetString("mongodb.collection.rider")
	riderCollection := u.mongoClient.Database(dbName).Collection(ridercollectionName)

	location.Latitude = restaurant.Location.Latitude
	location.Longitude = restaurant.Location.Longitude
	rider, err := model.GetRiderLocation(c, location, riderCollection)
	if err != nil {
		return nil, &utils.ErrorHandler{
			Message:    err.Message,
			DevMessage: err.DevMessage,
			Code:       err.Code,
		}
	}
	return &rider[0], nil
}

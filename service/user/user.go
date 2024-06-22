package user

import (
	"assignment/model"
	"assignment/utils"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	Register(c echo.Context, request *model.User) *utils.ErrorHandler
}

type userService struct {
	config      *viper.Viper
	mongoClient *mongo.Client
}

func NewUserService(config *viper.Viper, mongoClient *mongo.Client) UserService {
	return &userService{config: config, mongoClient: mongoClient}
}

func (u *userService) Register(c echo.Context, request *model.User) *utils.ErrorHandler {
	//we will verify user's email with otp (assuming it is already there)
	err := model.Register(c, request, *u.mongoClient)
	if err != nil {
		return &utils.ErrorHandler{Message: err.Message, DevMessage: err.DevMessage}
	}
	return nil
}


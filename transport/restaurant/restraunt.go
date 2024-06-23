package restaurant

import (
	"assignment/model"
	restaurantService "assignment/service/restaurant"
	"assignment/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type restaurantEndpoint struct {
	restaurantService restaurantService.RestaurantService
}

type RestaurantEnpoint interface {
	Register(c echo.Context) error
	SuggestRestaurant(c echo.Context) error
	GetRestaurantMenu(c echo.Context) error
	AcceptOrder(c echo.Context) error
	GetRiderLocation(c echo.Context) error
	NearestRider(c echo.Context) error
}

func NewRestaurantEndpoint(restaurantService restaurantService.RestaurantService) RestaurantEnpoint {
	return &restaurantEndpoint{restaurantService: restaurantService}
}

func (r *restaurantEndpoint) Register(c echo.Context) error {
	restaurant := new(model.Restaurant)

	err := c.Bind(restaurant)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.GenericResponse{
			Message:          utils.SOMETHING_WENT_WRONG,
			DeveloperMessage: utils.DATA_INVALID,
		})
	}

	validationError := utils.Validator(restaurant)
	if validationError != nil {
		return c.JSON(http.StatusBadRequest, &utils.GenericResponse{
			Message:          utils.DATA_INVALID,
			DeveloperMessage: validationError.Message,
		})
	}

	errRegister := r.restaurantService.Resgister(c, restaurant)
	if errRegister != nil {
		return c.JSON(errRegister.Code, utils.GenericResponse{
			Message: errRegister.Message,
		})
	}

	return c.JSON(http.StatusCreated, utils.GenericResponse{
		Message: utils.RESTAURANT_CREATED_SUCCESSFULLY,
	})

}

func (u *restaurantEndpoint) SuggestRestaurant(c echo.Context) error {
	preferences := new(model.UserPreferences)
	err := c.Bind(preferences)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.GenericResponse{
			Message:          utils.SOMETHING_WENT_WRONG,
			DeveloperMessage: utils.DATA_INVALID,
		})
	}
	validationError := utils.Validator(preferences)
	if validationError != nil {
		return c.JSON(http.StatusBadRequest, &utils.GenericResponse{
			Message:          utils.DATA_INVALID,
			DeveloperMessage: validationError.Message,
		})
	}
	response, errResponse := u.restaurantService.SuggestRestaurant(c, *preferences)
	if errResponse != nil {
		return c.JSON(errResponse.Code, utils.GenericResponse{
			Message: errResponse.Message,
		})
	}
	return c.JSON(http.StatusCreated, response)
}

func (r *restaurantEndpoint) GetRestaurantMenu(c echo.Context) error {
	restaurantId := c.Param("id")
	response, errResponse := r.restaurantService.GetRestaurantMenu(c, restaurantId)
	if errResponse != nil {
		return c.JSON(errResponse.Code, utils.GenericResponse{
			Message: utils.DATA_NOT_FOUND,
		})
	}
	return c.JSON(http.StatusOK, response.Menu)
}

func (r *restaurantEndpoint) AcceptOrder(c echo.Context) error {
	order := new(model.Order)
	err := c.Bind(order)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.GenericResponse{
			Message:          utils.SOMETHING_WENT_WRONG,
			DeveloperMessage: utils.DATA_INVALID,
		})
	}

	validationError := utils.Validator(order)
	if validationError != nil {
		return c.JSON(http.StatusBadRequest, &utils.GenericResponse{
			Message:          utils.DATA_INVALID,
			DeveloperMessage: validationError.Message,
		})
	}

	errAcceptOrder := r.restaurantService.AcceptOrder(c, *order)
	if errAcceptOrder != nil {
		return c.JSON(errAcceptOrder.Code, utils.GenericResponse{
			Message: errAcceptOrder.Message,
		})
	}

	return c.JSON(http.StatusCreated, utils.GenericResponse{
		Message: utils.ORDER_ACCEPTED_SUCCESSFULLY,
	})
}

func (r *restaurantEndpoint) GetRiderLocation(c echo.Context) error {
	location := new(model.Location)
	err := c.Bind(location)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.GenericResponse{
			Message:          utils.SOMETHING_WENT_WRONG,
			DeveloperMessage: utils.DATA_INVALID,
		})
	}
	validationError := utils.Validator(location)
	if validationError != nil {
		return c.JSON(http.StatusBadRequest, &utils.GenericResponse{
			Message:          utils.DATA_INVALID,
			DeveloperMessage: validationError.Message,
		})
	}
	rider, errLocation := r.restaurantService.GetRiderLocation(c, location)
	if errLocation != nil {
		return c.JSON(errLocation.Code, utils.GenericResponse{
			Message:          errLocation.Message,
			DeveloperMessage: errLocation.DevMessage,
		})
	}
	return c.JSON(http.StatusOK, rider)
}

func (r restaurantEndpoint) NearestRider(c echo.Context) error {
	restaurantId := c.Param("id")
	response, errResponse := r.restaurantService.NearestRider(c, restaurantId)
	if errResponse != nil {
		return c.JSON(errResponse.Code, utils.GenericResponse{
			Message: utils.DATA_NOT_FOUND,
		})
	}
	return c.JSON(http.StatusOK, response)
}

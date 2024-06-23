package user

import (
	"net/http"

	"github.com/kartikayshankh/food-delivery/model"
	userService "github.com/kartikayshankh/food-delivery/service/user"
	"github.com/kartikayshankh/food-delivery/utils"

	"github.com/labstack/echo/v4"
)

type userEndpoint struct {
	userService userService.UserService
}

type UserEnpoint interface {
	Register(c echo.Context) error
	GetUserOrderHistory(c echo.Context) error
}

func NewUserEndpoint(userService userService.UserService) UserEnpoint {
	return &userEndpoint{userService: userService}
}

func (u *userEndpoint) Register(c echo.Context) error {
	user := new(model.User)
	err := c.Bind(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &utils.GenericResponse{
			Message:          utils.SOMETHING_WENT_WRONG,
			DeveloperMessage: utils.DATA_INVALID,
		})
	}

	validationError := utils.Validator(user)
	if validationError != nil {
		return c.JSON(http.StatusBadRequest, &utils.GenericResponse{
			Message:          utils.DATA_INVALID,
			DeveloperMessage: validationError.Message,
		})
	}

	errRegister := u.userService.Register(c, user)
	if errRegister != nil {
		return c.JSON(errRegister.Code, &utils.GenericResponse{
			Message: errRegister.Message,
		})
	}

	return c.JSON(http.StatusCreated, &utils.GenericResponse{
		Message: utils.USER_CREATED_SUCCESSFULLY,
	})
}

func (r *userEndpoint) GetUserOrderHistory(c echo.Context) error {
	userId := c.Param("id")
	response, errRegister := r.userService.GetUserOrderHistory(c, userId)
	if errRegister != nil {
		return c.JSON(errRegister.Code, &utils.GenericResponse{
			Message: errRegister.Message,
		})
	}
	return c.JSON(http.StatusOK, response)
}

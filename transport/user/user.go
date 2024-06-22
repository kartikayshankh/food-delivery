package user

import (
	"assignment/model"
	userService "assignment/service/user"
	"assignment/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type userEndpoint struct {
	userService userService.UserService
}

type UserEnpoint interface {
	Register(c echo.Context) error
}

func NewUserEndpoint(userService userService.UserService) UserEnpoint {
	return &userEndpoint{userService: userService}
}

func (u *userEndpoint) Register(c echo.Context) error {
	user := new(model.User)

	err := c.Bind(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.GenericResponse{Message: "data not valid"})
	}

	errRegister := u.userService.Register(c, user)
	if errRegister != nil {
		return c.JSON(http.StatusBadRequest, utils.GenericResponse{Message: "data not valid"})
	}

	return c.JSON(http.StatusCreated, utils.GenericResponse{Message: "user registered successfully"})

}


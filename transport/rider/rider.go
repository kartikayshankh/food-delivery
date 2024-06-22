package rider

import (
	"assignment/model"
	riderService "assignment/service/rider"
	"assignment/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type riderEndpoint struct {
	riderService riderService.RiderService
}

type RiderEndpoint interface {
	RegisterRider(c echo.Context) error
	UpdateRiderLocation(c echo.Context) error
	GetRiderOrderHistory(c echo.Context) error
}

func NewRiderEndpoint(riderService riderService.RiderService) RiderEndpoint {
	return &riderEndpoint{riderService: riderService}
}

func (r *riderEndpoint) RegisterRider(c echo.Context) error {
	rider := new(model.Rider)
	if err := c.Bind(rider); err != nil {
		return c.JSON(http.StatusBadRequest, utils.GenericResponse{Message: "data not valid"})
	}

	errRegister := r.riderService.ResgisterRider(c, rider)
	if errRegister != nil {
		return c.JSON(http.StatusBadRequest, utils.GenericResponse{Message: utils.DATA_NOT_FOUND})
	}

	return c.JSON(http.StatusCreated, utils.GenericResponse{Message: "Rider registered successfully"})

}

func (r *riderEndpoint) UpdateRiderLocation(c echo.Context) error {
	updateLoaction := new(model.Location)
	riderId := c.Param("id")
	if err := c.Bind(updateLoaction); err != nil {
		return c.JSON(http.StatusBadRequest, utils.GenericResponse{Message: utils.DATA_NOT_FOUND})
	}
	errRegister := r.riderService.UpdateRiderLocation(c, *updateLoaction, riderId)
	if errRegister != nil {
		return c.JSON(http.StatusBadRequest, utils.GenericResponse{Message: utils.DATA_NOT_FOUND})
	}

	return c.JSON(http.StatusOK, utils.GenericResponse{Message: "Rider registered successfully"})

}

func (r *riderEndpoint) GetRiderOrderHistory(c echo.Context) error {
	riderId := c.Param("id")
	response, errRegister := r.riderService.GetRiderOrderHistory(c, riderId)
	if errRegister != nil {
		return c.JSON(errRegister.Code, utils.GenericResponse{Message: utils.DATA_NOT_FOUND})
	}

	return c.JSON(http.StatusOK, response)
}

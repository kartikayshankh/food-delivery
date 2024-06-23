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
	Register(c echo.Context) error
	UpdateRiderLocation(c echo.Context) error
	GetRiderOrderHistory(c echo.Context) error
}

func NewRiderEndpoint(riderService riderService.RiderService) RiderEndpoint {
	return &riderEndpoint{riderService: riderService}
}

func (r *riderEndpoint) Register(c echo.Context) error {
	rider := new(model.Rider)
	if err := c.Bind(rider); err != nil {
		return c.JSON(http.StatusBadRequest, &utils.GenericResponse{
			Message:          utils.SOMETHING_WENT_WRONG,
			DeveloperMessage: utils.DATA_INVALID,
		})
	}

	validationError := utils.Validator(rider)
	if validationError != nil {
		return c.JSON(http.StatusBadRequest, &utils.GenericResponse{
			Message:          utils.DATA_INVALID,
			DeveloperMessage: validationError.Message,
		})
	}

	errRegister := r.riderService.ResgisterRider(c, rider)
	if errRegister != nil {
		return c.JSON(errRegister.Code, &utils.GenericResponse{
			Message:          errRegister.Message,
			DeveloperMessage: errRegister.DevMessage,
		})
	}

	return c.JSON(http.StatusCreated, &utils.GenericResponse{
		Message: utils.RIDER_CREATED_SUCCESSFULLY,
	})

}

func (r *riderEndpoint) UpdateRiderLocation(c echo.Context) error {
	updateLoaction := new(model.Location)
	riderId := c.Param("id")
	if err := c.Bind(updateLoaction); err != nil {
		return c.JSON(http.StatusBadRequest, &utils.GenericResponse{
			Message:          utils.SOMETHING_WENT_WRONG,
			DeveloperMessage: utils.DATA_INVALID,
		})
	}

	validationError := utils.Validator(updateLoaction)
	if validationError != nil {
		return c.JSON(http.StatusBadRequest, &utils.GenericResponse{
			Message:          utils.DATA_INVALID,
			DeveloperMessage: validationError.Message,
		})
	}

	errUpdateLocation := r.riderService.UpdateRiderLocation(c, *updateLoaction, riderId)
	if errUpdateLocation != nil {
		return c.JSON(errUpdateLocation.Code, &utils.GenericResponse{

			Message: errUpdateLocation.Message,
		})
	}

	return c.JSON(http.StatusOK, &utils.GenericResponse{
		Message: utils.LOCATION_UPDATED_SUCCESSFULLY,
	})

}

func (r *riderEndpoint) GetRiderOrderHistory(c echo.Context) error {
	riderId := c.Param("id")
	response, errRegister := r.riderService.GetRiderOrderHistory(c, riderId)
	if errRegister != nil {
		return c.JSON(errRegister.Code, &utils.GenericResponse{
			Message: errRegister.Message,
		})
	}
	return c.JSON(http.StatusOK, response)
}

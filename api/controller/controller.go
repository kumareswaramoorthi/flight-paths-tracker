package controller

import (
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/constants"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/dto"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/errors"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/logging"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/service"
)

type FlightTrackerController interface {
	FindSourceAndDestination(c *gin.Context)
}

type flightTrackerController struct {
	flightTrackerService service.FlightTrackerService
}

func NewFlightTrackerController(flightTrackerService service.FlightTrackerService) FlightTrackerController {
	return flightTrackerController{
		flightTrackerService: flightTrackerService,
	}
}

// @title FLIGHT PATHS TRACKER API
// Find Flight Source And Destination godoc
// @Tags Find Source And Destination
// @Accept json
// @Produce  json
// @Description Find source and destination
// @Success 200 {object} []string
// @Failure 400 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Param Tickets body dto.Tickets true "request body"
// @Router /track [POST]
func (ftc flightTrackerController) FindSourceAndDestination(c *gin.Context) {
	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "FlightTrackerController").
		WithField(constants.Method, "FindSourceAndDestination")

	tickets := new(dto.Tickets)

	//Bind json to tickets object
	if err := c.ShouldBindJSON(tickets); err != nil {
		logger.Errorf("ShouldBindJSON - %s", err.Error())
		c.AbortWithStatusJSON(errors.ErrBadRequest.HttpStatusCode, errors.ErrBadRequest)
		return
	}

	//validate tickets
	err := ftc.flightTrackerService.ValidateTickets(c, tickets.Tickets)
	if err != nil {
		logger.Errorf("ValidateTickets - %s", err.Error())
		c.AbortWithStatusJSON(err.HttpStatusCode, err)
		return
	}

	//find source and destination
	srcdst, err := ftc.flightTrackerService.FindSourceAndDestination(c, tickets.Tickets)
	if err != nil {
		logger.Errorf("FindSourceAndDestination - %s", err.Error())
		c.AbortWithStatusJSON(err.HttpStatusCode, err)
		return
	}

	c.JSON(http.StatusOK, srcdst)
	logger.Info("FindSourceAndDestination call completed")
}

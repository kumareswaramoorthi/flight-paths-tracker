package service

import (
	"strings"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/constants"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/errors"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/logging"
)

type FlightTrackerService interface {
	FindSourceAndDestination(c *gin.Context, tickets [][]string) ([]string, *errors.ErrorResponse)
	ValidateTickets(c *gin.Context, tickets [][]string) *errors.ErrorResponse
}

type flightTrackerService struct {
}

func NewFlightTrackerService() FlightTrackerService {
	return &flightTrackerService{}
}

func (fts *flightTrackerService) FindSourceAndDestination(c *gin.Context, tickets [][]string) ([]string, *errors.ErrorResponse) {
	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "FlightTrackerService").
		WithField(constants.Method, "FindSourceAndDestination")

	srcdst := make([]string, 2)
	flightPath := make(map[string]int)

	for _, ticket := range tickets {
		flightPath[ticket[0]]--
		flightPath[ticket[1]]++
	}

	//validate the number of source and destinations
	for k, v := range flightPath {
		if v > 1 || v < -1 {
			logger.Errorf("Error invalid number of source and destinations - %s", errors.ErrUnableToTrack.Error())
			return nil, errors.ErrUnableToTrack
		}
		if v == 0 {
			delete(flightPath, k)
		}
	}
	//check if there is only one source and one destination
	if len(flightPath) != 2 {
		logger.Errorf("Error invalid flght paths -  %s", errors.ErrUnableToTrack.Error())
		return nil, errors.ErrUnableToTrack
	}

	//fetch the source which has value as -1 and destination which has value as 1 from the flightPath map
	for k, v := range flightPath {
		switch v {
		case -1:
			srcdst[0] = k
		case 1:
			srcdst[1] = k
		}
	}
	return srcdst, nil
}

func (fts *flightTrackerService) ValidateTickets(c *gin.Context, tickets [][]string) *errors.ErrorResponse {
	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "FlightTrackerService").
		WithField(constants.Method, "ValidateTickets")

	for _, ticket := range tickets {
		//check if each ticket has exactly one source and one destination
		if len(ticket) != 2 {
			logger.Errorf("Error ticket format - %s", errors.ErrInvalidTicket.Error())
			return errors.ErrInvalidTicket
		}
		for _, place := range ticket {
			//check the source and destination's naming convention
			if strings.ToUpper(place) != place || len(place) != 3 {
				logger.Errorf("Error in source or destination name - %s", errors.ErrInvalidTicket.Error())
				return errors.ErrInvalidTicket
			}
		}
	}
	return nil
}

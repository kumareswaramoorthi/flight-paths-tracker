package service

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/errors"
	"github.com/stretchr/testify/suite"
)

type FlightTrackerServiceTestSuite struct {
	suite.Suite
	context              *gin.Context
	recorder             *httptest.ResponseRecorder
	mockCtrl             *gomock.Controller
	flightTrackerService FlightTrackerService
}

func TestFlightTrackerService(t *testing.T) {
	suite.Run(t, new(FlightTrackerServiceTestSuite))
}

func (suite FlightTrackerServiceTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *FlightTrackerServiceTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.flightTrackerService = NewFlightTrackerService()
}

func (suite *FlightTrackerServiceTestSuite) TestGetSourceAndDestinationSuccessfully() {
	var tickets [][]string
	tickets = append(tickets, []string{"IND", "EWR"}, []string{"SFO", "ATL"}, []string{"GSO", "IND"}, []string{"ATL", "GSO"})

	actualResponse, err := suite.flightTrackerService.FindSourceAndDestination(suite.context, tickets)
	expectedResponse := []string{"SFO", "EWR"}

	suite.Equal(expectedResponse, actualResponse)
	suite.Nil(err)
}

func (suite *FlightTrackerServiceTestSuite) TestGetSrcDstReturnsErrorIfPathsInvalid() {
	var tickets [][]string
	tickets = append(tickets, []string{"IND", "EWR"}, []string{"SFO", "ATL"}, []string{"GSO", "IND"}, []string{"ATL", "GSO"}, []string{"IND", "EWR"})
	_, err := suite.flightTrackerService.FindSourceAndDestination(suite.context, tickets)

	suite.Equal(errors.ErrUnableToTrack, err)
	suite.NotNil(err)
}

func (suite *FlightTrackerServiceTestSuite) TestGetSrcDstnReturnsErrIfSrcAndDstInvalid() {
	var tickets [][]string
	tickets = append(tickets, []string{"IND", "EWR"}, []string{"IND", "EWR"}, []string{"IND", "EWR"})
	_, err := suite.flightTrackerService.FindSourceAndDestination(suite.context, tickets)

	suite.Equal(errors.ErrUnableToTrack, err)
	suite.NotNil(err)
}

func (suite *FlightTrackerServiceTestSuite) TestValidateTicketsSuccessfullyIfNoError() {
	var tickets [][]string
	tickets = append(tickets, []string{"IND", "EWR"}, []string{"IND", "EWR"}, []string{"IND", "EWR"})
	err := suite.flightTrackerService.ValidateTickets(suite.context, tickets)

	suite.Nil(err)
}

func (suite *FlightTrackerServiceTestSuite) TestValidateTicketsReturnsErrIfTicketInvalid() {
	var tickets [][]string
	tickets = append(tickets, []string{"IND", "EWR", "EWR"}, []string{"IND", "EWR", "ATL"}, []string{"IND", "EWR"})
	err := suite.flightTrackerService.ValidateTickets(suite.context, tickets)

	suite.Equal(errors.ErrInvalidTicket, err)
	suite.NotNil(err)
}

func (suite *FlightTrackerServiceTestSuite) TestValidateTicketsReturnsErrIfPlaceNameIsInvalid() {
	var tickets [][]string
	tickets = append(tickets, []string{"ind", "ewr"}, []string{"IND", "ATL"}, []string{"IND", "EWR"})
	err := suite.flightTrackerService.ValidateTickets(suite.context, tickets)

	suite.Equal(errors.ErrInvalidTicket, err)
	suite.NotNil(err)
}

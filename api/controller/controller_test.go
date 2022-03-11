package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/controller/mocks"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/dto"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/errors"
	"github.com/stretchr/testify/suite"
)

type FlightTrackerControllerTestSuite struct {
	suite.Suite
	context                  *gin.Context
	recorder                 *httptest.ResponseRecorder
	mockCtrl                 *gomock.Controller
	mockFlightTrackerService *mocks.MockFlightTrackerService
	flightTrackerController  FlightTrackerController
}

func TestFlightTrackerController(t *testing.T) {
	suite.Run(t, new(FlightTrackerControllerTestSuite))
}

func (suite FlightTrackerControllerTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *FlightTrackerControllerTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.mockFlightTrackerService = mocks.NewMockFlightTrackerService(suite.mockCtrl)
	suite.flightTrackerController = NewFlightTrackerController(suite.mockFlightTrackerService)
}

func (suite *FlightTrackerControllerTestSuite) TestFindSourceAndDestinationSuccessfully() {
	var tickets [][]string
	tickets = append(tickets, []string{"IND", "EWR"}, []string{"SFO", "ATL"}, []string{"GSO", "IND"}, []string{"ATL", "GSO"})
	payload := dto.Tickets{
		Tickets: tickets,
	}

	req, _ := json.Marshal(payload)
	expectedResponse := []string{"SFO", "EWR"}
	response, _ := json.Marshal(expectedResponse)
	suite.context.Request, _ = http.NewRequest("POST", "/track", bytes.NewBufferString(string(req)))

	suite.mockFlightTrackerService.EXPECT().ValidateTickets(suite.context, payload.Tickets).Return(nil)
	suite.mockFlightTrackerService.EXPECT().FindSourceAndDestination(suite.context, payload.Tickets).Return(expectedResponse, nil)
	suite.flightTrackerController.FindSourceAndDestination(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.JSONEq(string(response), suite.recorder.Body.String())
}

func (suite *FlightTrackerControllerTestSuite) TestFindSourceAndDestinationFailsIfNoPayload() {
	suite.context.Request, _ = http.NewRequest("POST", "/track", nil)
	suite.flightTrackerController.FindSourceAndDestination(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *FlightTrackerControllerTestSuite) TestFindSourceAndDestinationFailsIfTicketInvalid() {
	var tickets [][]string
	tickets = append(tickets, []string{"IND", "EWR"}, []string{"SFO", "ATL", "EWR"}, []string{"GSO", "IND"}, []string{"ATL", "GSO"})
	payload := dto.Tickets{
		Tickets: tickets,
	}
	req, _ := json.Marshal(payload)
	suite.context.Request, _ = http.NewRequest("POST", "/track", bytes.NewBufferString(string(req)))

	suite.mockFlightTrackerService.EXPECT().ValidateTickets(suite.context, payload.Tickets).Return(errors.ErrInvalidTicket)
	suite.flightTrackerController.FindSourceAndDestination(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *FlightTrackerControllerTestSuite) TestFindSourceAndDestinationFailsIfPathInvalid() {
	var tickets [][]string
	tickets = append(tickets, []string{"IND", "EWR"}, []string{"IND", "EWR"})
	payload := dto.Tickets{
		Tickets: tickets,
	}
	req, _ := json.Marshal(payload)
	suite.context.Request, _ = http.NewRequest("POST", "/track", bytes.NewBufferString(string(req)))

	suite.mockFlightTrackerService.EXPECT().ValidateTickets(suite.context, payload.Tickets).Return(nil)
	suite.mockFlightTrackerService.EXPECT().FindSourceAndDestination(suite.context, payload.Tickets).Return(nil, errors.ErrUnableToTrack)
	suite.flightTrackerController.FindSourceAndDestination(suite.context)

	suite.Equal(http.StatusUnprocessableEntity, suite.recorder.Code)
}

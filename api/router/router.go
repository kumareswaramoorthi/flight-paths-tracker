package router

import (
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/controller"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/logging"
	"github.com/kumareswaramoorthi/flight-paths-tracker/api/service"
	"github.com/kumareswaramoorthi/flight-paths-tracker/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	//Get default router from gin
	router := gin.Default()

	//create a global logger for the server
	apiLoggerEntry := logging.NewLoggerEntry()

	//router use the global logger
	router.Use(logging.LoggingMiddleware(apiLoggerEntry))
	router.Use(requestid.New())

	//swagger init
	docs.SwaggerInfo.Title = "FLIGHT PATHS TRACKER API"
	docs.SwaggerInfo.Description = "This lists down the endpoints that are part of FLIGHT PATHS TRACKER API server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http"}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//health check API
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})

	trackService := service.NewFlightTrackerService()
	trackController := controller.NewFlightTrackerController(trackService)

	//route to fetch source and destination from tickets
	router.POST("/track", trackController.FindSourceAndDestination)

	return router
}

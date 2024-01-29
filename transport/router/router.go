package router

import (
	"net/http"
	"synapso/config"
	controller2 "synapso/transport/controller"
	middles "synapso/transport/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Init initialize the routing of this application.
func Init(e *echo.Echo, conf *config.Config) {
	if conf.Extension.CorsEnabled {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
			AllowOrigins:     []string{"*"},
			AllowHeaders: []string{
				echo.HeaderAccessControlAllowHeaders,
				echo.HeaderContentType,
				echo.HeaderContentLength,
				echo.HeaderAcceptEncoding,
				echo.HeaderXCSRFToken,
				echo.HeaderAuthorization,
			},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPost,
			},
			MaxAge: 86400,
		}))
	}

	e.HTTPErrorHandler = controller2.JSONErrorHandler
	e.Use(middleware.Recover())

	api := e.Group("/api")
	api.GET("/health", controller2.GetHealthCheck())
	api.GET("/login", controller2.GetLoginAccount(), middles.AuthAndExtractUserMiddleware)
	api.POST("/sign-up", controller2.UserSignUp)
	api.POST("/login", controller2.UserLogin)

	researcher := api.Group("/researcher")
	researcher.Use(middles.AuthAndExtractUserMiddleware, middles.AuthResearcherMiddleware)
	researcher.POST("/recall", controller2.SaveRecall)
	researcher.GET("/recall/:id", controller2.GetRecall)
	researcher.POST("/recognition", controller2.SaveRecognitionExperiment)
	researcher.GET("/recognition/:id", controller2.GetRecognitionExperiment)
	researcher.GET("/experiment", controller2.ListExperiments)
	researcher.GET("/result/:id", controller2.GetExperimentResultCsv)

	subject := api.Group("/subject")
	subject.Use(middles.AuthAndExtractUserMiddleware, middles.AuthSubjectMiddleware)
	subject.POST("/result", controller2.SaveExperimentResult)
	subject.GET("/result/:id", controller2.GetExperimentResult)
	subject.GET("/result", controller2.GetExperimentResults)
	subject.GET("/recognition/:id", controller2.GetRecognitionExperiment)
	subject.GET("/recall/:id", controller2.GetRecall)
	subject.GET("/experiment", controller2.ListExperiments)
	//subject.POST("/recall", controller2.SaveRecallResult, middles.AuthAndExtractUserMiddleware, middles.AuthSubjectMiddleware)

}

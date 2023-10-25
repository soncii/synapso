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
	//e.POST(controller2.APIAccountLogout, controller2.PostLogout())

}

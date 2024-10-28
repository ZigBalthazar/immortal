package routes

import (
	"github.com/dezh-tech/immortal/management"
	"github.com/dezh-tech/immortal/management/handlers"
	"github.com/dezh-tech/immortal/management/services/token"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigureRoutes(server *management.Server) {
	authHandler := handlers.NewAuthHandler(server)

	server.Echo.Use(middleware.Logger())

	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	server.Echo.POST("/login", authHandler.Login)

	r := server.Echo.Group("/a")
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(token.JwtCustomClaims)
		},
		SigningKey: []byte(server.Config.AuthConfig.AccessSecret),
	}
	r.Use(echojwt.WithConfig(config))
}
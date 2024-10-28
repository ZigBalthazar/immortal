package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dezh-tech/immortal/management"
	"github.com/dezh-tech/immortal/management/models"
	"github.com/dezh-tech/immortal/management/repositories"
	requests "github.com/dezh-tech/immortal/management/requestes"
	"github.com/dezh-tech/immortal/management/responses"
	tokenService "github.com/dezh-tech/immortal/management/services/token"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	server *management.Server
}

func NewAuthHandler(server *management.Server) *AuthHandler {
	return &AuthHandler{server: server}
}

// Login godoc
//
//	@Summary		Authenticate a user
//	@Description	Perform user login
//	@ID				user-login
//	@Tags			User Actions
//	@Accept			json
//	@Produce		json
//	@Param			params	body		requests.LoginRequest	true	"User's credentials"
//	@Success		200		{object}	responses.LoginResponse
//	@Failure		401		{object}	responses.Error
//	@Router			/login [post]
func (authHandler *AuthHandler) Login(c echo.Context) error {
	loginRequest := new(requests.LoginRequest)

	if err := c.Bind(loginRequest); err != nil {
		return err
	}
	if err := loginRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty or not valid")
	}

	user := models.User{}
	userRepository := repositories.NewUserRepository(authHandler.server.DataBase)
	userRepository.GetUserByEmail(context.Background(), &user, loginRequest.Email)

	// if user.ID == 0 || (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)) != nil) {
	// 	return responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
	// }

	fmt.Println(user)
	if user.ID == 0 {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
	}

	tokenService := tokenService.NewTokenService(&authHandler.server.Config.AuthConfig)
	accessToken, exp, err := tokenService.CreateAccessToken(&user)
	if err != nil {
		return err
	}
	res := responses.NewLoginResponse(accessToken, exp)

	return responses.Response(c, http.StatusOK, res)
}

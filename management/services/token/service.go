package token

import (
	"github.com/dezh-tech/immortal/management/models"
	"github.com/golang-jwt/jwt/v5"
)

const ExpireCount = 2
const ExpireRefreshCount = 168

type JwtCustomClaims struct {
	Name string `json:"name"`
	ID   int   `json:"id"`
	jwt.RegisteredClaims
}

type ServiceWrapper interface {
	CreateAccessToken(user *models.User) (accessToken string, exp int64, err error)
	CreateRefreshToken(user *models.User) (t string, err error)
}


type Config struct {
	AccessSecret string
}

type Service struct {
	config *Config
}

func NewTokenService(cfg *Config) *Service {
	return &Service{
		config: cfg,
	}
} 
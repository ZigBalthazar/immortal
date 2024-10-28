package management

import (
	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/management/services/token"
)

type Config struct {
	Bind           string `yaml:"bind"`
	Port           uint16 `yaml:"port"`
	DatabaseConfig database.Config
	AuthConfig token.Config
}

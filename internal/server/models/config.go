package models

import (
	"github.com/alexsniffin/website/pkg/models"
)

type Config struct {
	Environment string
	Logger      models.Logger
	HTTPRouter  HTTPRouterConfig
	HTTPServer  HTTPServerConfig
}

type HTTPRouterConfig struct {
	TimeoutSec     int
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type HTTPServerConfig struct {
	Port int
}

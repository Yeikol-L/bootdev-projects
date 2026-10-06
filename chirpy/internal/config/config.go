package config

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type ApiConfig struct {
	platform  string
	port      string
	jwtSecret string
	polkaKey  string
}

// Carga la configuracion de la rest api, es necesario de que ya esten cargadas las variables de entorno  cuando se ejecute esta funcion.
// Variables de entorno necesarias: PLATFORM  (dev en desarrollo), PORT (puerto en que core la res api), JWT_SECRET (secreto para firmar tokens)
func NewApiConfig() (ApiConfig, error) {
	platform := os.Getenv("PLATFORM")
	port := os.Getenv("PORT")
	jwt := os.Getenv("JWT_SECRET")
	polka := os.Getenv("POLKA_KEY")
	if platform == "" {
		return ApiConfig{}, fmt.Errorf("variable de entorno platform obligatoria")
	}
	if port == "" {
		return ApiConfig{}, fmt.Errorf("variable de entorno port obligatoria")
	}
	if jwt == "" {
		return ApiConfig{}, fmt.Errorf("variable de entorno jwt_secret obligatoria")
	}
	if polka == "" {
		return ApiConfig{}, fmt.Errorf("variable de entorno polka_key obligatoria")
	}
	return ApiConfig{
		platform:  platform,
		port:      port,
		jwtSecret: jwt,
		polkaKey:  polka,
	}, nil
}

func (c *ApiConfig) IsDev() bool {
	return c.platform == "dev"
}
func (c *ApiConfig) GetPort() string {
	return c.port
}
func (c *ApiConfig) GetJWTSecret() string {
	return c.jwtSecret
}
func (c *ApiConfig) GetPolkaKey() string {
	return c.polkaKey
}

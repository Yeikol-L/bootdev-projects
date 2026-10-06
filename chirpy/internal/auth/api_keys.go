package auth

import (
	"errors"
	"net/http"
	"strings"
)

const prefix = "apikey"

func GetApiKey(header http.Header) (string, error) {
	authorizationContent := header.Get("Authorization")
	if authorizationContent == "" {
		return "", errors.New("header Authorization obligatorio")
	}

	if len(authorizationContent) < len(prefix) || !strings.EqualFold(prefix, authorizationContent[:len(prefix)]) {
		return "", errors.New("formato ApiKey: valor")
	}

	apiKey := strings.Trim(authorizationContent[len(prefix):], " ")
	return apiKey, nil
}

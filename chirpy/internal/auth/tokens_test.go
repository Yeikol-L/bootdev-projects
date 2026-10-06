package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	defaultID := uuid.New()
	duration := time.Millisecond * 100
	tests := []struct {
		name             string
		userID           uuid.UUID
		signTokenSecret  string
		sleep            *time.Duration
		signDuration     time.Duration
		checkTokenSecret string
		want             struct {
			id  uuid.UUID
			err bool
		}
	}{
		{
			"Firma correctamente con happy path",
			defaultID, "123456", nil, time.Minute * 10, "123456", struct {
				id  uuid.UUID
				err bool
			}{defaultID, false},
		},
		{
			"Falla por secretos diferentes",
			defaultID, "123456", nil, time.Minute * 10, "12345678", struct {
				id  uuid.UUID
				err bool
			}{defaultID, true},
		},
		{
			"Falla por token expirado",
			defaultID, "123456", &duration, duration / 2, "123456", struct {
				id  uuid.UUID
				err bool
			}{defaultID, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.userID, tt.signTokenSecret, tt.signDuration)
			if err != nil && !tt.want.err {
				t.Fatalf("Error al firmar token: %v", err)
			}
			if tt.sleep != nil {
				time.Sleep(*tt.sleep)
			}
			id, err := ValidateJWT(token, tt.checkTokenSecret)
			if err != nil && !tt.want.err {
				t.Fatalf("Error al validar token: %v", err)
			}
			if id != tt.want.id && !tt.want.err {
				t.Fatalf("El id %s no es igual a %s", tt.userID, tt.want.id)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		token   string
		err     bool
	}{
		{
			"Happy PATH",
			map[string][]string{"Authorization": []string{"Bearer 123456"}},
			"123456",
			false,
		},
		{
			"Happy PATH",
			map[string][]string{"Authorization": []string{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiZDBkN2Y3Y2QtNDYxMi00NWRhLTkzODAtNjczMmJiYzdjZDY1IiwiZXhwIjoxNzkxMjU1MTUyLCJpYXQiOjE3OTEyNTE1NTJ9.5jpX-ZJlwIa51Ydh-zppgyeBGTfYKGIise32vWnsOmI"}},
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiZDBkN2Y3Y2QtNDYxMi00NWRhLTkzODAtNjczMmJiYzdjZDY1IiwiZXhwIjoxNzkxMjU1MTUyLCJpYXQiOjE3OTEyNTE1NTJ9.5jpX-ZJlwIa51Ydh-zppgyeBGTfYKGIise32vWnsOmI",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)
			if err != nil && !tt.err {
				t.Fatalf("Error no esperado %v", err)
			}
			if token != tt.token && !tt.err {
				t.Fatalf("Token esperado: %s recibido: %s", tt.token, token)
			}
		})
	}
}

package httpx

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bootdotdev/learn-web-security/internal/templates"
	"github.com/bootdotdev/learn-web-security/internal/textutils"
)

func RespondWithError(responseWriter http.ResponseWriter, code int, message string) {
	message = textutils.StripANSI(message)
	msg := "The request failed. Try again, or return to the store."
	if code >= 500 {
		fmt.Printf("Responding with error code %v, message: %v\n", code, message)
	}
	if code < 500 {
		msg = message
	}
	RespondWithJSON(responseWriter, code, map[string]string{"error": msg})
}

func RespondWithJSON(responseWriter http.ResponseWriter, code int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(code)
	encoder := json.NewEncoder(responseWriter)
	if err := encoder.Encode(payload); err != nil {
		log.Printf("Error respondWithJSON: status code %v: %v", code, err)
	}
}

func RespondWithErrorPage(responseWriter http.ResponseWriter, renderer *templates.Renderer, statusCode int, title, message string) error {
	msg := "The request failed. Try again, or return to the store."
	t := "Something Went Wrong"
	if statusCode < 500 {
		msg = message
		t = title
	}
	return renderer.Render(responseWriter, statusCode, "error", templates.ErrorPage{
		Title:      t,
		StatusCode: statusCode,
		Message:    msg,
	})
}

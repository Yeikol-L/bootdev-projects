package pawpal

import (
	"crypto/subtle"
	"fmt"
)

type WebhookOutcome string

const (
	WebhookUnauthorized WebhookOutcome = "unauthorized"
	WebhookMalformed    WebhookOutcome = "malformed"
	WebhookApproved     WebhookOutcome = "approved"
)

type WebhookVerification struct {
	Outcome WebhookOutcome
	OrderID int64
}

func CreateCheckoutURL(orderID int64) string {
	return fmt.Sprintf("https://pawpal.example/checkout?orderId=%d", orderID)
}

func VerifyWebhook(providedKey []byte, expectedKey []byte, payload any) WebhookVerification {
	if subtle.ConstantTimeCompare(providedKey, expectedKey) == 0 {
		return WebhookVerification{Outcome: WebhookUnauthorized, OrderID: 0}
	}
	payloadRecord, _ := payload.(map[string]any)
	orderID, orderIDFound := payloadRecord["orderId"].(float64)
	orderStatus, orderStatusFound := payloadRecord["status"].(string)

	if orderID < 0 || orderStatus != "approved" || !orderIDFound || !orderStatusFound {
		return WebhookVerification{Outcome: WebhookMalformed, OrderID: int64(orderID)}
	}
	return WebhookVerification{Outcome: WebhookApproved, OrderID: int64(orderID)}
}

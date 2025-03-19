package config

import (
	"log"
	"net/http"
	"os"

	"github.com/yantology/retail-pro-be/pkg/customerror"
)

type ResendApi struct {
	apiKey       string
	resendDomain string
	resendName   string
}

func InitResendConfig() (*ResendApi, *customerror.CustomError) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Println("Resend API key is not set")
		return nil, customerror.NewCustomError(nil, "Resend API key is not set", http.StatusUnauthorized)
	}
	resendDomain := os.Getenv("RESEND_DOMAIN")
	log.Println("Resend domain => ", resendDomain)
	if resendDomain == "" {
		return nil, customerror.NewCustomError(nil, "Resend domain is not set", http.StatusUnauthorized)
	}

	resendName := os.Getenv("RESEND_NAME")
	if resendName == "" {
		return nil, customerror.NewCustomError(nil, "Resend name is not set", http.StatusUnauthorized)
	}

	return &ResendApi{apiKey: apiKey, resendDomain: resendDomain, resendName: resendName}, nil
}

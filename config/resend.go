package config

import (
	"log"
	"net/http"
	"os"

	"github.com/yantology/simple-pos/pkg/customerror"
)

type ResendApi struct {
	ApiKey       string
	ResendDomain string
	ResendName   string
}

func InitResendConfig() (*ResendApi, *customerror.CustomError) {
	apiKey := os.Getenv("RESEND_API_KEY")
	log.Println("Resend API key => ", apiKey != "")
	if apiKey == "" {
		log.Println("Resend API key is not set")
		return nil, customerror.NewCustomError(nil, "Resend API key is not set", http.StatusUnauthorized)
	}

	resendDomain := os.Getenv("RESEND_DOMAIN")
	log.Println("Resend domain => ", resendDomain)
	if resendDomain == "" {
		log.Println("Resend domain is not set")
		return nil, customerror.NewCustomError(nil, "Resend domain is not set", http.StatusUnauthorized)
	}

	resendName := os.Getenv("RESEND_NAME")
	log.Println("Resend name => ", resendName)
	if resendName == "" {
		log.Println("Resend name is not set")
		return nil, customerror.NewCustomError(nil, "Resend name is not set", http.StatusUnauthorized)
	}

	return &ResendApi{ApiKey: apiKey, ResendDomain: resendDomain, ResendName: resendName}, nil
}

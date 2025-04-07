package resendutils

import (
	"fmt"
	"net/http"

	"github.com/resend/resend-go/v2"
	"github.com/yantology/golang-starter-template/pkg/customerror"
)

type ResendUtilsInterface interface {
	Send(html, subject string, to []string) *customerror.CustomError
}

type ResendUtils struct {
	apiKey     string
	fromDomain string
}

func NewResendUtils(apiKey, fromDomain string) ResendUtilsInterface {
	return &ResendUtils{
		apiKey:     apiKey,
		fromDomain: fromDomain,
	}
}

func (r *ResendUtils) Send(html, subject string, to []string) *customerror.CustomError {
	client := resend.NewClient(r.apiKey)

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("Yantology <activation@%s>", r.fromDomain),
		To:      to,
		Subject: subject,
		Html:    html,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return customerror.NewCustomError(err, "Failed to send email", http.StatusInternalServerError)
	}

	return nil
}

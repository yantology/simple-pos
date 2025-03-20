package resendutils

import (
	"fmt"
	"net/http"

	"github.com/resend/resend-go/v2"
	"github.com/yantology/retail-pro-be/pkg/customerror"
)

// ResendUtils interface for sending emails using Resend API
type ResendUtilsInterface interface {
	Send(html, subject string, to []string) *customerror.CustomError
}

type ResendUtils struct {
	apiKey       string
	resendDomain string
	resendName   string
}

// NewResendUtils creates a new instance of ResendUtils
func NewResendUtils(apiKey, resendDomain, resendName string) ResendUtilsInterface {
	return &ResendUtils{
		apiKey:       apiKey,
		resendDomain: resendDomain,
		resendName:   resendName,
	}
}

func (r *ResendUtils) Send(html, subject string, to []string) *customerror.CustomError {
	client := resend.NewClient(r.apiKey).Emails

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <noreply@%s>", r.resendName, r.resendDomain),
		To:      to,
		Subject: subject,
		Html:    html,
	}

	sent, err := client.Send(params)
	if err != nil {
		return customerror.NewCustomError(err, "Failed to send email", http.StatusInternalServerError)
	}

	fmt.Printf("Email sent successfully with ID: %s\n", sent.Id)
	return nil
}

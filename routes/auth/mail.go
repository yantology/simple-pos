package auth

// EmailTemplateInterface defines methods for generating email content
type EmailTemplateInterface interface {
	GenerateRegistrationEmail(email, token string) string
	GeneratePasswordResetEmail(email, token string) string
}

type emailTemplate struct{}

// NewEmailTemplate creates a new email template generator
func NewEmailTemplate() EmailTemplateInterface {
	return &emailTemplate{}
}

// GenerateRegistrationEmail creates a simple email template for registration
func (e *emailTemplate) GenerateRegistrationEmail(email, token string) string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Registration Code</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2>Registration Activation Code</h2>
        <p>Hello,</p>
        <p>Here is your registration activation code:</p>
        <div style="font-size: 24px; font-weight: bold; text-align: center; 
                    letter-spacing: 5px; margin: 20px 0; color: #0066cc;">
            ` + token + `
        </div>
        <p>This code will expire in 15 minutes.</p>
        <p>If you did not request this code, please ignore this email.</p>
        <hr>
        <p style="font-size: 12px; color: #666;">
            This is an automated email, please do not reply.
        </p>
    </div>
</body>
</html>`
}

// GeneratePasswordResetEmail creates a simple email template for password reset
func (e *emailTemplate) GeneratePasswordResetEmail(email, token string) string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Reset Code</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2>Password Reset Code</h2>
        <p>Hello,</p>
        <p>Here is your password reset code:</p>
        <div style="font-size: 24px; font-weight: bold; text-align: center; 
                    letter-spacing: 5px; margin: 20px 0; color: #0066cc;">
            ` + token + `
        </div>
        <p>This code will expire in 15 minutes.</p>
        <p>If you did not request this code, please ignore this email.</p>
        <hr>
        <p style="font-size: 12px; color: #666;">
            This is an automated email, please do not reply.
        </p>
    </div>
</body>
</html>`
}

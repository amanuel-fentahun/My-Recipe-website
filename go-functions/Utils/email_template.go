package utils

import "fmt"

const (
	SubjectPasswordReset     = "Password reset code"
	SubjectEmailVerification = "Verify your email address"
)

func GetPasswordResetTemplate(newCode string) string {
	body := "<p>You need to insert this code in order to keep resetting your password.</p>"
	body += "<p>Your password reset code:</p>"
	body += fmt.Sprintf("<h3 style='color: #FF5733; font-size: 24px;'>%s</h3>", newCode)
	body += "<p>If you did not request this password reset code, please ignore this email.</p>"
	body += "<p>Thanks,<br/>The Tafach Kitchen Team</p>"
	return body
}

func GetEmailVerificationTemplate(newCode string) string {
	body := "<h2>Welcome to Tafach Kitchen!</h2>"
	body += "<p>Thank you for signing up. Please verify your email address to activate your account.</p>"
	body += "<p>Your email verification code:</p>"
	body += fmt.Sprintf("<h3 style='color: #2ECC71; font-size: 24px;'>%s</h3>", newCode)
	body += "<p>This code is valid for 15 minutes. If you did not create an account, please ignore this message.</p>"
	body += "<p>Happy Cooking,<br/>The Tafach Kitchen Team</p>"
	return body
}

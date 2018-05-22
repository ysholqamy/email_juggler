package email

import "os"

var (
	mailgunBaseURL = os.Getenv("MAILGUN_BASE_URL")
	mailgunDomain  = os.Getenv("MAILGUN_DOMAIN")
	mailgunKey     = os.Getenv("MAILGUN_KEY")

	sendgridBaseURL = os.Getenv("SENDGRID_BASE_URL")
	sendgridKey     = os.Getenv("SENDGRID_KEY")
)

package email

import "os"

var (
	mailgunBaseURL = "https://api.mailgun.net/v3"
	mailgunDomain  = os.Getenv("MAILGUN_DOMAIN")
	mailgunKey     = os.Getenv("MAILGUN_KEY")

	sendgridBaseURL = "https://api.sendgrid.com/api/mail.send.json"
	sendgridKey     = os.Getenv("SENDGRID_KEY")
)

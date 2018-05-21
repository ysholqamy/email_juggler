package email_test

import (
	"testing"

	. "github.com/ysholqamy/uber-email-challenge/email"
)

var validMessage = Message{
	From:    "email@example.com",
	To:      "youssefsholqamy@gmail.com",
	Subject: "subject",
	Text:    "Hello from tests",
}

func TestValidMailgun(t *testing.T) {
	err := DefaultMailgun.Send(validMessage)
	if err != nil {
		t.Errorf("Failed to send valid message using mailgun. got: %v", err)
	}
}

func TestValidSendgrid(t *testing.T) {
	err := DefaultSendgrid.Send(validMessage)
	if err != nil {
		t.Errorf("Failed to send valid message using Sendgrid. got: %v", err)
	}
}

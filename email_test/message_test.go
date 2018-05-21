package email_test

import (
	"reflect"
	"testing"

	. "github.com/ysholqamy/email_juggler/email"
)

var validMessage = Message{
	From:    "email@example.com",
	To:      "youssefsholqamy@gmail.com",
	Subject: "subject",
	Text:    "Hello from tests",
}

var missingFieldMessage = Message{
	From:    "email@example.com",
	To:      "email@example.com",
	Subject: "subject",
}

func TestValidateValid(t *testing.T) {
	err := validMessage.Validate()
	if err != nil {
		t.Error("Validate: valid message is marked as invalid")
	}
}

func TestValidateMissingField(t *testing.T) {
	err := missingFieldMessage.Validate()
	if err == nil {
		t.Error("Validate: message with a missing field is marked as valid")
	}
}

func TestValidateFromField(t *testing.T) {
	var badFromMessage = Message{
		From:    "email",
		To:      "email@example.com",
		Subject: "subject",
		Text:    "Hello from tests",
	}

	err := badFromMessage.Validate()
	if err == nil {
		t.Error("Validate: message with bad *from* field is marked as valid")
	}
}

func TestValidateToField(t *testing.T) {
	var badToMessage = Message{
		From:    "email@example.com",
		To:      "email",
		Subject: "subject",
		Text:    "Hello from tests",
	}

	err := badToMessage.Validate()
	if err == nil {
		t.Error("Validate: message with bad *from* field is marked as valid")
	}
}

func TestToDict(t *testing.T) {
	mDict := validMessage.ToDict()
	dict := map[string]string{
		"from":    validMessage.From,
		"to":      validMessage.To,
		"subject": validMessage.Subject,
		"text":    validMessage.Text,
	}

	equal := reflect.DeepEqual(mDict, dict)
	if !equal {
		t.Errorf(
			"ToDict: returns a wrong dict representation of the message.\nGot: %+v Expected: %+v",
			mDict,
			dict,
		)
	}
}

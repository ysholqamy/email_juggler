package email

import (
	"fmt"
	"regexp"
)

var emailRegexp = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// Message represents a basic email message.
type Message struct {
	From    string // email address of sender
	To      string // email address of receiver
	Subject string // subject line
	Text    string // email body
}

// ToDict converts the message to a dictionary
func (m Message) ToDict() map[string]string {
	dict := map[string]string{
		"from":    m.From,
		"to":      m.To,
		"subject": m.Subject,
		"text":    m.Text,
	}

	return dict
}

// Validate From and To email addresses are wellformed.
func (m Message) Validate() error {
	mDict := m.ToDict()

	// validate format of From and To emails
	emailFields := []string{"from", "to"}
	for _, emailField := range emailFields {
		value := mDict[emailField]
		if !emailRegexp.MatchString(value) {
			return fmt.Errorf("field %s has an invalid value %s", emailField, value)
		}
	}

	// required fields are present and email fields are well formed.
	return nil
}

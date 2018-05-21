package email

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

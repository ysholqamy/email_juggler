package email

// Message represents a basic email message.
type Message struct {
	From    string // email address of sender
	To      string // email address of receiver
	Subject string // subject line
	Text    string // email body
}

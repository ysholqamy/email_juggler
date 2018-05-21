package email

// Provider is an abstraction over email providers
type Provider interface {
	// Send takes a message and sends it, returning errors if any.
	Send(Message) error
}

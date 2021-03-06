package email

import (
	"errors"
	"log"
	"sync"
)

// Common provider errors
var (
	ErrServiceUnavailable = errors.New("Service Unavailable, try again later")
)

// DefaultProvider is a roundrobin provider over default Mailgun and Sendgrid.
var DefaultProvider = NewRoundRobinProvider(DefaultMailgun, DefaultSendgrid)

// Provider is an abstraction over email providers
type Provider interface {
	// Send takes a message and sends it, returning errors if any.
	Send(Message) error
}

// roundrobinProvider represents an email provider
// routing emails to be sent to its subproviders
// in a roundrobin fashion.
type roundrobinProvider struct {
	*sync.Mutex
	subProviders []Provider
	turn         int
}

// NewRoundRobinProvider creates a roundrobin email provider.
// Requires at least a single subprovider to function properly.
func NewRoundRobinProvider(providers ...Provider) Provider {
	if len(providers) == 0 {
		return nil
	}

	return &roundrobinProvider{
		Mutex:        &sync.Mutex{},
		subProviders: providers,
		turn:         0,
	}
}

// Send implements the Provider Send method.
// All subproviders will be tried for a given message until one succeeds.
// Steps through the given subproviders in turn, doing at most a full cycle.
func (rr *roundrobinProvider) Send(m Message) error {
	// do a single cycle through subproviders
	// starting at the subprovider currently in turn
	for range rr.subProviders {
		err := rr.nextProvider().Send(m)

		// message sent successfully
		if err == nil {
			return nil
		}

		// logged for internal use, will not be returned to caller
		log.Println(err)
	}

	// all subproviders failed to send message
	return ErrServiceUnavailable
}

// nextProvider returns the provider currently in turn
// advances the turn pointer of the roundrobinProvider
// safe for concurrent use
func (rr *roundrobinProvider) nextProvider() Provider {
	rr.Lock()
	defer rr.Unlock()

	p := rr.subProviders[rr.turn]
	rr.turn = (rr.turn + 1) % len(rr.subProviders) // advances in ring

	return p
}

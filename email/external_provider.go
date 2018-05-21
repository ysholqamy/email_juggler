package email

import "net/http"

type externalProvider interface {
	// an external provider is also a provider. i.e implements Send.
	Provider

	// url of this external provider.
	url() string

	// adds the necessary headers/fields to authorize the request.
	authorize(*http.Request)

	// normalizes response errors, if any.
	normalizeErrors(*http.Response) error
}

// Mailgun represents Mailgun external email provider
type Mailgun struct {
	Key     string
	Domain  string
	BaseURL string
}

// Sendgrid represents Sendgrid external email provider
type Sendgrid struct {
	Key     string
	BaseURL string
}

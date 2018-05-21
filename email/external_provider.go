package email

import (
	"net/http"
	"net/url"
	"strings"
)

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

func generateFormRequest(URL string, m Message) (*http.Request, error) {
	// encode message as a form
	mDict := m.ToDict()
	form := url.Values{}
	for key, val := range mDict {
		form.Set(key, val)
	}

	// create request
	req, err := http.NewRequest("POST", URL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func processMessage(ep externalProvider, m Message) error {
	// generate request
	req, err := generateFormRequest(ep.url(), m)
	if err != nil {
		return err
	}

	// authorization header
	ep.authorize(req)

	// execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// Provider specific error occured
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		return ep.normalizeErrors(res)
	}

	// message sent successfully.
	return nil
}

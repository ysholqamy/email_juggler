package email

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Common external provider errors
var (
	ErrProviderBadResponse = errors.New("Could not parse provider response")
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

// Send implements the Provider Send method.
func (mg *Mailgun) Send(m Message) error        { return processMessage(mg, m) }
func (mg *Mailgun) url() string                 { return mg.BaseURL + "/" + mg.Domain + "/messages" }
func (mg *Mailgun) authorize(req *http.Request) { req.SetBasicAuth("api", mg.Key) }
func (mg *Mailgun) name() string                { return "Mailgun" }

// Handles Mailgun specific response errors
func (mg *Mailgun) normalizeErrors(res *http.Response) error {
	// parseBody
	body, err := parseJSONBody(res)
	if err != nil {
		return errors.New(ErrProviderBadResponse.Error() + ": " + err.Error())
	}

	// extract error message
	errMessage, ok := body["message"].(string)
	if !ok {
		return ErrProviderBadResponse
	}

	return errors.New(errMessage)
}

// Sendgrid represents Sendgrid external email provider
type Sendgrid struct {
	Key     string
	BaseURL string
}

// Send implements the Provider Send method.
func (sg *Sendgrid) Send(m Message) error        { return processMessage(sg, m) }
func (sg *Sendgrid) url() string                 { return sg.BaseURL }
func (sg *Sendgrid) authorize(req *http.Request) { req.Header.Set("Authorization", "Bearer "+sg.Key) }
func (sg *Sendgrid) name() string                { return "Sendgrid" }

// Handles Sendgrid specific response errors
func (sg *Sendgrid) normalizeErrors(res *http.Response) error {
	// parseBody
	body, err := parseJSONBody(res)
	if err != nil {
		return errors.New(ErrProviderBadResponse.Error() + ": " + err.Error())
	}

	// extract error message
	errs, ok := body["errors"].([]interface{})
	if !ok {
		return ErrProviderBadResponse
	}

	// convert error messages to strings and join them into a single message.
	errMessages := make([]string, len(errs))
	for i, v := range errs {
		errMessages[i] = fmt.Sprint(v)
	}
	errMessage := strings.Join(errMessages, ". ")

	return errors.New(errMessage)
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

// Attempts to parse the response body as JSON
func parseJSONBody(res *http.Response) (map[string]interface{}, error) {
	var body map[string]interface{}

	bodyBuffer, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBuffer, &body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

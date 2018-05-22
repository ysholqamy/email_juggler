package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	contentJSON       = "application/json"
	contentURLEncoded = "application/x-www-form-urlencoded"
)

// DefaultService represents the default provider wrapped into an HTTP service.
var DefaultService = CreateService(DefaultProvider)

// CreateService Wraps a Provider into an http Handler.
// Enforces http Method and Content-Type.
// Handles body parsing to message.
// Delegates sending message to Provider.
func CreateService(p Provider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// only handle POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "only POST is supported", http.StatusMethodNotAllowed)
			return
		}

		// only handle JSON and URLEncoded requests
		if r.Header.Get("Content-Type") != contentJSON &&
			r.Header.Get("Content-Type") != contentURLEncoded {
			http.Error(w, "supports "+contentJSON+" and "+contentURLEncoded+" only",
				http.StatusUnsupportedMediaType)
			return
		}

		// parse body to message
		message, err := parseBodyMessage(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// validate message
		if err = message.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// delegate sending message to provider
		err = p.Send(message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "message sent successfully\n")
	})
}

func parseBodyMessage(r *http.Request) (Message, error) {
	if r.Header.Get("Content-Type") == contentJSON {
		return parseJSONMessage(r)
	}

	return parseURLEncodedMessage(r)
}

func parseURLEncodedMessage(r *http.Request) (Message, error) {
	m := Message{}

	err := r.ParseForm()
	if err != nil {
		return m, err
	}

	m.To = r.Form.Get("to")
	m.From = r.Form.Get("from")
	m.Text = r.Form.Get("text")
	m.Subject = r.Form.Get("subject")

	return m, nil
}

func parseJSONMessage(r *http.Request) (Message, error) {
	var m Message

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return m, err
	}

	err = json.Unmarshal(body, &m)
	if err != nil {
		return m, err
	}

	return m, err
}

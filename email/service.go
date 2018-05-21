package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	contentJSON = "application/json"
)

// CreateService Wraps a Provider into an http Handler.
// Handles body parsing to message.
// Delegates sending message to Provider.
func CreateService(p Provider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// only handle POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST is supported", http.StatusMethodNotAllowed)
			return
		}

		// only handle JSON requests
		if r.Header.Get("Content-Type") != contentJSON {
			http.Error(w, "Supports "+contentJSON+" only", http.StatusUnsupportedMediaType)
			return
		}

		// parse body to message
		message, err := parseJSONMessage(r)
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
		fmt.Fprintf(w, "Message sent successfully.\n")
	})
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

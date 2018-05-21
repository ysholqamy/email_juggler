package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// CreateService Wraps a Provider into an http Handler.
// Handles body parsing to message.
// Delegates sending message to Provider.
func CreateService(p Provider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// no errors, message was sent.
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

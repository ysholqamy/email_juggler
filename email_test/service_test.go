package email_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/ysholqamy/email_juggler/email"
)

// Wrapping a roundrobin provider that will always succeed into a service
var mockService = CreateService(NewRoundRobinProvider(&MockProvider{Succeed: true}))

func mockServeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.Handler(mockService)
	handler.ServeHTTP(rr, req)

	return rr
}

func TestValidRequest(t *testing.T) {
	jsonBody, err := json.Marshal(validMessage)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/email", bytes.NewBuffer(jsonBody))

	if err != nil {
		t.Fatal(err)
	}

	res := mockServeRequest(req)

	if res.Code != http.StatusOK {
		t.Errorf("Failed to process valid message. got: %d", res.Code)
	}
}

func TestNoAvailableProviders(t *testing.T) {
	jsonBody, err := json.Marshal(validMessage)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/email", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	// A roundrobin provided over a single always failing provider.
	service := CreateService(NewRoundRobinProvider(&MockProvider{Succeed: false}))
	res := httptest.NewRecorder()
	handler := http.Handler(service)
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected service unavailable, got: %d", res.Code)
	}
}

// Propagate invalid message errors as bad request
func TestBadMessage(t *testing.T) {
	jsonBody, err := json.Marshal(missingFieldMessage)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/email", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	res := mockServeRequest(req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("Malformed message proccessed. got: %d", res.Code)
	}
}

// Handle POST only
func TestWrongVerb(t *testing.T) {
	req, err := http.NewRequest("GET", "/email", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	res := mockServeRequest(req)
	if res.Code != http.StatusMethodNotAllowed {
		t.Error("Wrong method allowed")
	}
}

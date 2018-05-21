package email_test

import (
	"testing"

	. "github.com/ysholqamy/email_juggler/email"
)

type MockProvider struct {
	Trials  int  // how many times send was called
	Succeed bool // given a valid message, should always succeed or not?
}

func (mp *MockProvider) Send(m Message) error {
	mp.Trials++
	if !mp.Succeed {
		return ErrServiceUnavailable
	}

	return nil
}

var (
	mockP1 = &MockProvider{Succeed: false}
	mockP2 = &MockProvider{Succeed: false}
)

var mocks = []*MockProvider{mockP1, mockP2}
var rrMockProvider = NewRoundRobinProvider(mockP1, mockP2)

func resetAllMocks() {
	for _, mock := range mocks {
		mock.Trials = 0
		mock.Succeed = false
	}
	rrMockProvider = NewRoundRobinProvider(mockP1, mockP2)
}

func TestAllFail(t *testing.T) {
	resetAllMocks()
	err := rrMockProvider.Send(validMessage)

	// assert it failed
	if err == nil {
		t.Errorf("Message was not sent and no error was reported")
	}

	// all tried
	for i, mock := range mocks {
		if mock.Trials == 0 {
			t.Errorf("Failed without trying provider number %d", i)
		}
	}

	// correct err reported
	if err != ErrServiceUnavailable {
		t.Errorf("All subproviders failed. Expected err: %v, found: %v.", ErrServiceUnavailable, err)
	}
}

func TestStopsOnSuccess(t *testing.T) {
	resetAllMocks()
	mockP1.Succeed = true
	err := rrMockProvider.Send(validMessage)

	// assert it succeeds
	if err != nil {
		t.Errorf("Message was sent and an error was reported. err: %v", err)
	}

	// only first tried
	if mockP2.Trials != 0 {
		t.Errorf("first provider succeeded. Attempted to send message using second provider")
	}
}

// Only try second provider after first providerfails
func TestFailsOver(t *testing.T) {
	resetAllMocks()
	mockP2.Succeed = true
	err := rrMockProvider.Send(validMessage)

	// assert it succeeds
	if err != nil {
		t.Errorf("Message was sent and an error was reported. err: %v", err)
	}

	// tried first provider
	if mockP1.Trials == 0 {
		t.Error("Failed over without attempting first provider")
	}

	// tried second provider
	if mockP2.Trials == 0 {
		t.Error("Did not attempt second provider when provider first fails")
	}
}

func TestRoundRobin(t *testing.T) {
	resetAllMocks()
	mockP1.Succeed = true
	mockP2.Succeed = true
	rrMockProvider.Send(validMessage)
	rrMockProvider.Send(validMessage)

	if mockP1.Trials != 1 || mockP2.Trials != 1 {
		t.Error("Messages are not routed in a roundrobin fashion")
	}
}

package twilio

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("foo", "bar")

	if client.sid != "foo" {
		t.Errorf("Client sid was incorrect, got: %v, want: %v.", client.sid, "foo")
	}

	if client.token != "bar" {
		t.Errorf("Client sid was incorrect, got: %v, want: %v.", client.token, "bar")
	}
}
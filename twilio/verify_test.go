package twilio

import (
	"net/url"
	"testing"
)

func TestVerifySignature(t *testing.T) {
	client := NewClient("", "12345")

	u := "https://mycompany.com/myapp.php?foo=1&bar=2"

	p := make(url.Values)
	p.Add("CallSid", "CA1234567890ABCDE")
	p.Add("Caller", "+12349013030")
	p.Add("Digits", "1234")
	p.Add("From", "+12349013030")
	p.Add("To", "+18005551212")

	s := "0/KCTR6DLpKmkAf8muzZqo1nDgQ="

	if !client.VerifySignature(s, u, p) {
		t.Errorf("Client sid was incorrect, got: %v, want: %v.", "false", "true")
	}
}
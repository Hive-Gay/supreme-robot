package webapp

import (
	"fmt"
	"net/http"
)

func (s *Server)HandleWebhookSMSPost(w http.ResponseWriter, r *http.Request) {
	// Get Signature from Header
	signature := ""
	idempotencyToken := ""

	for name, values := range r.Header {
		if name == "X-Twilio-Signature" {
			signature = values[0]
		}
		if name == "I-Twilio-Idempotency-Token" {
			idempotencyToken = values[0]
		}
	}

	// If no signature send unauthorized
	if signature == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// If no idempotencyToken send bad request
	if idempotencyToken == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// parse form data
	err := r.ParseForm()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Loop over header names
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println(name, value)
		}
	}

	// Write
	_ = fmt.Sprintf("https://%s%s", s.apphostname, r.URL.String())

	w.WriteHeader(http.StatusAccepted)
}
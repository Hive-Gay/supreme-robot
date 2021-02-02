package webapp

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func (s *Server)HandleWebhookSMSPost(w http.ResponseWriter, r *http.Request) {
	// make URL
	uri := fmt.Sprintf("https://%s%s", s.apphostname, r.URL.String())

	// parse form data
	err := r.ParseForm()
	if err != nil {
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

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
	if !s.twilioClient.VerifySignature(signature, uri,  r.Form) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// If no idempotencyToken send bad request
	if idempotencyToken == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	_, err = uuid.Parse(idempotencyToken)
	if err != nil {
		logger.Debugf("couldn't parse uuid: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// parse form data
	err = r.ParseForm()
	if err != nil {
		logger.Debugf("couldn't parse form: %s", err.Error())
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	paramsJson, err := json.Marshal(r.Form)
	if err != nil {
		logger.Debugf("couldn't marshal json: %s", err.Error())
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Enqueue
	err = s.enqueuer.ReceivedSMS(string(paramsJson), idempotencyToken)
	if err != nil {
		logger.Warningf("couldn't enqueue sms: %s", err.Error())
		s.returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
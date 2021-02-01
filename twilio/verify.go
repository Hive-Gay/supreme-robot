package twilio

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"sort"
)

func (c *Client)VerifySignature(signature, url string, params url.Values) bool {
	// Add URL to signature string
	sigString := url

	// Alphabetize Keys
	keys := make([]string, 0)
	for key, _ := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Add to signature string
	for _, key := range keys {
		values := params[key]
		sigString = sigString + key + values[0]
	}

	// Calculate signature
	keyForSign := []byte(c.token)
	h := hmac.New(sha1.New, keyForSign)
	h.Write([]byte(sigString))
	calculatedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature == calculatedSignature
}
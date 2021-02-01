package twilio

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
)

func (c *Client)VerifySignature(signature, url string, params url.Values) bool {
	sigString := url

	fmt.Println()
	keys := make([]string, 0)
	for key, _ := range params {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		values := params[key]
		fmt.Println(key, values[0])
		sigString = sigString + key + values[0]
	}

	fmt.Println()

	// Calculate Signature
	key_for_sign := []byte(c.token)
	h := hmac.New(sha1.New, key_for_sign)
	h.Write([]byte(sigString))
	calculatedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	fmt.Println("NewSig: ", calculatedSignature)


	return signature == calculatedSignature
}
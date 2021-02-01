package webapp

import (
	"fmt"
	"net/http"
)

func HandleWebhookSMSPost(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://%s%s", apphostname, r.URL.String())

	// parse form data
	err := r.ParseForm()
	if err != nil {
		returnErrorPage(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// Get Signature from Header
	sendSignature := ""
	for name, values := range r.Header {
		if name == "X-Twilio-Signature" {
			sendSignature = values[0]
		}
	}

	fmt.Println(twilioClient.VerifySignature(sendSignature, url, r.Form))

}
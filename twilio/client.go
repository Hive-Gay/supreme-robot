package twilio

import "github.com/juju/loggo"

type Client struct {
	sid      string
	token    string

	logger *loggo.Logger
}

func NewClient(sid, token string) *Client {
	logger := loggo.GetLogger("twilio")

	return &Client{
		sid:      sid,
		token:    token,

		logger: &logger,
	}

}

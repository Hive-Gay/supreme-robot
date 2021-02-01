package twilio

import "github.com/juju/loggo"

type Client struct {
	sid      string
	token    string
}

var logger = loggo.GetLogger("twilio")

func NewClient(sid, token string) *Client {
	logger.Debugf("creating new twilio client")
	return &Client{
		sid:      sid,
		token:    token,
	}

}

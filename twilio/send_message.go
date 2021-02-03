package twilio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type SendMessageResponse struct {
	AccountSid        string  `json:"account_sid"`
	ApiVersion        string  `json:"api_version"`
	Body              string  `json:"body"`
	Direction         string  `json:"direction"`
	DateCreated       string  `json:"date_created"`
	DateSent          *string `json:"date_sent"`
	DateUpdated       string  `json:"date_updated"`
	ErrorCode         *string `json:"error_code"`
	ErrorMessage      *string `json:"error_message"`
	From              string  `json:"from"`
	MessageServiceSid *string `json:"messaging_service_sid"`
	NumMedia          string  `json:"num_media"`
	NumSegments       string  `json:"num_segments"`
	Price             *string `json:"price"`
	PriceUnit         *string `json:"price_unit"`
	Sid               string  `json:"sid"`
	Status            string  `json:"status"`
	To                string  `json:"to"`
}

func (c *Client) SendMessage(from, to, body, callback string) (*SendMessageResponse, error) {
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + c.sid + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("From", from)
	msgData.Set("To", to)
	msgData.Set("Body", body)
	if callback != "" {
		msgData.Set("StatusCallback", callback)
	}
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(c.sid, c.token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("couldn't send message request to twilio: %s", err.Error())
		return nil, err
	}

	var data SendMessageResponse
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			logger.Errorf("couldn't decode json: %s", err.Error())
			return nil, errors.New(fmt.Sprintf("couldn't decode json: %s", err.Error()))
		} else {
			logger.Tracef("twilio data: %#v\n", data)
		}

	} else {
		return nil, errors.New(fmt.Sprintf("twilio returning non 200 status: %d", resp.StatusCode))
	}

	return &data, nil
}

package twilio

import (
	"github.com/twilio/twilio-go"
	"sync"
	"user_service/internal/config"
)

type Client struct {
	TwilioClient *twilio.RestClient
	ServiceSID   string
	FromPhone    string
}

var instance *Client
var once sync.Once

func GetClient() *Client {
	once.Do(func() {
		cfg := config.GetConfig()
		username := cfg.Twilio.AccountSID
		password := cfg.Twilio.AuthToken
		instance = &Client{
			TwilioClient: twilio.NewRestClientWithParams(twilio.ClientParams{
				Username: username,
				Password: password,
			}),
			ServiceSID: cfg.Twilio.ServiceSID,
			FromPhone:  cfg.Twilio.TwilioPhone,
		}
	})
	return instance
}

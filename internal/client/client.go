package client

import (
	"net/http"

	"github.com/NarthurN/GoXML_JSON/settings"
)

type Client struct {
	URL    string
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		URL: settings.ClientURL,
		client: &http.Client{
			Timeout: settings.ClientTimeout,
		},
	}
}

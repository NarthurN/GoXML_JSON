package client

import (
	"net/http"

	"github.com/NarthurN/GoXML_JSON/pkg/logger"
	"github.com/NarthurN/GoXML_JSON/settings"
)

type Client struct {
	URL    string
	logger *logger.Logger
	client *http.Client
}

func NewClient(logger *logger.Logger) *Client {
	return &Client{
		URL: settings.ClientURL,
		logger: logger,
		client: &http.Client{
			Timeout: settings.ClientTimeout,
		},
	}
}

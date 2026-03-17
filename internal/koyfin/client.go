package koyfin

import (
	"entext-applications/internal/utils"
	"net/http"
	"net/http/cookiejar"
)

// Client provides a wrapper around http.Client from which http requests will be run
type Client struct {
	client *http.Client
}

func NewClient() *Client {
	client := utils.NewClient()
	jar, _ := cookiejar.New(nil)
	client.Jar = jar

	return &Client{
		client: client,
	}
}

type KoyfinError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

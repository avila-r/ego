package httpx

import "net/http"

var (
	DefaultClient = &Client{
		client: http.DefaultClient,
	}
)

type Client struct {
	client *http.Client
}

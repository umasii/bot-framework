package client

import (
	"net/url"

	"github.com/cicadaaio/httpclient/net/http"
)

type Response struct {
	Headers    http.Header
	Body       string
	Status     string
	StatusCode int
	Url        *url.URL
}

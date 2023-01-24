package Client

import (
	"github.com/cicadaaio/httpclient/net/http"
	"net/url"
)

type Response struct {
	Headers    http.Header
	Body       string
	Status     string
	StatusCode int
	Url        *url.URL
}

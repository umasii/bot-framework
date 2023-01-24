package Client

import (
	"io"

	"github.com/cicadaaio/httpclient/net/http"
)

type Client struct {
	Client  *http.Client
	LastUrl string
}

func (c *Client) NewRequest() *Request {
	return &Request{
		Client: c,
	}
}

func (c *Client) Do(r *http.Request) (*Response, error) {

	resp, err := c.Client.Do(r)
	if err != nil {
		
		return nil, err
	}
	response := &Response{
		Headers:    resp.Header,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Url:        resp.Request.URL,
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	response.Body = string(body)
	
	r.Close = true

	return response, nil
}

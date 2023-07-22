package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"strings"

	"github.com/cicadaaio/httpclient/net/http"
	tls "github.com/cicadaaio/utls"
)

type Request struct {
	Client                 *Client
	Headers                []map[string]string
	Method                 string
	Url                    string
	Body                   io.Reader
	Host                   string
}

func (r *Request) SetJSONBody(body interface{}) {
	b, _ := json.Marshal(body)
	r.Body = bytes.NewBuffer(b)
}

func (r *Request) SetFormBody(body url.Values) {
	r.Body = strings.NewReader(body.Encode())
}

func (r *Request) SetHost(Host string) {
	r.Host = Host
}

func (r *Request) Do() (*Response, error) {
	req, err := http.NewRequest(r.Method, r.Url, r.Body)
	if err != nil {
		return nil, err
	}
	var Headers = map[string][]string{}
	var Order = []string{}
	var spoof bool

	for _, header := range r.Headers {
		for key, val := range header {
			Headers[key] = []string{val}
			if key != ":authority" {
				Order = append(Order, key)
			} else {
				spoof = true
			}
		}
	}
	Order = append(Order, "Cookie", "Content-Length")
	if r.Client.Client.Transport.(*http.Transport).ClientHelloID != &tls.HelloRandomizedNoALPN {
		Headers["HEADERORDER"] = Order
		if spoof {
			Headers["PSEUDOORDER"] = []string{
				":method",
				":authority",
				":scheme",
				":path",
			}

		}
	}
	req.Header = Headers
	r.Client.LastUrl = r.Url
	req.Host = r.Host
	req.Close = true
	resp, err := r.Client.Do(req)
	
	return resp, err
}

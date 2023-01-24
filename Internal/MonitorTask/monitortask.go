package MonitorTask

import (
	"github.com/cicadaaio/LVBot/CMD/DataStores/ProxyStore"
	Client3 "github.com/cicadaaio/LVBot/Internal/Client"
	Monitor "github.com/cicadaaio/LVBot/Internal/MonitorEngine"

	//"github.com/cicadaaio/LVBot/Internal/MonitorEngine"
	"net/url"
	"time"

	"github.com/cicadaaio/LVBot/Internal/Proxies"
	"github.com/cicadaaio/httpclient/net/http"
	tls "github.com/cicadaaio/utls"
	utls "github.com/cicadaaio/utls"
	"github.com/gwatts/rootcerts"

	"github.com/cicadaaio/httpclient/net/http/cookiejar"
)

type RunningStatus string

const (
	Active  RunningStatus = "Active"
	Paused                = "Paused"
	Stopped               = "Stopped"
)

type MonitorTask struct {
	Client      Client3.Client `json:"-"`
	Jar         *cookiejar.Jar `json:"-"`
	Proxy       Proxies.Proxy  `json:"-"`
	ProxyList   []string       `json:"-"`
	Delay       time.Duration
	StockStatus *Monitor.MonitorResp
	Running     RunningStatus
}

func (t *MonitorTask) WaitDelay() {
	time.Sleep(t.Delay * time.Millisecond)
}

func (t *MonitorTask) InitializeClient() {
	var err error
	for {
		t.Jar, err = cookiejar.New(nil)
		if err != nil {
			continue
		}
		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: false,
			Proxy:              http.ProxyURL(nil),
			ClientHelloID:      &tls.HelloChrome_83,
			ForceAttemptHTTP2:  true,
		}
		tr.TLSClientConfig = &utls.Config{
			InsecureSkipVerify: false,
			RootCAs:            rootcerts.ServerCertPool(),
		}
		t.Client.Client = &http.Client{Transport: tr, Jar: t.Jar}
		if err != nil {
			continue
		}
		return
	}
}

func (t *MonitorTask) RotateProxy() {
	if t.ProxyList == nil || (t.ProxyList)[0] == "" {
		return
	} else {
		for {
			if t.Proxy.Raw != "" {
				Proxies.SafeProxy.ReleaseProxy(t.Proxy.Raw)
			}

			rawProxy := Proxies.SafeProxy.GetProxy(t.ProxyList)

			// Proxy is empty wait some time then get new
			if rawProxy == "" {
				t.WaitDelay()
				continue
			}

			parsedProxy, err := Proxies.SafeProxy.ParseProxy(rawProxy)
			if err != nil {
				continue
			}

			t.Proxy = parsedProxy

			t.Client.Client.Transport.(*http.Transport).Proxy = http.ProxyURL(t.Proxy.Formatted)
			t.Client.Client.Transport.(*http.Transport).CloseIdleConnections()

			return
		}
	}
}

func (t *MonitorTask) Charles() {
	charles, _ := url.Parse("http://localhost:8888")
	t.Client.Client.Transport.(*http.Transport).Proxy = http.ProxyURL(charles)
	t.Client.Client.Transport.(*http.Transport).CloseIdleConnections()

}

func (m *MonitorTask) Initialize() {

	m.ProxyList = ProxyStore.GetMonitoringProxies()
	m.InitializeClient()
	m.RotateProxy()
}

func (m *MonitorTask) Start(checkStock func()) {
	m.Running = Active

	for {

		m.WaitDelay()

		if m.Running == Stopped {
			return
		}

		if m.Running == Paused {
			continue
		}

		checkStock()
	}

}

func (m *MonitorTask) Stop() {
	m.Running = Stopped
}

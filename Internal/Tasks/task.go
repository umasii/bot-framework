package tasks

import (
	"context"
	goErrors "errors"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/umasii/bot-framework/cmd/datastores/profilestore"
	"github.com/umasii/bot-framework/cmd/datastores/proxystore"
	"github.com/umasii/bot-framework/internal/activityapi"
	"github.com/umasii/bot-framework/internal/certs"
	"github.com/umasii/bot-framework/internal/client"
	"github.com/umasii/bot-framework/internal/errors"
	"github.com/umasii/bot-framework/internal/helpers/utilities/cardidentification"
	"github.com/umasii/bot-framework/internal/products"
	"github.com/umasii/bot-framework/internal/profiles"
	"github.com/umasii/bot-framework/internal/proxies"

	tls "github.com/cicadaaio/utls"

	"sync"
	"time"

	"github.com/cicadaaio/httpclient/net/http"
)

type Stage string

const (
	Start   Stage = "Start"
	Stop          = "Stop"
	InStock       = "InStock"
	Success       = "Success"
)

type Task struct {
	Stage         Stage              `json:"-"`
	Client        client.Client     `json:"-"`
	Proxy         proxies.Proxy      `json:"-"`
	ProxyList     []string           `json:"-"`
	Status        string             `json:"-"`
	Profile       *profiles.Profile  `json:"-"`
	Jar           *client.FJar      `json:"-"`
	Tries         int                `json:"-"`
	Ctx           context.Context    `json:"-"`
	Cancel        context.CancelFunc `json:"-"`
	DeclineReason string             `json:"-"`
	OrderId       string             `json:"-"`

	TaskID         int
	GroupID        int
	GroupName      string
	ProfileGroupID int
	ProfileID      int
	ProxyGroupID   int
	ProxyGroupName string
	Site           string
	Product        products.Product
	MonitorDelay   time.Duration
	RetryDelay     time.Duration
}

type IBotTask interface {
	Get() *Task
	Initialize()
	InjectTaskData()
	//Start(wg *sync.WaitGroup)
	Stop()
	Execute()
	WrapExecutor(f func(), wg *sync.WaitGroup)
}

type TaskGroup struct {
	GroupName string
	GroupID   int
	Tasks     []IBotTask
}

func (t Task) Get() *Task {
	return &t

}

func (t *Task) Stop() {
	t.Cancel()
}

func (t *Task) WrapExecutor(f func(), wg *sync.WaitGroup) {
	stage := Start

	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				log.Println(err)
			}
		}
	}()
	defer wg.Done()

	t.Initialize()
	for {
		select {
		case <-t.Ctx.Done():
			panic(goErrors.New("stopped"))

		default:
			f()
			if stage == t.Stage {
				t.Tries++
				//handle tries in the methods themselves
				continue
			}
			if t.Stage == Success {
				//TODO: Webhook func here that access t.Product
				return
			}
			t.Tries = 0
			stage = t.Stage
		}
	}
}

func (t *Task) Restart() {
	t.UpdateStatus("Restarting", activityapi.LogLevel)
	t.Stage = Start
	t.Initialize()
}

func (t *Task) UpdateStatus(status string, level string) {

	t.Status = status
	fmt.Println(fmt.Sprintf("Task ID: %s | Status: %s | Product: %s | Status: %s", strconv.Itoa(t.TaskID), t.Status, t.Product.Identifier, t.Status))
}

func (t Task) WaitM() {
	time.Sleep(t.MonitorDelay * time.Millisecond)
}

func (t Task) WaitR() {
	time.Sleep(t.RetryDelay * time.Millisecond)
}

func (t *Task) Charles() {
	charles, _ := url.Parse("http://localhost:8888")
	t.Client.Client.Transport.(*http.Transport).Proxy = http.ProxyURL(charles)
	t.Client.Client.Transport.(*http.Transport).CloseIdleConnections()

}

func (t *Task) RotateProxy() {
	//support localhost through nil proxylist
	if t.ProxyList == nil || (t.ProxyList)[0] == "" {
		return
	} else {
		for {
			if t.Proxy.Raw != "" {
				proxies.SafeProxy.ReleaseProxy(t.Proxy.Raw)
			}

			rawProxy := proxies.SafeProxy.GetProxy((t.ProxyList))

			// Proxy is empty wait some time then get new
			if rawProxy == "" {
				t.UpdateStatus("Waiting for Proxy", activityapi.WarningLevel)
				t.WaitR()
				continue
			}

			parsedProxy, err := proxies.SafeProxy.ParseProxy(rawProxy)
			if err != nil {
				t.UpdateStatus(err.Error()+", retrying", activityapi.ErrorLevel)
				continue
			}

			t.Proxy = parsedProxy
			t.Client.Client.Transport.(*http.Transport).Proxy = http.ProxyURL(t.Proxy.Formatted)
			t.Client.Client.Transport.(*http.Transport).CloseIdleConnections()
			return
		}
	}
}

func (t *Task) InitializeClient() {
	var err error
	for {
		t.Jar = client.New()
		if err != nil {
			continue
		}
		tr := &http.Transport{
			MaxIdleConns:          5,
			MaxConnsPerHost:       5,
			IdleConnTimeout:       60 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
			DisableCompression:    false,
			Proxy:                 http.ProxyURL(nil),
			ClientHelloID:         &tls.HelloChrome_99,
			ForceAttemptHTTP2:     true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				RootCAs:            certs.ServerCertPool(),
			},
		}
		t.Client.Client = &http.Client{Transport: tr, Jar: t.Jar}
		if err != nil {
			continue
		}

		return
	}
}

func (t *Task) Initialize() {
	Ctx, Cancel := context.WithCancel(context.Background())
	t.Ctx = Ctx
	t.Cancel = Cancel

	t.InjectTaskData()
	t.InitializeClient()
	//t.RotateProxy()
	t.Stage = Start
}

func (t *Task) InjectTaskData() {
	for {
		profileData, err := profilestore.GetProfileByID((*t).ProfileGroupID, (*t).ProfileID)
		proxyGroup, err := proxystore.GetProxyGroupByID((*t).ProxyGroupID)
		proxy, err := proxies.SafeProxy.ParseProxy(proxies.SafeProxy.GetProxy(proxyGroup.Proxies))

		if err != nil {
			continue
		}

		(*t).Profile = profileData

		(*t).Profile.Billing.CardType, err = cardidentification.CreditCardType((*t).Profile)

		if err != nil {
			errors.Handler(err)
			return
		}

		t.ProxyGroupName = proxyGroup.GroupName
		proxy.GroupID = (*t).ProxyGroupID
		proxy.GroupName = (*t).ProxyGroupName
		(*t).Proxy = proxy
		(*t).ProxyList = (proxyGroup.Proxies)
		return
	}
}

func (t *Task) SendCheckoutData(checkoutStatus bool, checkoutResponse *client.Response, additionalInfo interface{}) {
	req := t.Client.NewRequest()
	req.Url = "http://127.0.0.1:3000/activity/"
	payload := activityapi.RecpData{
		UserInfo: activityapi.UserData{
			UserID:        "TESTING", // TODO: get this from bot
			ShippingState: t.Profile.Shipping.State,
		},
		Settings: activityapi.TaskSettings{
			Site:    t.Site,
			Product: t.Product.Identifier,
			Mode:    "Mode one",
		},
		Results: activityapi.TaskResults{
			CheckedOut:         checkoutStatus,
			CheckoutStatusCode: checkoutResponse.StatusCode,
			CheckoutMessage:    checkoutResponse.Body,
		},
		Instance: activityapi.InstanceInfo{
			OS:                   "Mac",
			TotalTasks:           0,
			TotalTasksForProduct: 0,
			Time:                 time.Now().Unix(),
		},
	}

	if additionalInfo != nil {
		payload.AdditionalData = additionalInfo
	}
	req.SetJSONBody(payload)

	resp, err := req.Do()

	if err != nil {
		errors.Handler(err)
	}

	if resp.StatusCode != 200 {
		errors.Handler(goErrors.New("Failed to send checkout stat!, resp code" + resp.Status))
	}
}

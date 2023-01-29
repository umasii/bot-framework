package walmart

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cicadaaio/LVBot/Internal/ActivityApi"
	"github.com/cicadaaio/LVBot/Internal/Client"
	Errors "github.com/cicadaaio/LVBot/Internal/Errors"
	"github.com/cicadaaio/LVBot/Internal/Tasks"

	api2captcha "github.com/cicadaaio/2cap"

	WmMonitor "github.com/cicadaaio/LVBot/Monitors/walmart"

	uuid "github.com/PrismAIO/go.uuid"

	encrypt "github.com/cicadaaio/walmart-encryption"

	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cicadaaio/httpclient/net/http"
)

type WalmartTask struct {
	Tasks.Task
	Px            PxInfo `json:"-"`
	Mode          WmMonitor.Mode
	ItemId        string `json:"-"`
	UserAgent     string `json:"-"` // Set specifically for WM as it requires PX cookies which are UA bound
	secCookieWg   *sync.WaitGroup
	checkoutRetry int         `json:"-"`
	PrevStage     Tasks.Stage `json:"-"`

	StoreData SavedShippingRates `json:"-"`

	PieKeyId1     string    `json:"-"`
	PieKeyId2     string    `json:"-"`
	PiePhase1     string    `json:"-"`
	PiePhase2     string    `json:"-"`
	CardData1     [3]string `json:"-"` //TODO: Need to check if objects are actually strings
	CardData2     [3]string `json:"-"`
	PiHash        string    `json:"-"`
	StockStatus   bool      `json:"-"`
	CaptchaClient *api2captcha.Client
	CaptchaCount  int
	Language      string
	SechUa        string
	Accept        string
}

func (e *WalmartTask) Execute() {
	if e.CaptchaCount > 3 && e.Stage != Cart {
		e.RotateProxy()
		return
	}
	switch e.Stage {
	case Tasks.Start:
		var err error

		if err != nil {
			fmt.Println(err)
		}
		e.Language = randomAcceptLanguage()
		e.Accept = randomAccept()
		e.SechUa = randomChUa()

		e.setWmCookie()

		e.Mode = WmMonitor.ModeTwo
		e.secCookieWg = &sync.WaitGroup{}
		e.pxLoop()

	case Account:
		e.Stage = ShippingRates
		return
		e.genAccount()

	case ShippingRates:
		e.genShippingRates()

	case GetCart:
		e.getCart()

	case PreEncrypt:
		e.preEncrypt()

	case CC:
		e.Stage = Cart
		return
		e.submitCC()

	case Cart:
		e.cart()

	case Tasks.InStock:
		e.startCheckout()

	case Checkout:

		e.startCheckout()

	case Shipping:
		e.submitShipping()

	case Payment:
		e.submitPayment()
	case Order:
		e.submitOrder()

	case DeleteItem:
		e.deleteItem()
	}
}

func (e *WalmartTask) AutoSolve() string {
	e.UpdateStatus("Solving captcha...", ActivityApi.LogLevel)

	cap := api2captcha.ReCaptcha{
		Version:   "RecaptchaV3TaskProxyless",
		SiteKey:   "6Lc8-RIaAAAAAPWSm2FVTyBg-Zkz2UjsWWfrkgYN",
		Url:       "https://www.walmart.com/cart",
		Invisible: false,
		Action:    "handleCaptcha",
		Score:     0.9,
	}

	for {
		req := cap.ToRequest()

		code, err := e.CaptchaClient.Solve(req)

		if err == nil {
			e.UpdateStatus("SOLVED Captcha "+code, ActivityApi.LogLevel)
			return code
		} else {
			Errors.Handler(err)
		}

	}
}
func (e *WalmartTask) pxLoop() {

	pxWait := make(chan bool) // We do not want to progress until the initial PX cookie has been set, so a channel is set after the first loop
	i := false

	go func() {
		for {
			if i == true {
				time.Sleep(270 * time.Second)
			}

			err := e.getPxCookies()

			if err == nil {

				pxWait <- true
				i = true

				continue
			} else {

				Errors.Handler(err)

				continue
			}
		}
	}()

	<-pxWait

	e.Stage = Account

}

func (e *WalmartTask) getPxCookies() error {

	uuid := string(hex.EncodeToString(uuid.NewV1().Bytes()))
	uuidstr := uuid[:8] + "-" + uuid[8:12] + "-" + uuid[12:16] + "-" + uuid[16:20] + "-" + uuid[20:]
	e.UpdateStatus("Getting PX cookies", ActivityApi.LogLevel)

	pxPacket := PxInfo{
		Uuid:  &uuidstr,
		Site:  "walmart",
		Proxy: e.Proxy.Formatted.String(),
		Key:   "CDDF7993D938682E9D319592E79EC"}

	req := e.Client.NewRequest()
	req.Method = "POST"

	req.Url = "https://px.cicadabots.com/px"

	req.SetJSONBody(pxPacket)

	req.Headers = randomizeHeaders([]map[string]string{
		{"content-type": "application/json"},
	})
	resp, err := req.Do()
	if err != nil {
		
		return errors.New("failed to send px request")
	}

	if resp.StatusCode == 200 {

		var pxApiResp pxResp

		err = json.Unmarshal([]byte(resp.Body), &pxApiResp)
		if err != nil {

			return errors.New("failed to process px data")
		}

		e.setPxCookie(pxApiResp.Px3)

		e.UserAgent = pxApiResp.Ua

		e.Px.SetID = &pxApiResp.SetID
		e.Px.Uuid = &pxApiResp.Uuid
		e.Px.Vid = &pxApiResp.Vid

		e.UpdateStatus("Successfully genned PX cookies", ActivityApi.LogLevel)
		return nil

	} else {

		e.UpdateStatus("Failed to gen PX cookie", ActivityApi.ErrorLevel)
		
		e.WaitR()
		return errors.New("failed to gen px cookie")
	}

}

func (e *WalmartTask) PxCheck(resp *Client.Response) bool {

	if resp.StatusCode == 200 || resp.StatusCode == 201 { // I figure we should do this first thing so we don't waste time doing ioutil on fine requests
		return false
	}

	if (resp.StatusCode == 412 || resp.StatusCode == 307) && strings.Contains(string(resp.Body), "blocked") {
		e.UpdateStatus("PX Captcha blocked  ["+strconv.Itoa(resp.StatusCode)+"]", ActivityApi.ErrorLevel)
		e.CaptchaCount++
		e.RotateProxy()
		e.getPxCookies()

	}

	return true
}

func (e *WalmartTask) pxCap( /*captchaToken string*/) error {

	e.UpdateStatus("Getting PX captcha cookies", ActivityApi.LogLevel)

	e.Px.Site = "walmart"
	e.Px.Proxy = e.Proxy.Formatted.String()
	e.Px.Key = "CDDF7993D938682E9D319592E79EC"
	pxPacket := e.Px

	req := e.Client.NewRequest()

	req.Method = "POST"
	req.SetJSONBody(pxPacket)
	req.Url = "https://px.cicadabots.com/px"
	req.Headers = randomizeHeaders([]map[string]string{

		{"content-type": "application/json"},
	})

	resp, err := req.Do()

	if err != nil {

		return errors.New("failed to send PX Captcha request")
	}

	if resp.StatusCode == 200 {

		var pxApiResp pxResp

		err = json.Unmarshal([]byte(resp.Body), &pxApiResp)
		if err != nil {

			Errors.Handler(err)

			return errors.New("failed to process PX Captcha request")
		}

		e.setPxCookie(pxApiResp.Px3)

		e.UpdateStatus("Successfully solved PX captcha", ActivityApi.LogLevel)
		return nil

	} else {

		e.UpdateStatus("Failed to process PX Captcha data", ActivityApi.ErrorLevel)

		e.WaitR()
		return errors.New("")
	}

}

func (e *WalmartTask) setPxCookie(px3 string) {
	wmUrl, _ := url.Parse("https://www.walmart.com/")

	px3Cookie := &http.Cookie{
		Name:     "_px3",
		Value:    px3,
		MaxAge:   999999,
		Secure:   true,
		HttpOnly: true,
	}

	var newCookies []*http.Cookie //
	newCookies = append(newCookies, px3Cookie)

	e.Jar.SetCookies(wmUrl, newCookies)

}

func (e *WalmartTask) genShippingRates() {

	shippingInfoToSend := ShippingData{
		e.Profile.Shipping.PostalCode,
		"STOREMETAPLUS",
		true,
		true,
		"Web-Checkout-ShippingAddress",
		true,
		true,
	}

	req := e.Client.NewRequest()
	req.Url = "https://www.walmart.com/account/api/location"
	req.Method = "PUT"
	req.SetJSONBody(shippingInfoToSend)

	req.Headers = randomizeHeaders([]map[string]string{
		{"sec-ch-ua": e.SechUa},
		{"Accept": e.Accept},
		{"sec-ch-ua-mobile": "?0"},
		{"user-agent": e.UserAgent},
		{"Content-Type": "application/json"},
		{"origin": "https://www.walmart.com"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-dest": "empty"},
		{"referer": "https://www.walmart.com/checkout/"},
		{"Accept-Language": e.Language},
	})

	resp, err := req.Do()

	if err != nil {
		e.UpdateStatus("Failed to send shipping rates request", ActivityApi.ErrorLevel)
		Errors.Handler(err)
		return
	}

	shippingReqResp := ShippingResp{}

	err = json.Unmarshal([]byte(resp.Body), &shippingReqResp)

	if err != nil {
		e.UpdateStatus("Failed to unpack shipping rate response", ActivityApi.ErrorLevel)
		Errors.Handler(err)
		e.WaitR()
		return
	}

	if len(shippingReqResp.Stores) == 0 {
		e.UpdateStatus("Failed to get valid shipping rates response... Rotating proxy", ActivityApi.ErrorLevel) // I've found this happens on bad ips for whatever reason
		e.RotateProxy()
		e.WaitR()
		return
	}

	var goodStore string

	for _, store := range shippingReqResp.Stores {
		if strings.Contains(store.Types[0], "gsf_store") {
			goodStore = store.StoreId
		}
	}

	if goodStore == "" {
		e.UpdateStatus("Failed to get valid shipping rates... Please make sure your shipping address is supported by US Walmart.", ActivityApi.ErrorLevel)
		e.WaitR()
		return // We need to stop this task here
	}

	storeList := SavedStoreList{goodStore}

	e.StoreData = SavedShippingRates{
		[]SavedStoreList{storeList},
		e.Profile.Shipping.PostalCode,
		e.Profile.Shipping.City,
		e.Profile.Shipping.State,
		true,
		"",
		"",
		"",
		"",
	}

	e.UpdateStatus("Successfully genned shipping rates", ActivityApi.LogLevel)
	e.Stage = GetCart

}

func (e *WalmartTask) genAccount() {

	e.UpdateStatus("Genning Walmart account", ActivityApi.LogLevel)
	emailSplit := strings.Split(e.Profile.Email, "@")
	rand.Seed(time.Now().UnixNano())

	accountEmail := fmt.Sprintf("%s+%d@%s", emailSplit[0], rand.Intn(100000), emailSplit[1])

	capToSend := signupCaptcha{""}

	Person := personName{
		e.Profile.Shipping.FirstName,
		e.Profile.Shipping.LastName,
	}
	accountToSend := accountData{
		Person,
		accountEmail,
		"TempPassword",
		false,
		"true",
		false,
		capToSend,
	}

	req := e.Client.NewRequest()

	req.Method = "POST"
	req.Url = "https://www.walmart.com/account/electrode/api/signup?ref=domain"
	req.SetJSONBody(accountToSend)

	req.Headers = []map[string]string{
		{"content-length": ""},
		{"sec-ch-ua": `" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`},
		{"sec-ch-ua-mobile": "?0"},
		{"user-agent": e.UserAgent},
		{"content-type": "application/json; charset=utf-8"},
		{"accept": "*/*"},
		{"origin": "https://www.walmart.com"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-dest": "empty"},
		{"referer": "https://www.walmart.com/account/signup?ref=domain"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
	}

	resp, err := req.Do()

	if err != nil {
		Errors.Handler(err)
		e.UpdateStatus("Failed to send account request", ActivityApi.ErrorLevel)
		e.WaitR()
		return
	}

	if e.PxCheck(resp) == false {
		e.Stage = ShippingRates
		e.UpdateStatus("Successfully genned account", ActivityApi.LogLevel)
	} else {
		Errors.Handler(err)

		e.UpdateStatus("Failed to gen account", ActivityApi.ErrorLevel)
		e.WaitR()
		return
	}

}

func (e *WalmartTask) getCart() {
	e.UpdateStatus("Getting payment cookies..", ActivityApi.LogLevel)

	req := e.Client.NewRequest()
	req.Url = "https://www.walmart.com/cart"
	req.Method = "GET"

	req.Headers = []map[string]string{

		{"accept-encoding": "gzip, deflate, br"},
		{"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		{"accept-language": "en-US,en;q=0.9"},
		{"cache-control": "max-age=0"},
		{"dnt": "1"},
		{"sec-ch-ua": `" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-fetch-dest": "document"},
		{"sec-fetch-mode": "navigate"},
		{"sec-fetch-site": "none"},
		{"sec-fetch-user": "?1"},
		{"service-worker-navigation-preload": "true"},
		{"upgrade-insecure-requests": "1"},
		{"user-agent": e.UserAgent},
	}

	resp, err := req.Do()

	if err != nil {
		Errors.Handler(err)
		e.WaitR()
		return
	}

	if e.PxCheck(resp) == false {
		e.Stage = PreEncrypt
		e.UpdateStatus("Successfully got payment cookies", ActivityApi.LogLevel)
	}
	return

}

func (e *WalmartTask) getCheckout(wg *sync.WaitGroup) {
	e.UpdateStatus("Getting security cookies...", ActivityApi.LogLevel)
	req := e.Client.NewRequest()
	req.Url = "https://www.walmart.com/checkout"
	req.Method = "GET"
	req.Body = nil

	req.Headers = []map[string]string{

		{"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		{"accept-language": "en-US,en;q=0.9"},
		{"cache-control": "no-cache"},
		{"dnt": "1"},
		{"pragma": "no-cache"},
		{"referer": "https://www.walmart.com"},
		{"sec-ch-ua": `" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-fetch-dest": "document"},
		{"sec-fetch-mode": "navigate"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-user": "?1"},
		{"service-worker-navigation-preload": "true"},
		{"upgrade-insecure-requests": "1"},
		{"user-agent": e.UserAgent},
	}

	resp, err := req.Do()

	if err != nil {
		Errors.Handler(err)
		e.WaitR()
		return
	}

	if e.PxCheck(resp) == false {

		wg.Done()
		e.UpdateStatus("Successfully got security cookies!", ActivityApi.LogLevel)
		return

	} else {

		return
	}

}

func (e *WalmartTask) setWmCookie() {
	wmUrl, _ := url.Parse("https://www.walmart.com/")

	wmValue := fmt.Sprintf("reflectorid:0000000000000000000000@lastupd:%d@firstcreate:%d", time.Now().Unix(), time.Now().Unix())

	px3Cookie := &http.Cookie{
		Name:     "com.wm.reflector",
		Value:    wmValue,
		MaxAge:   999999,
		Secure:   true,
		HttpOnly: true,
	}
	gCookie := &http.Cookie{
		Name:     "g",
		Value:    "0",
		MaxAge:   999999,
		Secure:   true,
		HttpOnly: true,
	}

	var newCookies []*http.Cookie //
	newCookies = append(newCookies, px3Cookie)
	newCookies = append(newCookies, gCookie)

	e.Jar.SetCookies(wmUrl, newCookies)

}

func (e *WalmartTask) getPieKeys() (error, PieKeyResp) {

	req := e.Client.NewRequest()
	req.Url = "https://securedataweb.walmart.com/pie/v1/wmcom_us_vtg_pie/getkey.js?bust=" + strconv.FormatInt(time.Now().UnixNano(), 10)

	req.Method = "GET"

	req.Headers = []map[string]string{
		{"accept": "application/json"},
		{"accept-language": "en-US,en;q=0.9"},
		{"cache-control": "no-cache"},
		{"content-type": "application/json"},
		{"dnt": "1"},
		{"origin": "https://www.walmart.com"},
		{"pragma": "no-cache"},
		{"referer": "https://www.walmart.com/checkout/"},
		{"sec-ch-ua": `" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-origin"},
		{"user-agent": e.UserAgent},
	}

	resp, err := req.Do()

	if err != nil {

		Errors.Handler(err)

		return errors.New("failed to send payment keys request"), PieKeyResp{}
		e.WaitR()
	}

	res := PieKeyResp{}

	kRegx := regexp.MustCompile(`PIE.K = "(.*?)";`) // This repetive regex matching is very ugly, does anyone know a better way?
	kMatch := kRegx.FindStringSubmatch(resp.Body)
	res.K = kMatch[1]

	LRegx := regexp.MustCompile(`PIE.L = (\d*?);`)
	LMatch := LRegx.FindStringSubmatch(resp.Body)
	res.L, _ = strconv.Atoi(LMatch[1])

	eRegx := regexp.MustCompile(`PIE.E = (\d*?);`)
	eMatch := eRegx.FindStringSubmatch(resp.Body)
	res.E, _ = strconv.Atoi(eMatch[1])

	keyRegx := regexp.MustCompile(`PIE.key_id = "(.*?)";`)
	keyMatch := keyRegx.FindStringSubmatch(resp.Body)
	res.KeyId = keyMatch[1]

	phaseRegx := regexp.MustCompile(`PIE.phase = (.*?);`)
	phaseMatch := phaseRegx.FindStringSubmatch(resp.Body)
	res.PhaseId, _ = strconv.Atoi(phaseMatch[1])

	return nil, res

}

func (e *WalmartTask) preEncrypt() {
	e.UpdateStatus("Encrypting payment", ActivityApi.LogLevel)

	for i := 1; i <= 2; i++ {
		err, res := e.getPieKeys()

		if err != nil {
			Errors.Handler(err)
			return
		}

		PIE := encrypt.Pie{
			L: res.L,
			E: res.E,
			K: res.K,
		}

		if i == 1 {
			cardData := encrypt.ProtectPanAndCvv("4111111111111111", e.Profile.Billing.CVC, PIE)
			e.PieKeyId2 = res.KeyId
			e.PiePhase2 = strconv.Itoa(res.PhaseId)
			e.CardData2 = cardData

		} else if i == 2 {
			cardData := encrypt.ProtectPanAndCvv(e.Profile.Billing.CardNumber, e.Profile.Billing.CVC, PIE)

			e.PieKeyId1 = res.KeyId
			e.PiePhase1 = strconv.Itoa(res.PhaseId)
			e.CardData1 = cardData

		}

	}

	e.UpdateStatus("Successfully encrypted payment", ActivityApi.LogLevel)

	e.Stage = CC

}

func (e *WalmartTask) submitCC() {

	e.UpdateStatus("Submitting first round of payment", ActivityApi.LogLevel)

	submitCC := ccdata{
		e.Profile.Shipping.Address,
		e.Profile.Shipping.Address2,
		strings.ToUpper(e.Profile.Billing.CardType),
		e.CardData1[1],
		e.CardData1[0],
		e.Profile.Billing.ExpMonth,
		e.Profile.Billing.ExpYear,
		e.Profile.Billing.BillingAddress.FirstName,
		e.CardData1[2],
		true,
		e.PieKeyId1,
		e.Profile.Billing.BillingAddress.LastName,
		e.PiePhase1,
		e.Profile.Billing.BillingAddress.Phone,
		e.Profile.Billing.BillingAddress.PostalCode,
		e.Profile.Billing.BillingAddress.State,
	}
	req := e.Client.NewRequest()
	req.Method = "POST"
	req.SetJSONBody(submitCC)
	req.Url = "https://www.walmart.com/api/checkout-customer/:CID/credit-card"
	req.Headers = []map[string]string{
		{"content-length": ""},
		{"sec-ch-ua": `" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`},
		{"accept": "application/json"},
		{"sec-ch-ua-mobile": "?0"},
		{"user-agent": e.UserAgent},
		{"content-type": "application/json"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-dest": "empty"},
		{"referer": "https://www.walmart.com/checkout/"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
	}

	resp, err := req.Do()

	if err != nil {
		Errors.Handler(err)
		e.UpdateStatus("Failed to send round 1 payment request", ActivityApi.ErrorLevel)
		e.WaitR()

		return
	}

	if e.PxCheck(resp) == false {

		e.UpdateStatus("Successfully submitted round 1 payment", ActivityApi.LogLevel)
		var ccResp PaymentResp

		err = json.Unmarshal([]byte(resp.Body), &ccResp)

		e.PiHash = ccResp.PiHash

		e.Stage = Cart

	} else {

		e.UpdateStatus("Failed to submit round 1 payment", ActivityApi.ErrorLevel)

		e.WaitR()
		return
	}

}

func (e *WalmartTask) preAdd() {

	addToListData := atcParams{e.Product.Identifier, e.Product.Qty}

	e.UpdateStatus("Pre-adding product", ActivityApi.LogLevel)

	req := e.Client.NewRequest()
	req.Url = "https://www.walmart.com/api/v3/cart/:CRT/items"
	req.Method = "POST"
	req.SetJSONBody(addToListData)

	req.Headers = []map[string]string{

		{"content-length": ""},
		{"sec-ch-ua": `" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`},
		{"accept": "application/json"},
		{"sec-ch-ua-mobile": "?0"},
		{"user-agent": e.UserAgent},
		{"content-type": "application/json"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-dest": "empty"},
		{"referer": "https://www.walmart.com/checkout/"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
	}

	resp, err := req.Do()

	if err != nil {
		Errors.Handler(err)
		e.UpdateStatus("Failed to send pre-add request", ActivityApi.ErrorLevel)

		e.WaitR()

		return
	}

	if e.PxCheck(resp) == false {

		preAddResp := preAdd{}
		err = json.Unmarshal([]byte(resp.Body), &preAddResp)
		e.ItemId = preAddResp.SavedItems[0].ItemId
		e.UpdateStatus("Successfully pre-added product", ActivityApi.LogLevel)
		e.Stage = Cart
	} else {

		e.UpdateStatus("Failed to pre-add product", ActivityApi.ErrorLevel)
		if e.Tries > 3 {
			e.Stop()
			return
		}
		e.WaitR()

		return
	}

}

func (e *WalmartTask) cart() {

	e.UpdateStatus("Adding to cart", ActivityApi.LogLevel)

	atcData := atcParams{e.Product.Identifier, 1}

	req := e.Client.NewRequest()
	req.Url = "https://www.walmart.com/api/v3/cart/:CID/items"
	req.SetJSONBody(atcData)
	req.Method = "POST"
	req.Headers = []map[string]string{
		{"content-length": ""},
		{"sec-ch-ua": `" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`},
		{"accept": "application/json"},
		{"sec-ch-ua-mobile": "?0"},
		{"user-agent": e.UserAgent},
		{"content-type": "application/json"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-dest": "empty"},
		{"referer": "https://www.walmart.com/checkout/"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
	}

	resp, err := req.Do()
	if err != nil {
		Errors.Handler(err)
		e.UpdateStatus("Failed to send pre-cart request", ActivityApi.ErrorLevel)
		e.WaitR()
		return
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {

		var atcJsonResp AtcResp

		err = json.Unmarshal([]byte(resp.Body), &atcJsonResp)

		if err != nil {
			e.UpdateStatus("Failed to get valid cart resp...  ["+strconv.Itoa(resp.StatusCode)+"]", ActivityApi.ErrorLevel)
			e.WaitR()
			return
		}

		e.ItemId = atcJsonResp.Items[0].Id

		if int(atcJsonResp.Cart.Totals.SubTotal) > e.Product.PriceCheck {
			e.UpdateStatus("Added to cart but exceeded price limit! ["+strconv.Itoa(resp.StatusCode)+"]", ActivityApi.ErrorLevel)

			e.PrevStage = Cart
			e.Stage = DeleteItem
			return
		}

		e.PrevStage = Cart
		e.Stage = Checkout

		e.UpdateStatus("Successfully pre Carted  ["+strconv.Itoa(resp.StatusCode)+"]", ActivityApi.LogLevel)

	} else if strings.Contains(resp.Body, "ITEM_COUNT_MAX_LIMIT") {

		e.Product.Qty = 1 // TODO: Should this be an option? IE: if the bot can't add the qty a user specified, should we set the qty to 1 or just not go for it?
		e.UpdateStatus("Failed to add to cart due to item limit, retrying with qty 1 ["+strconv.Itoa(resp.StatusCode)+"]", ActivityApi.ErrorLevel)
		e.WaitR()
		return

	} else if strings.Contains(resp.Body, "No fulfillment option has availability") {

		e.UpdateStatus("Product is out of stock ["+strconv.Itoa(resp.StatusCode)+"]", ActivityApi.ErrorLevel)
		e.WaitR()
		return

	} else if e.PxCheck(resp) {
		e.UpdateStatus("Failed to add to cart... ["+strconv.Itoa(resp.StatusCode)+"]", ActivityApi.ErrorLevel)

		e.WaitR()
		return
	}

}

func (e *WalmartTask) startCheckout() {

	e.UpdateStatus("Starting checkout", ActivityApi.LogLevel)
	req := e.Client.NewRequest()
	req.SetJSONBody(e.StoreData)
	req.Method = "POST"
	req.Url = "https://www.walmart.com/api/checkout/v3/contract?page=CHECKOUT_VIEW"

	req.Headers = []map[string]string{
		{"accept-encoding": "gzip, deflate, br"},
		{"accept": "application/json, text/javascript, */*; q=0.01"},
		{"accept-language": "en-US,en;q=0.9"},
		{"content-type": "application/json; charset=UTF-8"},
		{"host": "www.walmart.com"},
		{"origin": "https://www.walmart.com"},
		{"referer": "https://www.walmart.com/checkout/"},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-origin"},
		{"user-agent": e.UserAgent},
		{"wm_cvv_in_session": "true"},
		{"wm_vertical_id": "0"},
		{"content-length": ""},
		{"sec-ch-ua": `" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`},
		{"accept": "application/json"},
		{"sec-ch-ua-mobile": "?0"},
		{"user-agent": e.UserAgent},
		{"content-type": "application/json"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-dest": "empty"},
		{"referer": "https://www.walmart.com/checkout/"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
	}

	resp, err := req.Do()

	if err != nil {
		Errors.Handler(err)
		e.WaitR()

		return
	}

	if e.PxCheck(resp) == false {

		var checkoutViewRespJson CheckoutViewResp

		err = json.Unmarshal([]byte(resp.Body), &checkoutViewRespJson)

		if err != nil {
			e.UpdateStatus("Failed to get valid checkout resp... ["+strconv.Itoa(resp.StatusCode)+"]", ActivityApi.ErrorLevel)
			e.WaitR()
			return
		}

		e.UpdateStatus("Successfully loaded checkout ["+strconv.Itoa(resp.StatusCode)+"]", ActivityApi.LogLevel)

		if int(checkoutViewRespJson.Summary.SubTotal) > e.Product.PriceCheck {
			e.UpdateStatus("Got to checkout but exceeded price check! restarting session ["+strconv.Itoa(resp.StatusCode)+"]", ActivityApi.ErrorLevel)
			e.PrevStage = Checkout
			e.Stage = DeleteItem
			return

		}

		e.Stage = Shipping

		e.secCookieWg.Add(1)

		go e.getCheckout(e.secCookieWg)
		return

	} else {
		e.UpdateStatus("Failed to go to checkout... Retrying... ["+strconv.Itoa(resp.StatusCode)+"]", ActivityApi.ErrorLevel)
		e.WaitR()
		return
	}

}

func (e *WalmartTask) submitShipping() {

	e.UpdateStatus("Submitting shipping", ActivityApi.LogLevel)
	shippingData := shippingInfo{
		e.Profile.Shipping.Address,
		e.Profile.Shipping.Address2,
		e.Profile.Shipping.City,
		e.Profile.Shipping.FirstName,
		e.Profile.Shipping.LastName,
		e.Profile.Shipping.Phone,
		e.Profile.Email,
		false,
		e.Profile.Shipping.PostalCode,
		e.Profile.Shipping.State,
		e.Profile.Shipping.Country,
		"RESIDENTIAL",
		[]string{}, // TODO: Check if this actuall

	}

	req := e.Client.NewRequest()
	req.Url = "https://www.walmart.com/api/checkout/v3/contract/:PCID/shipping-address"
	req.Method = "POST"
	req.SetJSONBody(shippingData)

	req.Headers = []map[string]string{
		{"content-length": ""},
		{"sec-ch-ua": `" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`},
		{"accept": "application/json"},
		{"sec-ch-ua-mobile": "?0"},
		{"user-agent": e.UserAgent},
		{"content-type": "application/json"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-dest": "empty"},
		{"referer": "https://www.walmart.com/checkout/"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
	}

	resp, err := req.Do()

	if err != nil {
		Errors.Handler(err)
		e.UpdateStatus("Failed to send shipping request", ActivityApi.ErrorLevel)
		e.WaitR()
		return
	}

	if e.PxCheck(resp) == false {

		e.UpdateStatus("Successfully submitted shipping", ActivityApi.LogLevel)
		e.Stage = Payment

	} else if e.Tries == 3 {
		e.UpdateStatus("Failed to submit shipping 3 times... ", ActivityApi.ErrorLevel)
		e.WaitR()

	} else {
		e.UpdateStatus("Failed to submit shipping... Retrying...", ActivityApi.ErrorLevel)
		e.WaitR()
		return
	}

}

func (e *WalmartTask) submitPayment() {

	e.UpdateStatus("Submitting payment", ActivityApi.LogLevel)

	myPayment := payments{"CREDITCARD",
		strings.ToUpper(e.Profile.Billing.CardType),
		e.Profile.Billing.BillingAddress.LastName,
		e.Profile.Billing.BillingAddress.FirstName,
		e.Profile.Billing.BillingAddress.Address,
		e.Profile.Billing.BillingAddress.Address2,
		e.Profile.Billing.BillingAddress.City,
		e.Profile.Billing.BillingAddress.State,
		e.Profile.Billing.BillingAddress.PostalCode,
		e.Profile.Billing.ExpMonth,
		e.Profile.Billing.ExpYear,
		e.Profile.Email,
		e.Profile.Billing.BillingAddress.Phone,
		e.CardData1[0],
		e.CardData1[1],
		e.CardData1[2],
		e.PieKeyId1,
		e.PiePhase1,

		"",
	}
	paymentToSend := paymentInfo{[]payments{myPayment}, true}

	req := e.Client.NewRequest()
	req.Url = "https://www.walmart.com/api/checkout/v3/contract/:PCID/payment"
	req.Method = "POST"
	req.SetJSONBody(paymentToSend)

	req.Headers = []map[string]string{
		{"content-length": ""},
		{"sec-ch-ua": `" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`},
		{"accept": "application/json"},
		{"sec-ch-ua-mobile": "?0"},
		{"user-agent": e.UserAgent},
		{"content-type": "application/json"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-dest": "empty"},
		{"referer": "https://www.walmart.com/checkout/"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
	}
	resp, err := req.Do()

	if err != nil {
		Errors.Handler(err)
		e.UpdateStatus("Failed to send payment request", ActivityApi.ErrorLevel)
		e.WaitR()

		return
	}

	if e.PxCheck(resp) == false {

		e.UpdateStatus("Successfully submitted payment", ActivityApi.LogLevel)
		e.Stage = Order

	} else if e.Tries == 3 {

		e.UpdateStatus("Failed to submit payment 3 times... ", ActivityApi.ErrorLevel)
		e.WaitR()

	} else {
		e.UpdateStatus("Failed to submit payment.. Retrying..", ActivityApi.ErrorLevel)
		e.WaitR()
		return
	}

}

func (e *WalmartTask) submitOrder() {
	e.checkoutRetry = e.checkoutRetry + 1

	e.UpdateStatus("Submitting order", ActivityApi.LogLevel)

	paymentInfo := voltagePayments{
		"CREDITCARD",
		e.CardData1[1],
		e.CardData1[0],
		e.CardData1[2],
		e.PieKeyId1,
		e.PiePhase1,
	}

	OrderData := orderInfo{
		true,
		[]voltagePayments{paymentInfo},
	}

	req := e.Client.NewRequest()
	req.SetJSONBody(OrderData)
	req.Url = "https://www.walmart.com/api/checkout/v3/contract/:PCID/order"
	req.Method = "PUT"

	req.Headers = []map[string]string{
		{"content-length": ""},
		{"sec-ch-ua": `" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`},
		{"accept": "application/json"},
		{"sec-ch-ua-mobile": "?0"},
		{"user-agent": e.UserAgent},
		{"content-type": "application/json"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-dest": "empty"},
		{"referer": "https://www.walmart.com/checkout/"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
	}

	e.secCookieWg.Wait() // Waits until we have gotten cookies from getting /checkout

	resp, err := req.Do()

	if err != nil {
		e.UpdateStatus("Failed to send order request", ActivityApi.ErrorLevel)
		e.WaitR()

		return
	}

	if e.PxCheck(resp) == false {

		e.UpdateStatus("Successfully submitted order", ActivityApi.LogLevel)
		// e.FireWebHook(true, []*discordhook.EmbedField{{
		// 	Name:   "Mode",
		// 	Value:  string(e.Mode),
		// 	Inline: false,
		// }})
		e.Stop()

	} else {

		var orderFail OrderResp

		err = json.Unmarshal([]byte(resp.Body), &orderFail)

		e.UpdateStatus(fmt.Sprintf("Failed to submit order. Reason: %s . Attempt %d out of 3 [%d]", orderFail.FailedReason, e.checkoutRetry, resp.StatusCode), ActivityApi.ErrorLevel)
		e.WaitR()

		if e.Tries == 2 {
			// e.FireWebHook(false, []*discordhook.EmbedField{{ // you can pass additional fields if you want to
			// 	Name:   "Failed Reason",
			// 	Value:  orderFail.FailedReason,
			// 	Inline: false,
			// }, {
			// 	Name:   "Mode",
			// 	Value:  string(e.Mode),
			// 	Inline: false,
			// },
			// 	{ // you can pass additional fields if you want to
			// 		Name:   "Resp status code",
			// 		Value:  strconv.Itoa(resp.StatusCode),
			// 		Inline: false,
			// 	},
			// })

			if strings.Contains(orderFail.FailedReason, "Item is no longer in stock") {
				e.UpdateStatus("Failed to checkout due to out of stock error.", ActivityApi.ErrorLevel)
				e.WaitR()

			} else {
				e.Stop()

			}

		}
		return
	}

}

func (e *WalmartTask) deleteItem() {
	req := e.Client.NewRequest()
	req.Url = "https://www.walmart.com/api/v3/cart/:CRT/items/" + e.ItemId
	req.Method = "DELETE"
	req.Headers = []map[string]string{
		{"authority": "www.walmart.com"},
		{"pragma": "no-cache"},
		{"cache-control": "no-cache"},
		{"sec-ch-ua": `" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`},
		{"dnt": "1"},
		{"sec-ch-ua-mobile": "?0"},
		{"user-agent": e.UserAgent},
		{"credentials": "include"},
		{"content-type": "application/json"},
		{"accept": "application/json, text/javascript, */*; q=0.01"},
		{"omitcsrfjwt": "true"},
		{"origin": "https://www.walmart.com"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-dest": "empty"},
		{"referer": "https://www.walmart.com/cart?action=SignIn&rm=true"},
		{"accept-language": "en-US,en;q=0.9"},
	}
	resp, err := req.Do()

	if err != nil {
		e.UpdateStatus("Failed to send delete req", ActivityApi.ErrorLevel)
		Errors.Handler(err)
		return
	}

	if resp.StatusCode != 200 {
		e.UpdateStatus("Failed to delete item... retrying", ActivityApi.ErrorLevel)

		Errors.Handler(errors.New("Failed to delete item, resp:" + resp.Body))
		e.WaitR()
		return
	}

	e.UpdateStatus("Succesfully deleted item out of price range", ActivityApi.LogLevel)
	e.Stage = e.PrevStage
	return

}

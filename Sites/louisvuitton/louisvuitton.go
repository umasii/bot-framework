package louisvuitton

import (
	"bytes"
	//"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"net/url"

	"github.com/tidwall/gjson"
	"github.com/umasii/bot-framework/internal/activityapi"
	"github.com/umasii/bot-framework/internal/tasks"

	api2captcha "github.com/cicadaaio/2cap"
)

var lvUrl *url.URL

func checkPage(shtml string) bool {
	if strings.Contains(shtml, "Thank You for Your Interest") {
		return false
	} else {
		return true
	}
}

type LVTask struct {
	tasks.Task
	aVal string `json:"-"`
	uVal string `json:"-"`
	tVal string `json:"-"`

	ItemId        string      `json:"-"`
	SkuId         string      `json:"-"`
	CatalogRefId  string      `json:"-"`
	UserAgent     string      `json:"-"` // we used website UA rather than LV app UA
	checkoutRetry int32       `json:"-"`
	PrevStage     tasks.Stage `json:"-"`
	Region        string      `json:"-"` // for use on Asia/Europe drops (we only ran US though)
	Mode          string      `json:"-"`
	SensorData    string      `json:"-"`
	PixelData     string      `json:"-"`

	CaptchaClient *api2captcha.Client
	CaptchaCount  int

	Username string `json:"username"`
	Password string `json:"password"`

	SensorURL string `json:"-"`
	CardId    string `json:"-"`

	PaymentRetries int    `json:"-"`
	OrderRetries   int    `json:"-"`
	PixelScriptURL string `json:"-"`
	PixelURL       string `json:"-"`
}

func (e *LVTask) Execute() {

	e.SensorURL = "https://secure.louisvuitton.com/Q9qt2k/lHAJ/mu/aofS/G8ehiyCHsvg/3Q1tXbwSESG5/aVNaAQ/HU/gCBTFTHEo"

	// for the sake of time, account information was hardcoded for Litty2 mode
	// Litty1 mode does not require accounts and was to be used on a larger scale
	e.Username = "[redacted]"
	e.Password = "[redacted]"

	//e.UserAgent = "Mozilla/5.0 (Nintendo Switch; WebApplet) AppleWebKit/609.4 (KHTML, like Gecko) NF/6.0.2.20.5 NintendoBrowser/5.1.0.22023"
	e.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1320.0"

	e.Mode = "Litty2"

	// A nested switch was used to handle mode selection and each step of the task flow
	switch e.Mode {
	case "Litty2":
		switch e.Stage {
		case tasks.Start: // tasks.Start is part of the Cicada framework, but no initialization is required
			var err error

			if err != nil {
				fmt.Println(err)
			}
			e.Stage = GetSecure

		case GetSecure: // sets initial Akamai cookies
			e.getSecure()

		case GetU: // gets Akamai pixel
			e.getU()

		case GetSensor:
			e.getAkamaiSensor()

		case Akamai:
			e.submitSensor()

		case AkamaiPixel:
			e.submitPixel()

		case Account: // Litty2 mode used accounts with the Air Forces precarted already
			e.login()

		case PrepareCart:
			e.prepareCart()

		case GetCart:
			e.getCart()

		case GetCard:
			e.getCard()

		case SubmitPayment:
			e.submitPayment()

		case Order:
			e.order()
		}

	case "Litty1":
		switch e.Stage {
		case tasks.Start:

			var err error

			if err != nil {
				fmt.Println(err)
			}

			e.Stage = GetSecure

		case GetSecure:
			e.getSecure()

		case GetU:
			e.getU()

		case GetSensor:
			e.getAkamaiSensor()

		case Akamai:
			e.submitSensor()

		case AkamaiPixel:
			e.submitPixel()

		case ProductUrl:
			e.productUrl()

		case LittyATC: // this was the ATC bypass
			e.littyATC()

		case GetCart:
			e.getCart()

		case GetCard:
			e.getCard()

		case SubmitPayment:
			e.submitPayment()

		case Order:
			e.order()

		}
	}

}

// We chose to fetch product information through an individual LVTask instead of a monitor task for simplicity
// and due to there being 9 different products
func (e *LVTask) getProduct() {
	e.UpdateStatus("Fetching product", activityapi.LogLevel)

	req := e.Client.NewRequest()
	req.Url = "https://pass.louisvuitton.com/api/facade/api/eng-us/catalog/product/" + e.Product.Identifier // fetching through facade bypassed LV's blocking of web endpoints during drop time
	req.Method = "GET"
	req.Headers = []map[string]string{
		{"accept": "application/json, text/plain, */*"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
		{"origin": "https://us.louisvuitton.com"},
		{"referer": "https://us.louisvuitton.com/eng-us/new/for-men/louis-vuitton-and-nike-air-force-1-by-virgil-abloh/_/N-t1769r8n"},
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": `"Windows"`},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-site"},
		{"user-agent": e.UserAgent},
	}

	resp, err := req.Do()
	if err != nil {
		// handle error, set next stage
		return
	}

	if resp.StatusCode == 200 {
		productResp := string(resp.Body)
		e.Product.ProductName = gjson.Get(productResp, "name").String()
		e.Product.Qty = 1
		e.Product.Image = gjson.Get(productResp, "model.0.images.0.contentUrl").String()

		return
	} else {
		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

func (e *LVTask) getVariant() {
	e.UpdateStatus("Fetching random in-stock variant", activityapi.LogLevel)

	req := e.Client.NewRequest()
	req.Url = "https://pass.louisvuitton.com/api/facade/api/eng-us/catalog/product/" + e.Product.Identifier
	req.Method = "GET"
	req.Headers = []map[string]string{
		{"accept": "application/json, text/plain, */*"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
		{"origin": "https://us.louisvuitton.com"},
		{"referer": "https://us.louisvuitton.com/eng-us/new/for-men/louis-vuitton-and-nike-air-force-1-by-virgil-abloh/_/N-t1769r8n"},
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": `"Windows"`},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-site"},
		{"user-agent": e.UserAgent},
	}

	resp, err := req.Do()
	if err != nil {
		// handle error, set next stage
		return
	}

	if resp.StatusCode == 200 {
		productResp := string(resp.Body)
		e.Product.ProductName = gjson.Get(productResp, "name").String()
		e.Product.Qty = 1
		e.Product.Image = gjson.Get(productResp, "model.0.images.0.contentUrl").String()

		return
	} else {
		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

func (e *LVTask) getU() {

	e.UpdateStatus("Getting Pixel", activityapi.LogLevel)
	req := e.Client.NewRequest()
	req.Url = e.PixelScriptURL
	req.Method = "GET"
	req.Headers = []map[string]string{
		{"accept": "application/json, text/plain, */*"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
		{"origin": "https://us.louisvuitton.com"},
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": `"Windows"`},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-site"},
		{"user-agent": e.UserAgent},
	}

	resp, err := req.Do()
	if err != nil {
		// handle error, set next stage
		log.Println(err)
		e.UpdateStatus("Getting Pixel Error", activityapi.LogLevel)

		return
	}

	if resp.StatusCode != 200 {
		return
	}

	body := resp.Body

	var uVal string

	re := regexp.MustCompile("g=_\\[(\\d{1,3})\\]")
	matches := re.FindAllStringSubmatch(body, -1)

	if len(matches) == 0 || len(matches[0]) < 2 {
		return
	}

	index, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return
	}

	re = regexp.MustCompile("\\[(.*?)\\]")
	matchesTwo := re.FindStringSubmatch(body)

	if len(matchesTwo) == 0 {
		return
	}

	match := matchesTwo[0]

	split := strings.Split(match, ",")

	formatted := split[index]
	formatted = strings.ReplaceAll(formatted, `"`, "")

	splitU := strings.Split(formatted, "\\x")[1:]

	for _, s := range splitU {
		s = strings.Trim(s, " \r\t\n[]")

		p, err := strconv.ParseInt("0x"+s, 0, 0)
		if err != nil {
			return
		}

		uVal += string(rune(p))
	}
	e.uVal = uVal
	e.Stage = AkamaiPixel
	return
}

func (e *LVTask) genPixelData() {
	e.UpdateStatus("Submitting Pixel", activityapi.LogLevel)

	payload := genPixelData{
		"https://us.louisvuitton.com/",
		e.uVal,
		e.aVal,
		e.UserAgent,
		"13",
	}

	req := e.Client.NewRequest()
	req.Url = "[redacted]"
	req.Method = "POST"
	req.Headers = []map[string]string{
		{"x-api-key": AK_API_KEY},
	}
	req.SetJSONBody(payload)

	resp, err := req.Do()
	if err != nil {
		// handle error, set next stage
		return
	}

	if resp.StatusCode == 200 {
		sensorResp := string(resp.Body)
		e.PixelData = gjson.Get(sensorResp, "pixel_payload").String()

		e.Stage = AkamaiPixel
		return
	} else {
		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

func (e *LVTask) getAkamaiSensor() {
	e.UpdateStatus("Generating Akamai sensor", activityapi.LogLevel)

	cookies := e.Jar.Cookies(lvUrl)

	var ak_cookie string

	for _, cookie := range cookies {
		if cookie.Name == "_abck" {
			ak_cookie = cookie.Value
			break
		}
	}

	if ak_cookie == "" {
		e.getSecure()
		return
	}

	payload := genSensorData{
		e.UserAgent,
		"https://secure.louisvuitton.com/",
		ak_cookie,
	}

	req := e.Client.NewRequest()
	req.Url = "[redacted]"
	req.Method = "POST"
	req.Headers = []map[string]string{
		{"x-api-key": AK_API_KEY},
	}
	req.SetJSONBody(payload)

	resp, err := req.Do()
	if err != nil {
		// handle error, set next stage
		return
	}

	respBody := string(resp.Body)
	fmt.Println(respBody)
	if resp.StatusCode == 200 {
		sensorResp := string(resp.Body)
		e.SensorData = gjson.Get(sensorResp, "sensor_data").String()

		e.Stage = Akamai
		return

	} else {
		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

func (e *LVTask) getSecure() {
	e.UpdateStatus("Getting homepage", activityapi.LogLevel)

	req := e.Client.NewRequest()
	req.Url = e.SensorURL
	req.Method = "GET"
	req.Headers = []map[string]string{
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": `"Windows"`},
		{"upgrade-insecure-requests": "1"},
		{"user-agent": e.UserAgent},
		{"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "navigate"},
		{"sec-fetch-user": "?1"},
		{"sec-fetch-dest": "document"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
	}

	resp, err := req.Do()
	if err != nil {
		// handle error, set next stage
		return
	}

	if resp.StatusCode == 200 {
		lvUrl = resp.Url
		if e.Mode == "Litty1" {
			e.Stage = GetU
		} else {
			e.Stage = GetSensor
		}

		return
	}

	e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
	e.WaitR()

	return
}

func (e *LVTask) submitSensor() {
	e.UpdateStatus("Submitting Akamai sensor", activityapi.LogLevel)

	payload := sensorPayload{
		e.SensorData,
	}

	req := e.Client.NewRequest()
	req.Url = e.SensorURL
	req.Method = "POST"
	req.Headers = []map[string]string{
		{"accept": "*/*"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
		//{"content-length":""},
		{"content-type": "text/plain;charset=UTF-8"},
		{"origin": "https://secure.louisvuitton.com"},
		{"referer": "https://secure.louisvuitton.com/eng-us/mylv"},
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": `"Windows"`},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-site"},
		{"user-agent": e.UserAgent},
	}
	req.SetJSONBody(payload)

	resp, err := req.Do()
	if err != nil {
		return
	}

	respBody := string(resp.Body)
	fmt.Println(respBody)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		var ak_cookie string
		cookies := e.Jar.Cookies(lvUrl)
		for _, cookie := range cookies {
			if cookie.Name == "_abck" {
				ak_cookie = cookie.Value
				break
			}
		}

		if ak_cookie == "" {
			e.getSecure()
			e.getAkamaiSensor()
			e.Stage = Akamai
			return
		} else if strings.Contains(ak_cookie, "~0~") {
			if e.Mode == "Litty1" {
				e.Stage = ProductUrl

			} else {
				e.Stage = Account
			}

			return
		} else {
			e.getAkamaiSensor()
		}

		// set next stage
		return
	} else {
		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

func (e *LVTask) submitPixel() {
	e.UpdateStatus("Submitting Akamai pixel", activityapi.LogLevel)

	req := e.Client.NewRequest()
	req.Url = e.PixelURL
	req.Method = "POST"
	req.Headers = []map[string]string{
		{"accept": "*/*"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
		//{"content-length":""},
		{"content-type": "application/x-www-form-urlencoded"},
		{"origin": "https://secure.louisvuitton.com"},
		{"referer": "https://secure.louisvuitton.com/eng-us/mylv"},
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": `"Windows"`},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-site"},
		{"user-agent": e.UserAgent},
	}
	req.Body = bytes.NewBuffer([]byte(e.PixelData))

	resp, err := req.Do()
	if err != nil {
		return
	}

	respBody := string(resp.Body)
	fmt.Println(respBody)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		e.Stage = Akamai
		return
	} else {
		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

func (e *LVTask) login() {
	e.UpdateStatus("Logging into account", activityapi.LogLevel)

	payload := loginPayload{
		e.Username,
		e.Password,
	}

	req := e.Client.NewRequest()
	req.Url = "https://api.louisvuitton.com/api/eng-us/account/login/"
	req.Method = "POST"
	req.Headers = []map[string]string{
		{"accept": "application/json, text/plain, */*"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
		//{"content-length":""},
		{"content-type": "application/json"},
		{"origin": "https://us.louisvuitton.com"},
		{"referer": "https://us.louisvuitton.com/eng-us/homepage"},
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": `"Windows"`},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-site"},
		{"user-agent": e.UserAgent},
	}
	req.SetJSONBody(payload)

	resp, err := req.Do()
	if err != nil {
		return
	}

	respBody := string(resp.Body)
	fmt.Println(respBody)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		ok := checkPage(resp.Body)
		if ok {
			e.Stage = PrepareCart
			return
		} else {
			e.WaitR()
			return
		}
		// set next stage
		return
	} else if resp.StatusCode == 403 {
		e.getAkamaiSensor()
		e.submitSensor()
		e.Stage = Account
		return

	} else {
		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

func (e *LVTask) prepareCart() {
	e.UpdateStatus("Preparing cart", activityapi.LogLevel)

	req := e.Client.NewRequest()
	req.Url = "https://api.louisvuitton.com/api/eng-us/cart/prepare-order"
	req.Method = "PUT"
	req.Headers = []map[string]string{
		{"accept": "application/json, text/plain, */*"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
		//{"content-length":""},
		{"content-type": "application/json"},
		{"origin": "https://us.louisvuitton.com"},
		{"referer": "https://us.louisvuitton.com/eng-us/cart"},
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": `"Windows"`},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-site"},
		{"user-agent": e.UserAgent},
	}
	req.Body = strings.NewReader("{}")

	resp, err := req.Do()
	if err != nil {
		return
	}

	respBody := string(resp.Body)
	fmt.Println(respBody)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		ok := checkPage(resp.Body)
		if ok {
			e.Stage = GetCart
			return
		} else {
			e.WaitR()

			return
		}

	} else if resp.StatusCode == 403 {
		e.getAkamaiSensor()
		e.submitSensor()
		e.Stage = PrepareCart
		return

	} else {
		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

func (e *LVTask) getCard() {
	e.UpdateStatus("Fetching payment method", activityapi.LogLevel)

	req := e.Client.NewRequest()
	req.Url = "https://api.louisvuitton.com/api/eng-us/account/credit-cards"
	req.Method = "GET"
	req.Headers = []map[string]string{
		{"accept": "application/json, text/plain, */*"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
		{"content-type": "application/json"},
		{"origin": "https://us.louisvuitton.com"},
		{"referer": "https://us.louisvuitton.com/eng-us/checkout/payment"},
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": `"Windows"`},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-site"},
		{"user-agent": e.UserAgent},
	}

	resp, err := req.Do()
	if err != nil {
		return
	}

	respBody := string(resp.Body)
	fmt.Println(respBody)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		cardResp := string(resp.Body)
		e.CardId = gjson.Get(cardResp, "creditCardList.0.creditCardName").String()
		e.Stage = SubmitPayment
		return

	} else if resp.StatusCode == 403 {
		e.getAkamaiSensor()
		e.submitSensor()
		e.Stage = GetCard
		return

	} else {
		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

func (e *LVTask) getCart() {
	e.UpdateStatus("Fetching cart", activityapi.LogLevel)

	req := e.Client.NewRequest()
	req.Url = "https://api.louisvuitton.com/api/eng-us/cart/full"
	req.Method = "GET"
	req.Headers = []map[string]string{
		{"accept": "application/json, text/plain, */*"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
		{"content-type": "application/json"},
		{"origin": "https://us.louisvuitton.com"},
		{"referer": "https://us.louisvuitton.com/eng-us/checkout"},
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": `"Windows"`},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-site"},
		{"user-agent": e.UserAgent},
	}

	resp, err := req.Do()
	if err != nil {
		return
	}

	respBody := string(resp.Body)
	fmt.Println(respBody)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		ok := checkPage(resp.Body)
		if ok {
			e.Stage = GetCard
			return
		} else {
			e.WaitR()

			return
		}

	} else if resp.StatusCode == 403 {
		e.getAkamaiSensor()
		e.submitSensor()
		e.Stage = GetCart
		return

	} else {
		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

func (e *LVTask) productUrl() {
	e.UpdateStatus("Fetching product URL", activityapi.LogLevel)

	req := e.Client.NewRequest()
	req.Url = "https://us.louisvuitton.com/eng-us/products/lv-trainer-sneaker-nvprod3710063v/1AAHS3"
	req.Method = "GET"
	req.Headers = []map[string]string{
		{"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
		{"content-type": "application/json"},
		{"origin": "https://us.louisvuitton.com"},
		{"referer": "https://us.louisvuitton.com/eng-us/homepage"},
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": `"Windows"`},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-site"},
		{"user-agent": e.UserAgent},
	}

	resp, err := req.Do()
	if err != nil {
		return
	}

	respBody := string(resp.Body)
	fmt.Println(respBody)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		e.Stage = LittyATC

		return

	} else if resp.StatusCode == 403 {
		e.getAkamaiSensor()
		e.submitSensor()
		e.Stage = ProductUrl
		return

	} else {
		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

func (e *LVTask) submitPayment() {
	e.UpdateStatus("Submitting payment method", activityapi.LogLevel)

	var cardType string

	if e.Profile.Billing.CardType == "amex" {
		cardType = "americanExpress"
	} else if e.Profile.Billing.CardType == "mastercard" {
		cardType = "masterCard"

	} else if e.Profile.Billing.CardType == "visa" {
		cardType = "visa"
	} else {
		cardType = strings.ToLower(e.Profile.Billing.CardType)
	}

	createStruct := Create{
		cardType,
		e.CardId,
		e.Profile.Billing.CVC,
		true,
		true,
		true,
		true,
		"main",
	}

	applyStruct := Apply{
		"",
		true,
	}

	payload := cardPayload{
		createStruct,
		applyStruct,
	}

	req := e.Client.NewRequest()
	req.Url = "https://api.louisvuitton.com/api/eng-us/checkout/payment/creditcard"
	req.Method = "POST"
	req.Headers = []map[string]string{
		{"accept": "application/json, text/plain, */*"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
		//{"content-length":""},
		{"content-type": "application/json"},
		{"origin": "https://us.louisvuitton.com"},
		{"referer": "https://us.louisvuitton.com/eng-us/checkout/payment"},
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": `"Windows"`},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-site"},
		{"user-agent": e.UserAgent},
	}

	req.SetJSONBody(payload)

	resp, err := req.Do()
	if err != nil {
		// handle error, set next stage
		return
	}

	respBody := string(resp.Body)
	fmt.Println(respBody)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		ok := checkPage(resp.Body)
		if ok {
			e.Stage = Order
			return
		} else {
			e.WaitR()

			return
		}

	} else if resp.StatusCode == 403 {
		e.PaymentRetries++
		if e.PaymentRetries > 4 {
			e.PaymentRetries = 0
			e.RotateProxy()
		}
		e.getAkamaiSensor()
		e.submitSensor()
		e.Stage = SubmitPayment
		return

	} else if resp.StatusCode == 500 {
		errorResp := string(resp.Body)
		errorMsg := gjson.Get(errorResp, "errors.0.errorCode").String()
		e.UpdateStatus(errorMsg, activityapi.LogLevel)
		e.WaitM()
		return

	} else {

		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

func (e *LVTask) order() {
	e.UpdateStatus("Submitting order", activityapi.LogLevel)

	req := e.Client.NewRequest()
	req.Url = "https://api.louisvuitton.com/api/eng-us/checkout/order/commit"
	req.Method = "POST"
	req.Headers = []map[string]string{
		{"accept": "application/json, text/plain, */*"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
		//{"content-length":""},
		{"content-type": "application/json"},
		{"origin": "https://us.louisvuitton.com"},
		{"referer": "https://us.louisvuitton.com/eng-us/checkout/review"},
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": `"Windows"`},
		{"sec-fetch-dest": "empty"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-site": "same-site"},
		{"user-agent": e.UserAgent},
	}
	req.Body = strings.NewReader("{}")

	resp, err := req.Do()
	if err != nil {
		// handle error, set next stage
		return
	}

	respBody := string(resp.Body)
	fmt.Println(respBody)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		ok := checkPage(resp.Body)
		if ok {
			e.Stop()
			return
		} else {
			e.WaitR()

			return
		}

	} else if resp.StatusCode == 403 {
		e.OrderRetries++
		if e.OrderRetries > 4 {
			e.OrderRetries = 0
			e.RotateProxy()
		}
		e.getAkamaiSensor()
		e.submitSensor()
		e.Stage = Order
		return

	} else if resp.StatusCode == 500 {
		errorResp := string(resp.Body)
		errorMsg := gjson.Get(errorResp, "errors.0.errorCode").String()
		e.UpdateStatus(errorMsg, activityapi.LogLevel)
		return

	} else {
		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

func (e *LVTask) littyATC() {
	e.UpdateStatus("Inverting cart marticies", activityapi.LogLevel) // quasi professional cart matrix inverter moment

	payload := littyLoad{ // only SkuId and Quantity are necessary, other keys from default ATC were included as a precautionary measure
		e.Product.Identifier,
		[]string{e.Product.Identifier},
		[]string{"c21b0b857c24b1925dc02c7070e066b3"},
		e.ItemId,
		1,
	}

	req := e.Client.NewRequest()
	req.Url = "https://secure.louisvuitton.com/rest/bean/vuitton/commerce/services/cart/1_0/CartService/addToCart?storeLang=eng-us"
	req.Method = "POST"
	req.Headers = []map[string]string{
		{"sec-ch-ua": `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		{"accept": "application/json, text/plain, */*"},
		{"content-type": "application/json"},
		{"sec-ch-ua-mobile": "?0"},
		{"user-agent": e.UserAgent},
		{"sec-ch-ua-platform": `"Windows"`},
		{"origin": "https://us.louisvuitton.com"},
		{"sec-fetch-site": "same-site"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-dest": "empty"},
		{"referer": "https://us.louisvuitton.com/eng-us/products/lv-trainer-sneaker-nvprod3710063v/1AAHS3"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
		//{"content-length":""},
	}
	req.SetJSONBody(payload)

	resp, err := req.Do()
	if err != nil {
		// handle error, set next stage
		return
	}

	respBody := string(resp.Body)
	fmt.Println(respBody)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {

		e.Stage = GetCart
		return

	} else if resp.StatusCode == 403 {
		e.getAkamaiSensor()
		e.submitSensor()
		e.Stage = LittyATC
		return

	} else {
		e.UpdateStatus(fmt.Sprintf("Bad status code: %s", resp.Status), activityapi.LogLevel)
		e.WaitR()
		return
	}
}

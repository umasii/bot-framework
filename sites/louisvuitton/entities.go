package louisvuitton

const (
	GetSecure     = "GetSecure"
	Account       = "Account"
	Product       = "Product"
	Variants      = "Variants"
	Login         = "Login"
	GetCart       = "GetCart"
	PrepareCart   = "PrepareCart"
	GetBilling    = "GetBilling"
	GetCard       = "GetCard"
	GetSensor     = "GetSensor"
	GetU          = "GetU"
	Akamai        = "Akamai"
	AkamaiPixel   = "AkamaiPixel"
	SubmitPayment = "SubmitPayment"
	Order         = "Order"
	LittyATC      = "LittyATC"
	ProductUrl    = "ProductUrl"

	AK_API_KEY = "[redacted]"
)

type AccountData struct {
	email    string
	password string
}

type genSensorData struct {
	UserAgent string `json:"userAgent"`
	Site      string `json:"site"`
	Cookie    string `json:"abck"`
}

type genPixelData struct {
	Site      string `json:"site"`
	GVal      string `json:"gVal"`
	UniqueVal string `json:"bazadebezolkohpepadr"`
	UserAgent string `json:"userAgent"`
	Version   string `json:"version"`
}

type sensorPayload struct {
	SensorData string `json:"sensor_data"`
}

type loginPayload struct {
	Email    string `json:"login"`
	Password string `json:"password"`
}

type Create struct {
	CreditCardType         string `json:"creditCardType"`
	CreditCardName         string `json:"creditCardName"`
	CardVerificationNumber string `json:"cardVerificationNumber"`
	CreatedFromSavedCard   bool   `json:"createdFromSavedCard"`
	UseExistingCreditCard  bool   `json:"useExistingCreditCard"`
	UseExistingAddress     bool   `json:"useExistingAddress"`
	ValidateCreditCard     bool   `json:"validateCreditCard"`
	BillingAddressName     string `json:"billingAddressName"`
}

type Apply struct {
	ThirdPartyPaymentTypeName string `json:"thirdPartyPaymentTypeName"`
	ApplyDefaultPaymentGroup  bool   `json:"applyDefaultPaymentGroup"`
}

type cardPayload struct {
	Createstruct Create `json:"create"`
	Applystruct  Apply  `json:"apply"`
}

type littyLoad struct {
	SkuId            string   `json:"skuId"`
	CatalogRefId     []string `json:"catalogRefIds"`
	CatalogRefIdKeys []string `json:"catalogRefIdKeys"`
	ProductId        string   `json:"productId"`
	Quantity         int      `json:"quantity"`
}

package ActivityApi

const (
	LogLevel          = "log"
	WarningLevel      = "warning"
	ErrorLevel        = "error"
	SuccessLevel      = "success"
	NotificationLevel = "notification"
)

type UserData struct {
	UserID        string `json:"userID"`
	ShippingState string `json:"State"`
}

type TaskSettings struct {
	Site          string `json:"Site"`
	Product       string `json:"product"`
	//ProductIDType string `json:"PIDType"` // ie: sku, offerid, keywords, etc.. *this is assuming this will be a task input
	Mode          string `json:"Mode"`
}

type TaskResults struct {
	CheckedOut         bool   `json:"checkedOut"`
	CheckoutStatusCode int    `json:"StatusCode"`
	CheckoutMessage    string `json:"CheckoutMessage,omitempty"`
	FailureReason      string `json:"FailureReason,omitempty"`
	CheckoutTime		int 	`json:"CheckoutTime"`
}

type InstanceInfo struct {
	OS                   string `json:"OS"`
	TotalTasks           int    `json:"totalTasks"`
	TotalTasksForProduct int    `json:"totalTasksPerProd"`
	Time                 int64    `json:"time"`
}

type RecpData struct {
	UserInfo       UserData     `json:"userInfo"`
	Settings       TaskSettings `json:"settings"`
	Results        TaskResults  `json:"results"`
	Instance       InstanceInfo `json:"instance"`
	AdditionalData interface{}  `json:"AddInfo,omitempty"`
}


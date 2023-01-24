package Profiles

type Profile struct {
	ProfileID   int
	GroupID     int
	ProfileName string
	GroupName 	string
	Email       string
	Billing     Billing
	Shipping    Address
	Options     Options
}

type Billing struct {
	CardHolder     string
	CardNumber     string
	ExpMonth       string
	ExpYear        string
	CVC            string
	CardType       string
	BillingAddress Address
}

type Address struct {
	FirstName  string
	LastName   string
	Phone      string
	Address    string
	Address2   string
	City       string
	State      string
	Country    string
	PostalCode string
}

type Options struct {
	SameBillingAsShiping bool
	OnlyOneCheckout      bool
}

type ProfileGroup struct {
	GroupName string
	GroupID   int
	Profiles  []Profile
}

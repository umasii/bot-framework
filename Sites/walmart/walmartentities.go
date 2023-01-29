package walmart

import "time"

type Mode string

const (
	ModeOne Mode = "ModeOne"
	ModeTwo Mode = "ModeTwo"
)

const (
	Account       = "Account"
	ShippingRates = "Shipping rates"
	GetCart       = "Get Cart"
	PreEncrypt    = "PreEncrypt"
	CC            = "Cc"
	MonitorStage  = "Monitor"
	PreAdd        = "Pre Add"
	Cart          = "cart"
	Checkout      = "Checkout"
	Shipping      = "Shipping"
	Payment       = "Payment"
	Order         = "Order"
	StartMonitors = "Start monitors"
	WaitForStock  = "Wait for stock"
	DeleteItem    = "Delete item"
)

type PxInfo struct {
	Site  string  `json:"site"`
	Proxy string  `json:"proxy"`
	Key   string  `json:"key"`
	SetID *string `json:"setID,omitempty"`
	Uuid  *string `json:"uuid,omitempty"`
	Vid   *string `json:"vid,omitempty"`
	Token *string `json:"token,omitempty"`
}
type personName struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type signupCaptcha struct {
	SensorData string `json:"sensorData"`
}
type accountData struct {
	PersonName                personName    `json:"personName"`
	Email                     string        `json:"email"`
	Password                  string        `json:"password"`
	Rememberme                bool          `json:"rememberme"`
	ShowRememberme            string        `json:"showRememberme"`
	EmailNotificationAccepted bool          `json:"emailNotificationAccepted"`
	Captcha                   signupCaptcha `json:"captcha"`
}
type ccdata struct {
	AddressLineOne string `json:"addressLineOne"`
	AddressLineTwo string `json:"addressLineTwo"`
	CardType       string `json:"cardType"`
	EncryptedCvv   string `json:"encryptedCvv"`
	EncryptedPan   string `json:"encryptedPan"`
	ExpiryMonth    string `json:"expiryMonth"`
	ExpiryYear     string `json:"expiryYear"`
	FirstName      string `json:"firstName"`
	IntegrityCheck string `json:"integrityCheck"`
	IsGuest        bool   `json:"isGuest"`
	KeyId          string `json:"keyId"`
	LastName       string `json:"lastName"`
	Phase          string `json:"phase"`
	Phone          string `json:"phone"`
	PostalCode     string `json:"postalCode"`
	State          string `json:"state"`
}

type atcParams struct {
	OfferId  string `json:"offerId"`
	Quantity int    `json:"quantity"`
}

type shippingInfo struct {
	AddressLineOne     string   `json:"addressLineOne"`
	AddressLineTwo     string   `json:"addressLineTwo"`
	City               string   `json:"city"`
	FirstName          string   `json:"firstName"`
	LastName           string   `json:"lastName"`
	Phone              string   `json:"phone"`
	Email              string   `json:"email"`
	MarketingEmailPref bool     `json:"marketingEmailPref"`
	PostalCode         string   `json:"postalCode"`
	State              string   `json:"state"`
	CountryCode        string   `json:"countryCode"`
	AddressType        string   `json:"addressType"`
	ChangedFields      []string `json:"changedFields"`
}

type payments struct {
	PaymentType    string `json:"paymentType"`
	CardType       string `json:"cardType"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	AddressLineOne string `json:"addressLineOne"`
	AddressLineTwo string `json:"addressLineTwo"`
	City           string `json:"city"`
	State          string `json:"state"`
	PostalCode     string `json:"postalCode"`
	ExpiryMonth    string `json:"expiryMonth"`
	ExpiryYear     string `json:"expiryYear"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	EncryptedPan   string `json:"encryptedPan"`
	EncryptedCvv   string `json:"encryptedCvv"`
	IntegrityCheck string `json:"integrityCheck"`
	KeyId          string `json:"keyId"`
	Phase          string `json:"phase"`
	PiHash         string `json:"piHash"`
}

type paymentInfo struct {
	Payments     []payments `json:"payments"`
	CvvInSession bool       `json:"cvvInSession"`
}
type voltagePayments struct {
	PaymentType    string `json:"paymentType"`
	EncryptedCvv   string `json:"encryptedCvv"`
	EncryptedPan   string `json:"encryptedPan"`
	IntegrityCheck string `json:"integrityCheck"`
	KeyId          string `json:"keyId"`
	Phase          string `json:"phase"`
}

type orderInfo struct {
	CvvInSession    bool              `json:"cvvInSession"`
	VoltagePayments []voltagePayments `json:"voltagePayments"`
}
type pxResp struct {
	Px3   string `json:"_px3"`
	Ua    string `json:"useragent"`
	SetID string `json:"setID"`
	Uuid  string `json:"uuid"`
	Vid   string `json:"vid"`
}

type preAdd struct {
	SavedItems []SavedItems `json:"savedItems"`
}

type SavedItems struct {
	ItemId string `json:"id"`
	Pid    string `json:"USItemId"`
}

type ShippingData struct {
	PostalCode            string `json:"postalCode"`
	ResponseGroup         string `json:"responseGroup"`
	IncludePickUpLocation bool   `json:"includePickUpLocation"`
	PersistLocation       bool   `json:"persistLocation"`
	ClientName            string `json:"clientName"`
	StoreMeta             bool   `json:"storeMeta"`
	Plus                  bool   `json:"plus"`
}

type ShippingResp struct {
	Stores []StoreResp `json:"stores"`
}

type StoreResp struct {
	Types   []string `json:"types"`
	StoreId string   `json:"storeId"`
}

type SavedShippingRates struct {
	StoreList    []SavedStoreList `json:"storeList"`
	PostalCode   string           `json:"postalCode"`
	City         string           `json:"city"`
	State        string           `json:"state"`
	IsZipLocated bool             `json:"isZipLocated"`
	Crt          string           `json:"crt:CRT"`
	CustomerId   string           `json:"customerId:CID"`
	CustomerType string           `json:"customerType:type"`
	Affiliate    string           `json:"affiliateInfo:com.wm.reflector"`
}

type SavedStoreList struct {
	Id string `json:"id"`
}

type PieKeyResp struct {
	L       int
	E       int
	K       string
	KeyId   string
	PhaseId int
}

type PaymentResp struct {
	PiHash string `json:"piHash"`
}

type OrderResp struct {
	FailedReason string `json:"message"`
}

type TerraFirmaResp struct {
	Errors []interface{} `json:"errors"`
	Data   struct {
		ItemByOfferId struct {
			OfferList []struct {
				ProductAvailability struct {
					AvailabilityStatus string `json:"availabilityStatus"`
				} `json:"productAvailability"`
			} `json:"offerList"`
		} `json:"itemByOfferId"`
	} `json:"data"`
}

type SFLResp struct {
	Checkoutable bool `json:"checkoutable"`
	Cart         struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Preview    bool   `json:"preview"`
		Customerid string `json:"customerId"`
		Location   struct {
			Postalcode          string `json:"postalCode"`
			City                string `json:"city"`
			Stateorprovincecode string `json:"stateOrProvinceCode"`
			Countrycode         string `json:"countryCode"`
			Isdefault           bool   `json:"isDefault"`
		} `json:"location"`
		Storeids        []int         `json:"storeIds"`
		Itemcount       int           `json:"itemCount"`
		Currencycode    string        `json:"currencyCode"`
		Entityerrors    []interface{} `json:"entityErrors"`
		Itemcountbytype struct {
			Regular int `json:"regular"`
			Group   int `json:"group"`
		} `json:"itemCountByType"`
		Saveditemcountbytype struct {
			Regular int `json:"regular"`
			Group   int `json:"group"`
		} `json:"savedItemCountByType"`
		Hassubmaptypeitem bool `json:"hasSubmapTypeItem"`
		Fulfillmenttotals struct {
		} `json:"fulfillmentTotals"`
		Totals struct {
		} `json:"totals"`
		Saveditemcount    int           `json:"savedItemCount"`
		Tenantid          int           `json:"tenantId"`
		Verticalid        int           `json:"verticalId"`
		Localeid          string        `json:"localeId"`
		Fulfillmentgroups []interface{} `json:"fulfillmentGroups"`
		Hasalcoholicitem  bool          `json:"hasAlcoholicItem"`
	} `json:"cart"`
	Sfllist struct {
		Itemcount       int    `json:"itemCount"`
		Name            string `json:"name"`
		Status          string `json:"status"`
		Itemcountbytype struct {
			Regular int `json:"regular"`
			Group   int `json:"group"`
		} `json:"itemCountByType"`
	} `json:"sflList"`
	Items      []interface{} `json:"items"`
	Saveditems []struct {
		Quantity              int           `json:"quantity"`
		Price                 float64       `json:"price"`
		Unitvaluepricetype    string        `json:"unitValuePriceType"`
		Pickupunitprice       float64       `json:"pickupUnitPrice"`
		Usitemid              string        `json:"USItemId"`
		Ussellerid            string        `json:"USSellerId"`
		Legacyitemid          string        `json:"legacyItemId"`
		Upc                   string        `json:"upc"`
		Wupc                  string        `json:"wupc"`
		Offerid               string        `json:"offerId"`
		Gtin                  string        `json:"gtin"`
		Alternatewupcs        []interface{} `json:"alternateWupcs"`
		Productid             string        `json:"productId"`
		Name                  string        `json:"name"`
		Manufacturerproductid string        `json:"manufacturerProductId"`
		Productclasstype      string        `json:"productClassType"`
		Seller                struct {
			Name string `json:"name"`
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"seller"`
		Currentpricetype string `json:"currentPriceType"`
		Comparisonprice  struct {
			Isstrikethrough                bool `json:"isStrikethrough"`
			Iseligibleforassociatediscount bool `json:"isEligibleForAssociateDiscount"`
			Isrollback                     bool `json:"isRollback"`
			Isreducedprice                 bool `json:"isReducedPrice"`
			Isclearance                    bool `json:"isClearance"`
			Hidepriceforsoi                bool `json:"hidePriceForSOI"`
		} `json:"comparisonPrice"`
		Assets struct {
			Primary []struct {
				Num60  string `json:"60"`
				Num100 string `json:"100"`
			} `json:"primary"`
		} `json:"assets"`
		Productmarketattributes struct {
			GpcWithFc                 string `json:"gpc_with_fc"`
			KarfSalesUnit             string `json:"karf_sales_unit"`
			BrandCode                 string `json:"brand_code"`
			RhPath                    string `json:"rh_path"`
			AlternateShelves          string `json:"alternate_shelves"`
			KarfSubstitutionsAllowed  string `json:"karf_substitutions_allowed"`
			PrimaryCategoryPath       string `json:"primary_category_path"`
			VerticalEligibility       string `json:"vertical_eligibility"`
			ShelfDescription          string `json:"shelf_description"`
			KarfMaximumQuantityFactor string `json:"karf_maximum_quantity_factor"`
			IsPrivateLabelUnbranded   string `json:"is_private_label_unbranded"`
			ProductURLText            string `json:"product_url_text"`
			PrimaryShelfID            string `json:"primary_shelf_id"`
			CharPrimaryCategoryPath   string `json:"char_primary_category_path"`
			Segregation               string `json:"segregation"`
		} `json:"productMarketAttributes"`
		Marketingattributes struct {
			RhPath              string `json:"rh_path"`
			AlignChannelPrice   string `json:"align_channel_price"`
			WmDeptNum           string `json:"wm_dept_num"`
			ItemClassID         string `json:"item_class_id"`
			AvailabilityMsgFlag string `json:"availability_msg_flag"`
		} `json:"marketingAttributes"`
		Offertype              string   `json:"offerType"`
		Offerpublishstatus     string   `json:"offerPublishStatus"`
		Availablequantity      int      `json:"availableQuantity"`
		Maxitemcountperorder   float64  `json:"maxItemCountPerOrder"`
		Twodayshippingeligible bool     `json:"twoDayShippingEligible"`
		Pickupdiscounteligible bool     `json:"pickupDiscountEligible"`
		Shippingtier           string   `json:"shippingTier"`
		Shippingslatiers       []string `json:"shippingSlaTiers"`
		Shippingsladetail      struct {
			Slatier               string `json:"slaTier"`
			Geoitemclassification string `json:"geoItemClassification"`
		} `json:"shippingSlaDetail"`
		Tooverrideiroprice bool          `json:"toOverrideIROPrice"`
		Itemclassid        string        `json:"itemClassId"`
		Manufacturername   string        `json:"manufacturerName"`
		Brand              string        `json:"brand"`
		Productsegment     string        `json:"productSegment"`
		Producttype        string        `json:"productType"`
		Isconsumable       bool          `json:"isConsumable"`
		Type               string        `json:"type"`
		Entityerrors       []interface{} `json:"entityErrors"`
		Shippingeligible   bool          `json:"shippingEligible"`
		Edeliveryeligible  bool          `json:"eDeliveryEligible"`
		Groupcomponents    []interface{} `json:"groupComponents"`
	} `json:"savedItems"`
}

type CheckoutViewResp struct {
	DbSessionTokenMap struct {
		CXOPCST string `json:"CXO_PC_ST"`
	} `json:"dbSessionTokenMap"`
	Id               string `json:"id"`
	CheckoutFlowType string `json:"checkoutFlowType"`
	CartId           string `json:"cartId"`
	Items            []struct {
		Id               string  `json:"id"`
		OfferId          string  `json:"offerId"`
		ProductId        string  `json:"productId"`
		ProductName      string  `json:"productName"`
		ItemId           int     `json:"itemId"`
		SellerId         string  `json:"sellerId"`
		ThumbnailUrl     string  `json:"thumbnailUrl"`
		LegacySellerId   int     `json:"legacySellerId"`
		ProductClassType string  `json:"productClassType"`
		Quantity         int     `json:"quantity"`
		UnitPrice        float64 `json:"unitPrice"`
		Type             string  `json:"type"`
		Price            float64 `json:"price"`
		UnitOfMeasure    string  `json:"unitOfMeasure"`
		HasCarePlan      bool    `json:"hasCarePlan"`
		Brand            string  `json:"brand"`
		Discount         struct {
		} `json:"discount"`
		RhPath                      string `json:"rhPath"`
		IsWarrantyEligible          bool   `json:"isWarrantyEligible"`
		Category                    string `json:"category"`
		PrimaryCategory             string `json:"primaryCategory"`
		IsCarePlan                  bool   `json:"isCarePlan"`
		IsEgiftCard                 bool   `json:"isEgiftCard"`
		IsAssociateDiscountEligible bool   `json:"isAssociateDiscountEligible"`
		IsShippingPassEligible      bool   `json:"isShippingPassEligible"`
		ShippingTier                string `json:"shippingTier"`
		IsTwoDayShippingEligible    bool   `json:"isTwoDayShippingEligible"`
		MeetsSla                    bool   `json:"meetsSla"`
		ClassId                     string `json:"classId"`
		MaxQuantityPerOrder         int    `json:"maxQuantityPerOrder"`
		IsSubstitutable             bool   `json:"isSubstitutable"`
		IsInstaWatch                bool   `json:"isInstaWatch"`
		IsAlcoholic                 bool   `json:"isAlcoholic"`
		IsSnapEligible              bool   `json:"isSnapEligible"`
		IsAgeRestricted             bool   `json:"isAgeRestricted"`
		IsSubstitutionsAllowed      bool   `json:"isSubstitutionsAllowed"`
		FulfillmentSelection        struct {
			FulfillmentOption string `json:"fulfillmentOption"`
			ShipMethod        string `json:"shipMethod"`
			AvailableQuantity int    `json:"availableQuantity"`
		} `json:"fulfillmentSelection"`
		ServicePlanType string        `json:"servicePlanType"`
		Errors          []interface{} `json:"errors"`
		WfsEnabled      bool          `json:"wfsEnabled"`
		IsAlcohol       bool          `json:"isAlcohol"`
	} `json:"items"`
	Shipping struct {
		PostalCode string `json:"postalCode"`
		City       string `json:"city"`
		State      string `json:"state"`
	} `json:"shipping"`
	Summary struct {
		SubTotal               float64 `json:"subTotal"`
		ShippingIsEstimate     bool    `json:"shippingIsEstimate"`
		TaxIsEstimate          bool    `json:"taxIsEstimate"`
		GrandTotal             float64 `json:"grandTotal"`
		QuantityTotal          int     `json:"quantityTotal"`
		AmountOwed             float64 `json:"amountOwed"`
		MerchandisingFeesTotal int     `json:"merchandisingFeesTotal"`
		ShippingCosts          []struct {
			Label  string  `json:"label"`
			Type   string  `json:"type"`
			Cost   float64 `json:"cost"`
			Method string  `json:"method"`
		} `json:"shippingCosts"`
		ShippingTotal      float64 `json:"shippingTotal"`
		HasSurcharge       bool    `json:"hasSurcharge"`
		PreTaxTotal        float64 `json:"preTaxTotal"`
		AddOnServicesTotal int     `json:"addOnServicesTotal"`
		ItemsSubTotal      float64 `json:"itemsSubTotal"`
	} `json:"summary"`
	PickupPeople []interface{} `json:"pickupPeople"`
	Email        string        `json:"email"`
	Buyer        struct {
		CustomerAccountId      string `json:"customerAccountId"`
		FirstName              string `json:"firstName"`
		LastName               string `json:"lastName"`
		Email                  string `json:"email"`
		IsGuest                bool   `json:"isGuest"`
		IsAssociate            bool   `json:"isAssociate"`
		ApplyAssociateDiscount bool   `json:"applyAssociateDiscount"`
		HasCapOneCard          bool   `json:"hasCapOneCard"`
	} `json:"buyer"`
	AllowedPaymentTypes []struct {
		Type        string `json:"type"`
		CvvRequired bool   `json:"cvvRequired"`
	} `json:"allowedPaymentTypes"`
	Registries                []interface{} `json:"registries"`
	Payments                  []interface{} `json:"payments"`
	CardsToDisable            []interface{} `json:"cardsToDisable"`
	AllowedPaymentPreferences []interface{} `json:"allowedPaymentPreferences"`
	IsRCFEligible             bool          `json:"isRCFEligible"`
	IsMarketPlaceItemsExist   bool          `json:"isMarketPlaceItemsExist"`
	Version                   string        `json:"version"`
	SharedCategory            struct {
		ShippingGroups []struct {
			ItemIds              []string `json:"itemIds"`
			Seller               string   `json:"seller"`
			ShippingTier         string   `json:"shippingTier"`
			MeetSla              bool     `json:"meetSla"`
			DefaultSelection     bool     `json:"defaultSelection"`
			FulfillmentOption    string   `json:"fulfillmentOption"`
			ShippingGroupOptions []struct {
				Method                string  `json:"method"`
				MethodDisplay         string  `json:"methodDisplay"`
				Selected              bool    `json:"selected"`
				Charge                float64 `json:"charge"`
				DeliveryDate          int64   `json:"deliveryDate"`
				AvailableDate         int64   `json:"availableDate"`
				FulfillmentOption     string  `json:"fulfillmentOption"`
				OnlineStoreId         int     `json:"onlineStoreId"`
				IsThresholdShipMethod bool    `json:"isThresholdShipMethod"`
			} `json:"shippingGroupOptions"`
			IsEdelivery      bool          `json:"isEdelivery"`
			HasWFSItem       bool          `json:"hasWFSItem"`
			ItemSellerGroups []interface{} `json:"itemSellerGroups"`
		} `json:"shippingGroups"`
		PickupGroupsByStore []struct {
			PickupGroups []struct {
				ItemIds            []string `json:"itemIds"`
				Seller             string   `json:"seller"`
				DefaultSelection   bool     `json:"defaultSelection"`
				FulfillmentOption  string   `json:"fulfillmentOption"`
				PickupGroupOptions []struct {
					Method             string `json:"method"`
					MethodDisplay      string `json:"methodDisplay"`
					Selected           bool   `json:"selected"`
					Charge             int    `json:"charge"`
					DeliveryDate       int64  `json:"deliveryDate"`
					AvailableDate      int64  `json:"availableDate"`
					FulfillmentOption  string `json:"fulfillmentOption"`
					Type               string `json:"type"`
					AvailabilityStatus string `json:"availabilityStatus"`
				} `json:"pickupGroupOptions"`
			} `json:"pickupGroups"`
			StoreId int `json:"storeId"`
			Address struct {
				PostalCode     string `json:"postalCode"`
				AddressLineOne string `json:"addressLineOne"`
				City           string `json:"city"`
				State          string `json:"state"`
				Address1       string `json:"address1"`
				Country        string `json:"country"`
			} `json:"address"`
			StoreType     string `json:"storeType"`
			StoreTypeId   int    `json:"storeTypeId"`
			StoreName     string `json:"storeName"`
			Selected      bool   `json:"selected"`
			StoreServices []struct {
				ServiceName string `json:"serviceName"`
				ServiceId   int    `json:"serviceId"`
				Active      bool   `json:"active"`
			} `json:"storeServices"`
			PickupTogetherEligibleStore bool `json:"pickupTogetherEligibleStore"`
		} `json:"pickupGroupsByStore"`
		IsShippingEnabled bool `json:"isShippingEnabled"`
	} `json:"sharedCategory"`
	BalanceToReachThreshold float64       `json:"balanceToReachThreshold"`
	EntityErrors            []interface{} `json:"entityErrors"`
	OneDaySelected          bool          `json:"oneDaySelected"`
	PaymentWithBagFee       bool          `json:"paymentWithBagFee"`
	GiftDetails             struct {
		GiftOrder           bool `json:"giftOrder"`
		HasGiftEligibleItem bool `json:"hasGiftEligibleItem"`
		XoGiftingOptIn      bool `json:"xoGiftingOptIn"`
	} `json:"giftDetails"`
	CanApplyDetails []interface{} `json:"canApplyDetails"`
	DbName          string        `json:"dbName"`
	Jwt             string        `json:"jwt"`
}

type AtcResp struct {
	DbSessionMap struct {
		CXOCARTST string `json:"CXO_CART_ST"`
	} `json:"dbSessionMap"`
	Checkoutable bool `json:"checkoutable"`
	Cart         struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Preview    bool   `json:"preview"`
		CustomerId string `json:"customerId"`
		Location   struct {
			PostalCode          string `json:"postalCode"`
			City                string `json:"city"`
			StateOrProvinceCode string `json:"stateOrProvinceCode"`
			CountryCode         string `json:"countryCode"`
			IsDefault           bool   `json:"isDefault"`
			IsZipLocated        bool   `json:"isZipLocated"`
		} `json:"location"`
		StoreIds        []int         `json:"storeIds"`
		ItemCount       int           `json:"itemCount"`
		CurrencyCode    string        `json:"currencyCode"`
		EntityErrors    []interface{} `json:"entityErrors"`
		ItemCountByType struct {
			Regular int `json:"regular"`
			Group   int `json:"group"`
		} `json:"itemCountByType"`
		SavedItemCountByType struct {
			Regular int `json:"regular"`
			Group   int `json:"group"`
		} `json:"savedItemCountByType"`
		HasSubmapTypeItem bool `json:"hasSubmapTypeItem"`
		FulfillmentTotals struct {
			S2H struct {
				Name       string `json:"name"`
				SellerId   string `json:"sellerId"`
				Type       string `json:"type"`
				SellerName string `json:"sellerName"`
				Methods    []struct {
					Method                  string        `json:"method"`
					Price                   float64       `json:"price"`
					ItemIds                 []interface{} `json:"itemIds"`
					ThresholdOrderTotal     float64       `json:"thresholdOrderTotal"`
					BalanceToReachThreshold float64       `json:"balanceToReachThreshold"`
				} `json:"methods"`
			} `json:"S2H"`
		} `json:"fulfillmentTotals"`
		Totals struct {
			SubTotal             float64 `json:"subTotal"`
			SellerShippingTotals []struct {
				SellerId   string `json:"sellerId"`
				SellerName string `json:"sellerName"`
				Methods    []struct {
					Price  float64 `json:"price"`
					Name   string  `json:"name"`
					Method string  `json:"method"`
				} `json:"methods"`
			} `json:"sellerShippingTotals"`
			ShippingTotal         float64 `json:"shippingTotal"`
			HasSurcharge          bool    `json:"hasSurcharge"`
			GrandTotal            float64 `json:"grandTotal"`
			AddOnServicesTotal    int     `json:"addOnServicesTotal"`
			ItemsSubTotal         float64 `json:"itemsSubTotal"`
			OriginalShippingTotal float64 `json:"originalShippingTotal"`
			MemberShipDiscount    int     `json:"memberShipDiscount"`
		} `json:"totals"`
		SavedItemCount        int           `json:"savedItemCount"`
		ShipMethodDefaultRule string        `json:"shipMethodDefaultRule"`
		TenantId              int           `json:"tenantId"`
		VerticalId            int           `json:"verticalId"`
		LocaleId              string        `json:"localeId"`
		FulfillmentGroups     []interface{} `json:"fulfillmentGroups"`
		HasAlcoholicItem      bool          `json:"hasAlcoholicItem"`
		GiftOrder             bool          `json:"giftOrder"`
		HasGiftEligibleItem   bool          `json:"hasGiftEligibleItem"`
	} `json:"cart"`
	Items []struct {
		Id                 string  `json:"id"`
		Quantity           int     `json:"quantity"`
		Price              float64 `json:"price"`
		LinePrice          float64 `json:"linePrice"`
		UnitValuePrice     float64 `json:"unitValuePrice,omitempty"`
		UnitValuePriceType string  `json:"unitValuePriceType"`
		StoreFrontType     string  `json:"storeFrontType"`
		StoreFrontId       struct {
			StoreUUID        string `json:"storeUUID"`
			StoreId          string `json:"storeId"`
			Preferred        bool   `json:"preferred"`
			SemStore         bool   `json:"semStore"`
			OnlineStoreFront bool   `json:"onlineStoreFront"`
			USStoreId        int    `json:"USStoreId"`
		} `json:"storeFrontId"`
		USItemId         string        `json:"USItemId"`
		USSellerId       string        `json:"USSellerId"`
		LegacyItemId     string        `json:"legacyItemId"`
		Upc              string        `json:"upc,omitempty"`
		Wupc             string        `json:"wupc"`
		OfferId          string        `json:"offerId"`
		Gtin             string        `json:"gtin"`
		AlternateWupcs   []interface{} `json:"alternateWupcs"`
		ProductId        string        `json:"productId"`
		Name             string        `json:"name"`
		ProductClassType string        `json:"productClassType"`
		Seller           struct {
			Name string `json:"name"`
			Type string `json:"type"`
			Id   string `json:"id"`
		} `json:"seller"`
		CurrentPriceType string `json:"currentPriceType"`
		ComparisonPrice  struct {
			IsStrikethrough                bool `json:"isStrikethrough"`
			IsEligibleForAssociateDiscount bool `json:"isEligibleForAssociateDiscount"`
			IsRollback                     bool `json:"isRollback"`
			IsReducedPrice                 bool `json:"isReducedPrice"`
			IsClearance                    bool `json:"isClearance"`
			HidePriceForSOI                bool `json:"hidePriceForSOI"`
		} `json:"comparisonPrice"`
		Assets struct {
			Primary []struct {
				Field1 string `json:"100"`
				Field2 string `json:"60"`
			} `json:"primary"`
		} `json:"assets"`
		ProductMarketAttributes struct {
			KarfSalesUnit             string    `json:"karf_sales_unit,omitempty"`
			BrandCode                 string    `json:"brand_code,omitempty"`
			New                       string    `json:"new,omitempty"`
			KarfMaximumOrderQuantity  string    `json:"karf_maximum_order_quantity,omitempty"`
			RhPath                    string    `json:"rh_path"`
			AlternateShelves          string    `json:"alternate_shelves"`
			KarfSubstitutionsAllowed  string    `json:"karf_substitutions_allowed,omitempty"`
			PrimaryCategoryPath       string    `json:"primary_category_path"`
			VerticalEligibility       string    `json:"vertical_eligibility,omitempty"`
			ShelfDescription          string    `json:"shelf_description"`
			KarfMaximumQuantityFactor string    `json:"karf_maximum_quantity_factor,omitempty"`
			IsPrivateLabelUnbranded   string    `json:"is_private_label_unbranded"`
			ProductUrlText            string    `json:"product_url_text"`
			Newuntildate              time.Time `json:"newuntildate,omitempty"`
			PrimaryShelfId            string    `json:"primary_shelf_id"`
			CharPrimaryCategoryPath   string    `json:"char_primary_category_path"`
			Segregation               string    `json:"segregation"`
			DisplayStatus             string    `json:"display_status,omitempty"`
		} `json:"productMarketAttributes"`
		MarketingAttributes struct {
			LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
			RhPath        string `json:"rh_path"`
			WmtDeptNum    string `json:"wmt_dept_num,omitempty"`
			WmDeptNum     string `json:"wm_dept_num"`
			ItemClassId   string `json:"item_class_id"`
		} `json:"marketingAttributes"`
		OfferType               string        `json:"offerType"`
		OfferPublishStatus      string        `json:"offerPublishStatus"`
		AvailableQuantity       int           `json:"availableQuantity"`
		MaxItemCountPerOrder    int           `json:"maxItemCountPerOrder"`
		TwoDayShippingEligible  bool          `json:"twoDayShippingEligible"`
		PickupDiscountEligible  bool          `json:"pickupDiscountEligible"`
		ShippingTier            string        `json:"shippingTier"`
		WfsEnabled              bool          `json:"wfsEnabled"`
		ShippingSlaTiers        []string      `json:"shippingSlaTiers"`
		ToOverrideIROPrice      bool          `json:"toOverrideIROPrice"`
		ItemClassId             string        `json:"itemClassId"`
		ManufacturerName        string        `json:"manufacturerName"`
		Brand                   string        `json:"brand,omitempty"`
		ProductSegment          string        `json:"productSegment"`
		ProductType             string        `json:"productType"`
		IsConsumable            bool          `json:"isConsumable"`
		PrimaryCategoryPath     string        `json:"primaryCategoryPath"`
		CharPrimaryCategoryPath string        `json:"charPrimaryCategoryPath"`
		Type                    string        `json:"type"`
		UnitOfMeasure           string        `json:"unitOfMeasure"`
		Sort                    int64         `json:"sort"`
		EntityErrors            []interface{} `json:"entityErrors"`
		ShippingOptions         struct {
			S2H []struct {
				Status   string  `json:"status"`
				Quantity float64 `json:"quantity"`
				Methods  []struct {
					Price        float64 `json:"price"`
					Name         string  `json:"name"`
					Method       string  `json:"method"`
					ShippingTier string  `json:"shippingTier"`
				} `json:"methods"`
				StoreAddress struct {
					City       string `json:"city"`
					Country    string `json:"country"`
					PostalCode string `json:"postalCode"`
					State      string `json:"state"`
				} `json:"storeAddress"`
			} `json:"S2H"`
			PUT []struct {
				Status   string  `json:"status"`
				Quantity float64 `json:"quantity"`
				Methods  []struct {
					Price        float64 `json:"price"`
					Name         string  `json:"name"`
					Method       string  `json:"method"`
					ShippingTier string  `json:"shippingTier"`
				} `json:"methods"`
				StoreAddress struct {
					Address1   string `json:"address1"`
					City       string `json:"city"`
					Country    string `json:"country"`
					PostalCode string `json:"postalCode"`
					State      string `json:"state"`
				} `json:"storeAddress"`
				UsstoreId int `json:"usstoreId"`
			} `json:"PUT,omitempty"`
		} `json:"shippingOptions"`
		IsWarrantyEligible     bool   `json:"isWarrantyEligible"`
		IsServicePlansEligible bool   `json:"isServicePlansEligible"`
		IsSnapEligible         bool   `json:"isSnapEligible"`
		WeightIncrement        int    `json:"weightIncrement"`
		ManufacturerProductId  string `json:"manufacturerProductId,omitempty"`
		GiftDetails            struct {
			Eligibility string `json:"eligibility"`
			GiftOpts    struct {
				GiftOverboxEligible bool `json:"giftOverboxEligible"`
			} `json:"giftOpts"`
		} `json:"giftDetails,omitempty"`
	} `json:"items"`
	NextDayEligible      bool  `json:"nextDayEligible"`
	EDeliveryCart        bool  `json:"eDeliveryCart"`
	CanAddWARPItemToCart bool  `json:"canAddWARPItemToCart"`
	CartVersion          int   `json:"cartVersion"`
	LastModifiedTime     int64 `json:"lastModifiedTime"`
}

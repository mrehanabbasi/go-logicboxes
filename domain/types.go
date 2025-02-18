package domain

import "github.com/mrehanabbasi/go-logicboxes/core"

type SortBy string

type (
	PrivacyState       string
	RegistrationStatus string
	SortOrder          map[SortBy]bool
)

type SuggestNames map[string]struct {
	Status string         `json:"status"`
	InGa   core.JSONBool  `json:"in_ga"`
	Score  core.JSONFloat `json:"score"`
	Spin   string         `json:"spin"`
}

type RegisterResponse struct {
	ActionTypeDesc          string            `json:"actiontypedesc"`
	UnutilisedSellingAmount core.JSONFloat    `json:"unutilisedsellingamount"`
	SellingAmount           core.JSONFloat    `json:"sellingamount"`
	EntityID                string            `json:"entityid"`
	ActionStatus            string            `json:"actionstatus"`
	Status                  string            `json:"status"`
	EaqID                   string            `json:"eaqid"`
	PremiumDNSDetails       *RegisterResponse `json:"premiumdnsdetails,omitempty"`
	CustomerID              string            `json:"customerid"`
	Description             string            `json:"description"`
	ActionType              string            `json:"actiontype"`
	InvoiceID               string            `json:"invoiceid"`
	SellingCurrencySymbol   string            `json:"sellingcurrencysymbol"`
	ActionStatusDesc        string            `json:"actionstatusdesc"`
}

type Contact struct {
	Company       string   `json:"company"`
	Address1      string   `json:"address1"`
	Telno         string   `json:"telno"`
	TelnoCC       string   `json:"telnocc"`
	ContactID     string   `json:"contactid"`
	Type          string   `json:"type"`
	ContactType   []string `json:"contacttype"`
	CustomerID    string   `json:"customerid"`
	Country       string   `json:"country"`
	ParentKey     string   `json:"parentkey"`
	ContactStatus string   `json:"contactstatus"`
	State         string   `json:"state"`
	EmailAddr     string   `json:"emailaddr"`
	City          string   `json:"city"`
	Name          string   `json:"name"`
	ZIP           string   `json:"zip"`
}

type OrderDetail struct {
	Classkey          string          `json:"classkey"`
	AllowDeletion     core.JSONBool   `json:"allowdeletion"`
	OrderID           string          `json:"orderid"`
	NoOfNameServers   core.JSONUint16 `json:"noOfNameServers"`
	ParentKey         string          `json:"parentkey"`
	ProductCategory   string          `json:"productcategory"`
	CurrentStatus     string          `json:"currentstatus"`
	DomainName        string          `json:"domainname"`
	Description       string          `json:"description"`
	MultilingualFlag  string          `json:"multilingualflag"`
	TNCRequired       core.JSONBool   `json:"tnc_required"`
	PremiumDNSEnabled core.JSONBool   `json:"premiumdnsenabled"`
	GDPR              struct {
		Enabled  core.JSONBool `json:"enabled"`
		Eligible core.JSONBool `json:"eligible"`
	} `json:"gdpr"`
	CustomerID                 string          `json:"customerid"`
	Addons                     []string        `json:"addons"`
	BulkWhoIsOptOut            string          `json:"bulkwhoisoptout"`
	TechContactID              string          `json:"techcontactid"`
	IsImmediateReseller        core.JSONBool   `json:"isImmediateReseller"`
	CreationTime               core.JSONTime   `json:"creationtime"`
	DNSSec                     []string        `json:"dnssec"`
	JumpConditions             []string        `json:"jumpConditions"`
	RaaVerificationStartTime   core.JSONTime   `json:"raaVerificationStartTime"`
	CNS                        struct{}        `json:"cns"`
	Paused                     core.JSONBool   `json:"paused"`
	Admincontact               Contact         `json:"admincontact"`
	BillingContactID           string          `json:"billingcontactid"`
	PrivacyProtectedAllowed    core.JSONBool   `json:"privacyprotectedallowed"`
	DomSecret                  string          `json:"domsecret"`
	PremiumDNSAllowed          core.JSONBool   `json:"premiumdnsallowed"`
	ServiceProviderID          string          `json:"serviceproviderid"`
	Classname                  string          `json:"classname"`
	ResellerCost               core.JSONUint16 `json:"resellercost"`
	OrderStatus                []string        `json:"orderstatus"`
	EaqID                      string          `json:"eaqid"`
	EndTime                    core.JSONTime   `json:"endtime"`
	BillingContact             Contact         `json:"billingcontact"`
	AutoRenewTermType          string          `json:"autoRenewTermType"`
	RaaVerificationStatus      string          `json:"raaVerificationStatus"`
	EntityID                   string          `json:"entityid"`
	Recurring                  core.JSONBool   `json:"recurring"`
	ProductKey                 string          `json:"productkey"`
	NS1                        string          `json:"ns1"`
	NS2                        string          `json:"ns2"`
	NS3                        string          `json:"ns3"`
	NS4                        string          `json:"ns4"`
	NS5                        string          `json:"ns5"`
	NS6                        string          `json:"ns6"`
	ActionCompleted            core.JSONUint16 `json:"actioncompleted"`
	RegistrantContact          Contact         `json:"registrantcontact"`
	EntityTypeID               string          `json:"entitytypeid"`
	AutoRenewAttemptDuration   core.JSONUint16 `json:"autoRenewAttemptDuration"`
	CustomerCost               core.JSONFloat  `json:"customercost"`
	DomainStatus               []string        `json:"domainstatus"`
	OrderSuspendedByParent     core.JSONBool   `json:"orderSuspendedByParent"`
	MoneyBackPeriod            core.JSONUint16 `json:"moneybackperiod"`
	TechContact                Contact         `json:"techcontact"`
	RegistrantContactID        string          `json:"registrantcontactid"`
	AdminContactID             string          `json:"admincontactid"`
	IsOrderSuspendedUponExpiry core.JSONBool   `json:"isOrderSuspendedUponExpiry"`
	IsPrivacyProtected         core.JSONBool   `json:"isprivacyprotected"`
}

type NameServersResponse struct {
	ActionTypeDesc   string `json:"actiontypedesc"`
	EntityID         string `json:"entityid"`
	ActionStatus     string `json:"actionstatus"`
	Status           string `json:"status"`
	EaqID            string `json:"eaqid"`
	CurrentAction    string `json:"currentaction"`
	Description      string `json:"description"`
	ActionType       string `json:"actiontype"`
	ActionStatusDesc string `json:"actionstatusdesc"`
}

type ModifyPrivacyProtectionStatusResponse struct {
	ActionTypeDesc          string         `json:"actiontypedesc"`
	UnutilisedSellingAmount core.JSONFloat `json:"unutilisedsellingamount"`
	SellingAmount           core.JSONFloat `json:"sellingamount"`
	EntityID                string         `json:"entityid"`
	ActionStatus            string         `json:"actionstatus"`
	Status                  string         `json:"status"`
	EaqID                   string         `json:"eaqid"`
	CustomerID              string         `json:"customerid"`
	Description             string         `json:"description"`
	ActionType              string         `json:"actiontype"`
	InvoiceID               string         `json:"invoiceid"`
	SellingCurrencySymbol   string         `json:"sellingcurrencysymbol"`
	ActionStatusDesc        string         `json:"actionstatusdesc"`
	Message                 string         `json:"message"`
}

type ModifyAuthCodeResponse struct {
	ActionTypeDesc   string `json:"actiontypedesc"`
	EntityID         string `json:"entityid"`
	ActionStatus     string `json:"actionstatus"`
	Status           string `json:"status"`
	EaqID            string `json:"eaqid"`
	CurrentAction    string `json:"currentaction"`
	Description      string `json:"description"`
	ActionType       string `json:"actiontype"`
	ActionStatusDesc string `json:"actionstatusdesc"`
}

type TheftProtectionLockResponse struct {
	ActionTypeDesc   string `json:"actiontypedesc"`
	EntityID         string `json:"entityid"`
	ActionStatus     string `json:"actionstatus"`
	Status           string `json:"status"`
	EaqID            string `json:"eaqid"`
	Error            string `json:"error"`
	Description      string `json:"description"`
	ActionType       string `json:"actiontype"`
	ActionStatusDesc string `json:"actionstatusdesc"`
}

type GetTheListOfLocksAppliedOnDomainNameResponse struct {
	TransferLock bool `json:"transferlock"`
	CustomerLock bool `json:"customerlock"`
}

type CancelTransferResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DeleteResponse struct {
	Status        string `json:"status"`
	EaqID         string `json:"eaqid"`
	CurrentAction string `json:"currentaction"`
}

// Const for sort order.
const (
	SortByOrderID          SortBy = "orderid"
	SortByCustomerID       SortBy = "customerid"
	SortByEndTime          SortBy = "endtime"
	SortByTimestamp        SortBy = "timestamp"
	SortByEntityTypeID     SortBy = "entitytypeid"
	SortByCreationTime     SortBy = "creationtime"
	SortByCreationDateTime SortBy = "creationdt"

	PrivacyEnabled     PrivacyState = "true"
	PrivacyDisabled    PrivacyState = "false"
	PrivacyUnsupported PrivacyState = "na"

	DomRegUnknown       RegistrationStatus = "unknown"
	DomRegUnregistered  RegistrationStatus = "available"
	DomRegThroughUs     RegistrationStatus = "regthroughus"
	DomRegThroughOthers RegistrationStatus = "regthroughothers"
)

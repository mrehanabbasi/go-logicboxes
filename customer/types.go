package customer

import (
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/xpartacvs/go-resellerclub/core"
)

type LoginToken struct {
	token   string
	baseUrl string
}

type SignUpForm struct {
	Username              string `validate:"required,email" query:"username"`
	Password              string `validate:"required,min=9,max=16,rcpassword" query:"passwd"`
	Name                  string `validate:"required" query:"name"`
	Company               string `validate:"required" query:"company"`
	Address               string `validate:"required" query:"address-line-1"`
	AddressLine2          string `validate:"omitempty" query:"address-line-2,omitempty"`
	AddressLine3          string `validate:"omitempty" query:"address-line-3,omitempty"`
	City                  string `validate:"required" query:"city"`
	State                 string `validate:"required" query:"state"`
	OtherState            string `validate:"omitempty" query:"other-state,omitempty"`
	Country               string `validate:"required,iso3166_1_alpha2" query:"country"`
	Zipcode               string `validate:"required" query:"zipcode"`
	LanguageCode          string `validate:"required" query:"lang-pref"`
	PhoneCountryCode      string `validate:"required,len=2" query:"phone-cc"`
	Phone                 string `validate:"required,number" query:"phone"`
	AltPhoneCountryCode   string `validate:"omitempty,len=2" query:"alt-phone-cc,omitempty"`
	AltPhone              string `validate:"omitempty,number" query:"alt-phone,omitempty"`
	FaxCountryCode        string `validate:"omitempty,len=2" query:"fax-cc,omitempty"`
	Fax                   string `validate:"omitempty,number" query:"fax,omitempty"`
	MobileCountryCode     string `validate:"omitempty,len=2" query:"mobile-cc,omitempty"`
	Mobile                string `validate:"omitempty,number" query:"mobile,omitempty"`
	VatID                 string `validate:"omitempty" query:"vat-id,omitempty"`
	SmsConcent            bool   `validate:"omitempty" query:"sms-consent,omitempty"`
	EmailMarketingConcent bool   `validate:"omitempty" query:"email-marketing-consent,omitempty"`
	AcceptPolicy          bool   `validate:"omitempty" query:"accept-policy,omitempty"`
	CustomerId            string `validate:"-"`
}

type CustomerDetail struct {
	Id                      string          `json:"customerid,omitempty" validate:"-" query:"-"`
	Username                string          `json:"username,omitempty" validate:"omitempty,email" query:"username"`
	ResellerId              string          `json:"resellerid,omitempty" validate:"-" query:"-"`
	ParentId                string          `json:"parentid,omitempty" validate:"-" query:"-"`
	Name                    string          `json:"name,omitempty" validate:"omitempty" query:"name"`
	Company                 string          `json:"company,omitempty" validate:"omitempty" query:"company"`
	Email                   string          `json:"useremail,omitempty" validate:"-" query:"-"`
	PhoneCountryCode        string          `json:"telnocc,omitempty" validate:"omitempty,len=2,number" query:"phone-cc"`
	Phone                   string          `json:"telno,omitempty" validate:"omitempty,number" query:"phone"`
	AltPhoneCountryCode     string          `json:"-" validate:"omitempty,len=2,number" query:"alt-phone-cc,omitempty"`
	AltPhone                string          `json:"-" validate:"omitempty,number" query:"alt-phone,omitempty"`
	MobileCountryCode       string          `json:"mobilenocc,omitempty" validate:"omitempty,len=2,number" query:"mobile-cc,omitempty"`
	Mobile                  string          `json:"mobileno,omitempty" validate:"omitempty,number" query:"mobile,omitempty"`
	FaxCountryCode          string          `json:"-" validate:"omitempty,len=2" query:"faxnocc,omitempty"`
	Fax                     string          `json:"-" validate:"omitempty,number" query:"faxno,omitempty"`
	Address                 string          `json:"address1,omitempty" validate:"omitempty" query:"address-line-1"`
	AddressLine2            string          `json:"address2,omitempty" validate:"omitempty" query:"address-line-2,omitempty"`
	AddressLine3            string          `json:"address3,omitempty" validate:"omitempty" query:"address-line-3,omitempty"`
	City                    string          `json:"city,omitempty" validate:"omitempty" query:"city"`
	StateId                 string          `json:"stateid,omitempty" validate:"-" query:"-"`
	State                   string          `json:"state,omitempty" validate:"omitempty" query:"state"`
	OtherState              string          `json:"-" validate:"omitempty" query:"other-state,omitempty"`
	CountryCode             string          `json:"country,omitempty" validate:"omitempty,iso3166_1_alpha2" query:"country"`
	Zipcode                 string          `json:"zip,omitempty" validate:"omitempty" query:"zipcode"`
	LanguagePreference      string          `json:"langpref,omitempty" validate:"omitempty" query:"lang-pref"`
	VatEurope               string          `json:"-" validate:"omitempty" query:"vat-id,omitempty"`
	VatRussia               string          `json:"-" validate:"omitempty" query:"russia-vat-id,omitempty"`
	GstIndia                string          `json:"-" validate:"omitempty" query:"indian-gst-id,omitempty"`
	GstAustralia            string          `json:"-" validate:"omitempty" query:"australia-gst-id,omitempty"`
	GstNewZealand           string          `json:"-" validate:"omitempty" query:"newzealand-gst-id,omitempty"`
	GstSingapore            string          `json:"-" validate:"omitempty" query:"singapore-gst-id,omitempty"`
	Pin                     string          `json:"pin,omitempty" validate:"-" query:"-"`
	TimeCreation            core.JSONTime   `json:"creationdt,omitempty" validate:"-" query:"-"`
	Status                  string          `json:"customerstatus,omitempty" validate:"-" query:"-"`
	SalesContactId          string          `json:"salescontactid,omitempty" validate:"-" query:"-"`
	WebsiteCount            core.JSONUint16 `json:"websitecount,omitempty" validate:"-" query:"-"`
	TotalReceipts           core.JSONFloat  `json:"totalreceipts,omitempty" validate:"-" query:"-"`
	Is2FA                   core.JSONBool   `json:"twofactorauth_enabled,omitempty" validate:"-" query:"-"`
	Is2FASms                core.JSONBool   `json:"twofactorsmsauth_enabled,omitempty" validate:"-" query:"-"`
	Is2FAGoogle             core.JSONBool   `json:"twofactorgoogleauth_enabled,omitempty" validate:"-" query:"-"`
	IsDominicanTaxConfgired core.JSONBool   `json:"isDominicanTaxConfiguredByParent,omitempty" validate:"-" query:"-"`
}

type CustomerCriteria struct {
	core.Criteria
	Username       string            `validate:"omitempty" query:"username,omitempty"`
	Status         core.EntityStatus `validate:"omitempty" query:"status,omitempty"`
	Name           string            `validate:"omitempty" query:"name,omitempty"`
	Company        string            `validate:"omitempty" query:"company,omitempty"`
	City           string            `validate:"omitempty" query:"city,omitempty"`
	State          string            `validate:"omitempty" query:"state,omitempty"`
	ReceiptLowest  float64           `validate:"omitempty" query:"total-receipt-start,omitempty"`
	ReceiptHighest float64           `validate:"omitempty" query:"total-receipt-end,omitempty"`
}

type CustomerSearchResult struct {
	RequestedLimit  uint16
	RequestedOffset uint16
	TotalMatched    int
	Customers       []CustomerDetail
}

type ErrorAuthentication struct {
	core.JSONStatusResponse
	AuthLimit     core.JSONUint16 `json:"maxAttempts"`
	AuthRemaining core.JSONUint16 `json:"remainingLoginAttempts"`
}

func (t LoginToken) String() string {
	return t.token
}

func (t LoginToken) UrlFullPath() string {
	data := url.Values{}
	data.Add("role", "customer")
	data.Add("userLoginId", t.String())
	return "servlet/AutoLoginServlet?" + data.Encode()
}

func (t LoginToken) LoginUrl() string {
	return strings.TrimRight(t.baseUrl, "/") + t.UrlFullPath()
}

func (c *CustomerDetail) mergePrevious(prev CustomerDetail) error {
	if err := validator.New().Struct(c); err != nil {
		return err
	}

	valueCurrent := reflect.ValueOf(c)
	typeCurrent := reflect.TypeOf(c)

	valuePrev := reflect.ValueOf(prev)

	wg := sync.WaitGroup{}

	for i := 0; i < valueCurrent.Elem().NumField(); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			vFieldCurrent := valueCurrent.Elem().Field(idx)
			vFieldPrev := valuePrev.Field(idx)

			tagFielCurrent := typeCurrent.Elem().Field(idx).Tag.Get("query")

			if len(tagFielCurrent) <= 0 || tagFielCurrent == "-" {
				return
			}

			if vFieldCurrent.IsZero() {
				if strings.HasSuffix(tagFielCurrent, "omitempty") {
					return
				}
				if vFieldCurrent.Kind() == reflect.String {
					if vFieldCurrent.CanSet() {
						vFieldCurrent.SetString(vFieldPrev.Interface().(string))
					}
				}
			}
		}(i)
	}

	wg.Wait()
	return nil
}

func (c CustomerDetail) UrlValues() (url.Values, error) {
	if err := validator.New().Struct(c); err != nil {
		return url.Values{}, err
	}

	wg := sync.WaitGroup{}
	rwMutex := sync.RWMutex{}

	urlValues := url.Values{}
	valueDetail := reflect.ValueOf(c)
	typeDetail := reflect.TypeOf(c)

	for i := 0; i < valueDetail.NumField(); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			vField := valueDetail.Field(idx)
			fieldTag := typeDetail.Field(idx).Tag.Get("query")

			if len(fieldTag) > 0 && fieldTag != "-" && vField.Kind() == reflect.String {
				if strings.HasSuffix(fieldTag, "omitempty") && vField.IsZero() {
					return
				}
				queryField := strings.TrimSuffix(fieldTag, ",omitempty")
				rwMutex.Lock()
				urlValues.Add(queryField, vField.Interface().(string))
				rwMutex.Unlock()

			}
		}(i)
	}

	wg.Wait()
	return urlValues, nil
}

func (c CustomerCriteria) UrlValues() (url.Values, error) {
	if err := validator.New().Struct(c); err != nil {
		return url.Values{}, err
	}

	wg := sync.WaitGroup{}
	rwMutex := sync.RWMutex{}

	urlValues := url.Values{}
	valueCriteria := reflect.ValueOf(c)
	typeCriteria := reflect.TypeOf(c)

	for i := 0; i < valueCriteria.NumField(); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			vField := valueCriteria.Field(idx)
			tField := typeCriteria.Field(idx)
			fieldTag := tField.Tag.Get("query")

			if len(fieldTag) <= 0 {
				if vField.Kind() == reflect.Struct && vField.Type().ConvertibleTo(reflect.TypeOf(core.Criteria{})) {
					coreCiteriaData, err := vField.Interface().(core.Criteria).UrlValues()
					if err != nil {
						return
					}

					wgCriteria := sync.WaitGroup{}
					for k, v := range coreCiteriaData {
						wgCriteria.Add(1)
						go func(key string, val []string) {
							defer wgCriteria.Done()
							wgSlc := sync.WaitGroup{}
							for _, v2 := range val {
								wgSlc.Add(1)
								go func(key2, strVal string) {
									defer wgSlc.Done()
									rwMutex.Lock()
									urlValues.Add(key2, strVal)
									rwMutex.Unlock()
								}(key, v2)
							}
							wgSlc.Wait()
						}(k, v)
					}
					wgCriteria.Wait()
				}
			} else {
				if strings.HasSuffix(fieldTag, "omitempty") && vField.IsZero() {
					return
				}
				queryField := strings.TrimSuffix(fieldTag, ",omitempty")

				switch vField.Kind() {
				case reflect.Float32, reflect.Float64:
					rwMutex.Lock()
					urlValues.Add(queryField, strconv.FormatFloat(vField.Float(), 'f', 2, 64))
					rwMutex.Unlock()
				case reflect.String:
					rwMutex.Lock()
					urlValues.Add(queryField, vField.String())
					rwMutex.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()
	return urlValues, nil
}

func (r SignUpForm) UrlValues() (url.Values, error) {
	valider := validator.New()
	if err := valider.RegisterValidation("rcpassword", validatePassword); err != nil {
		return url.Values{}, err
	}
	if err := valider.Struct(r); err != nil {
		return url.Values{}, err
	}

	wg := sync.WaitGroup{}
	rwMutex := sync.RWMutex{}

	urlValues := url.Values{}
	valueForm := reflect.ValueOf(r)
	typeForm := reflect.TypeOf(r)

	for i := 0; i < valueForm.NumField(); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			vField := valueForm.Field(idx)
			tField := typeForm.Field(idx)
			fieldTag := tField.Tag.Get("query")
			if len(fieldTag) > 0 {
				if strings.HasSuffix(fieldTag, "omitempty") && vField.IsZero() {
					return
				}
				queryField := strings.TrimSuffix(fieldTag, ",omitempty")
				switch vField.Kind() {
				case reflect.String:
					rwMutex.Lock()
					urlValues.Add(queryField, vField.String())
					rwMutex.Unlock()
				case reflect.Bool:
					rwMutex.Lock()
					urlValues.Add(queryField, strconv.FormatBool(vField.Bool()))
					rwMutex.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()
	return urlValues, nil
}

func validatePassword(fl validator.FieldLevel) bool {
	return matchPasswordWithPattern(fl.Field().String(), false)
}

func matchPasswordWithPattern(password string, withRangeOfLength bool) bool {
	if withRangeOfLength && (len(password) < 9 || len(password) > 16) {
		return false
	}
	rgxAlphaLower := regexp.MustCompile(`[a-z]`)
	rgxAlphaUpper := regexp.MustCompile(`[A-Z]`)
	rgxSymbol := regexp.MustCompile(`[\~\*\!\@\$\#\%\_\+\.\?\:\,\{\}]`)
	return rgxAlphaLower.MatchString(password) && rgxAlphaUpper.MatchString(password) && rgxSymbol.MatchString(password)
}

package contact

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/mrehanabbasi/go-logicboxes/core"
)

type (
	Type                 string
	Eligibility          string
	RegistrantValidation map[string]map[Eligibility]core.JSONBool
	Detail               struct {
		ID                   string             `json:"entityid,omitempty" query:"-" sponsor:"-"`
		Type                 Type               `json:"type,omitempty" query:"type" validate:"required" sponsor:"-"`
		CustomerID           string             `json:"customerid,omitempty" query:"customer-id" validate:"required,number" sponsor:"-"`
		StatusSystem         string             `json:"currentstatus,omitempty" query:"-" sponsor:"-"`
		StatusRegistry       string             `json:"contactstatus,omitempty" query:"-" sponsor:"-"`
		ParentKey            string             `json:"parentkey,omitempty" query:"-" sponsor:"-"`
		Name                 string             `json:"name,omitempty" query:"name" validate:"required,max=255" sponsor:"name"`
		Email                string             `json:"emailaddr,omitempty" query:"email" validate:"required,email" sponsor:"email"`
		Company              string             `json:"company,omitempty" query:"company" validate:"required,max=255" sponsor:"company"`
		Address              string             `json:"address1,omitempty" query:"address-line-1" validate:"required,max=64" sponsor:"address-line-1"`     //nolint:lll
		AddressLine2         string             `json:"address2,omitempty" query:"address-line-2,optional" validate:"-" sponsor:"address-line-2,optional"` //nolint:lll
		AddressLine3         string             `json:"address3,omitempty" query:"address-line-3,optional" validate:"-" sponsor:"address-line-3,optional"` //nolint:lll
		City                 string             `json:"city,omitempty" query:"city" validate:"required,max=64" sponsor:"city"`
		State                string             `json:"state,omitempty" query:"state,optional" validate:"omitempty,max=64" sponsor:"state,optional"` //nolint:lll
		CountryCode          string             `json:"country,omitempty" query:"country" validate:"required,iso3166_1_alpha2" sponsor:"country"`
		Zipcode              string             `json:"zip,omitempty" query:"zipcode" validate:"required,max=16" sponsor:"zipcode"`
		PhoneCountryCode     string             `json:"telnocc,omitempty" query:"phone-cc" validate:"required,min=1,max=3" sponsor:"phone-cc"`
		Phone                string             `json:"telno,omitempty" query:"phone" validate:"required,min=4,max=12" sponsor:"phone"`
		FaxCountryCode       string             `json:"faxnocc,omitempty" query:"fax-cc,optional" validate:"omitempty,min=1,max=3" sponsor:"fax-cc,optional"` //nolint:lll
		Fax                  string             `json:"faxno,omitempty" query:"fax,optional" validate:"omitempty,min=4,max=12" sponsor:"fax,optional"`        //nolint:lll
		ClassName            string             `json:"classname,omitempty" query:"-" sponsor:"-"`
		ClassKey             string             `json:"classkey,omitempty" query:"-" sponsor:"-"`
		EntityActionID       string             `json:"eaqid,omitempty" query:"-" sponsor:"-"`
		ActionCompleted      core.JSONUint16    `json:"actioncompleted,omitempty" query:"-" sponsor:"-"`
		ContactID            string             `json:"contactid,omitempty" query:"-" sponsor:"-"`
		EntityTypeID         string             `json:"entitytypeid,omitempty" query:"-" sponsor:"-"`
		Description          string             `json:"description,omitempty" query:"-" sponsor:"-"`
		TimeCreation         core.JSONTime      `json:"creationdt,omitempty" validate:"-" query:"-" sponsor:"-"`
		TimeCreationRegistry core.JSONTimestamp `json:"timestamp,omitempty" validate:"-" query:"-" sponsor:"-"`
		IsDesignatedAgent    core.JSONBool      `json:"designated-agent,omitempty" validate:"-" query:"-" sponsor:"-"`
		WhoisValidity        WHOISValidity      `json:"whoisValidity,omitempty" validate:"-" query:"-" sponsor:"-"`
	}
)

type Action struct {
	ID                string `json:"eaqid,omitempty"`
	EntityID          string `json:"entityid,omitempty"`
	Type              string `json:"actiontype,omitempty"`
	Description       string `json:"actiontypedesc,omitempty"`
	Status            string `json:"actionstatus,omitempty"`
	StatusDescription string `json:"actionstatusdesc,omitempty"`
}

type Criteria struct {
	CustomerID       string              `validate:"required,number" query:"customer-id"`
	ContactIDs       []string            `validate:"omitempty,dive,number" query:"contact-id,optional"`
	Statuses         []core.EntityStatus `validate:"omitempty" query:"status,optional"`
	Name             string              `validate:"omitempty" query:"name,optional"`
	Email            string              `validate:"omitempty,email" query:"email,optional"`
	Company          string              `validate:"omitempty" query:"company,optional"`
	Type             Type                `validate:"omitempty" query:"type,optional"`
	IsIncludeInvalid bool                `validate:"omitempty" query:"include-invalid,optional"`
}

type SearchResult struct {
	RequestedLimit  uint16
	RequestedOffset uint16
	TotalMatched    int
	Contacts        []Detail
}

type WHOISValidity struct {
	IsValid     core.JSONBool `json:"valid,omitempty"`
	InvalidData []string      `json:"invalidData,omitempty"`
}

// Const for contact types.
const (
	TypeContact   Type = "Contact"
	TypeAt        Type = "AtContact"
	TypeBr        Type = "BrContact"
	TypeBrOrg     Type = "BrOrgContact"
	TypeCa        Type = "CaContact"
	TypeCl        Type = "ClContact"
	TypeCn        Type = "CnContact"
	TypeCo        Type = "CoContact"
	TypeCoop      Type = "CoopContact"
	TypeDe        Type = "DeContact"
	TypeEs        Type = "EsContact"
	TypeEu        Type = "EuContact"
	TypeFr        Type = "FrContact"
	TypeMx        Type = "MxContact"
	TypeNl        Type = "NlContact"
	TypeNyc       Type = "NycContact"
	TypeUk        Type = "UkContact"
	TypeUKService Type = "UkServiceContact"

	EligibilityDotASIA1 Eligibility = "CED_ASIAN_COUNTRY"
	EligibilityDotASIA2 Eligibility = "CED_DETAILS"
	EligibilityDotCA    Eligibility = "CPR"
	EligibilityDotCOOP  Eligibility = "SPONSORS"
	EligibilityDotES    Eligibility = "ES_CONTACT_IDENTIFICATION_DETAILS"
	EligibilityDotEU    Eligibility = "EUROPEAN_COUNTRY"
	EligibilityDotRU    Eligibility = "RU_CONTACT_INFO"
	EligibilityDotUS    Eligibility = "APP_PREF_NEXUS"
)

func (c *Criteria) URLValues() (url.Values, error) {
	if err := validator.New().Struct(c); err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	rwMutex := sync.RWMutex{}

	urlValues := url.Values{}
	valueCiteria := reflect.ValueOf(c)
	typeCiteria := reflect.TypeOf(c)

	for i := 0; i < valueCiteria.Elem().NumField(); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			vField := valueCiteria.Elem().Field(idx)
			fieldTag := typeCiteria.Elem().Field(idx).Tag.Get("query")

			if (fieldTag == "" || fieldTag == "-") || (vField.IsZero() && strings.HasSuffix(fieldTag, ",optional")) {
				return
			}

			switch vField.Kind() {
			case reflect.String:
				rwMutex.Lock()
				urlValues.Add(strings.TrimSuffix(fieldTag, ",optional"), vField.Interface().(string))
				rwMutex.Unlock()
			case reflect.Bool:
				rwMutex.Lock()
				urlValues.Add(strings.TrimSuffix(fieldTag, ",optional"), strconv.FormatBool(vField.Interface().(bool)))
				rwMutex.Unlock()
			case reflect.Slice:
				wg2 := sync.WaitGroup{}
				for j := 0; j < vField.Len(); j++ {
					wg2.Add(1)
					go func(i2 int) {
						defer wg2.Done()
						vSlice := vField.Index(i2)
						if vSlice.Kind() == reflect.String {
							var queryValue string
							if vSlice.Type().String() == "string" {
								queryValue = vSlice.Interface().(string)
							} else {
								queryValue = string(vSlice.Interface().(core.EntityStatus))
							}
							rwMutex.Lock()
							urlValues.Add(strings.TrimSuffix(fieldTag, ",optional"), queryValue)
							rwMutex.Unlock()
						}
					}(j)
				}
				wg2.Wait()
			default:
				return
			}
		}(i)
	}

	wg.Wait()
	return urlValues, nil
}

// func extractSponsorData(c ContactDetail) (*url.Values, error) {
// 	valueCurrent := reflect.ValueOf(c)
// 	typeCurrent := reflect.TypeOf(c)

// 	ret := url.Values{}
// 	for i := 0; i < valueCurrent.NumField(); i++ {
// 		vFieldCurrent := valueCurrent.Field(i)
// 		tFieldCurrent := typeCurrent.Field(i)
// 		tagFieldCurrent := tFieldCurrent.Tag.Get("sponsor")
// 		if len(tagFieldCurrent) <= 0 || tagFieldCurrent == "-" || vFieldCurrent.Kind() != reflect.String {
// 			continue
// 		}
// 		if vFieldCurrent.IsZero() {
// 			if !strings.HasSuffix(tagFieldCurrent, ",optional") {
// 				return nil, errors.New(strings.ToLower(tFieldCurrent.Name) + " must not empty")
// 			}
// 			continue
// 		}
// 		ret.Add(strings.TrimSuffix(tagFieldCurrent, ",optional"), vFieldCurrent.String())
// 	}
// 	return &ret, nil
// }

func (c *Detail) URLValues() (*url.Values, error) {
	v := validator.New()
	if err := v.Struct(c); err != nil {
		return nil, err
	}

	valueCurrent := reflect.ValueOf(c)
	typeCurrent := reflect.TypeOf(c)

	ret := &url.Values{}
	for i := 0; i < valueCurrent.Elem().NumField(); i++ {
		vFieldCurrent := valueCurrent.Elem().Field(i)
		tFieldCurrent := typeCurrent.Elem().Field(i)
		tagFieldCurrent := tFieldCurrent.Tag.Get("query")
		if tagFieldCurrent == "" || tagFieldCurrent == "-" || vFieldCurrent.Kind() != reflect.String {
			continue
		}
		if vFieldCurrent.IsZero() {
			if !strings.HasSuffix(tagFieldCurrent, ",optional") {
				return nil, errors.New(strings.ToLower(tFieldCurrent.Name) + " must not empty")
			}
			continue
		}
		ret.Add(strings.TrimSuffix(tagFieldCurrent, ",optional"), vFieldCurrent.String())
	}

	return ret, nil
}

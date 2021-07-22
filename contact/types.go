package contact

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/xpartacvs/go-resellerclub/core"
)

type ContactType string
type ContactDetail struct {
	Id                   string             `json:"entityid,omitempty" query:"-"`
	Type                 ContactType        `json:"type,omitempty" query:"type" validate:"required"`
	CustomerId           string             `json:"customerid,omitempty" query:"customer-id" validate:"required,number"`
	StatusSystem         string             `json:"currentstatus,omitempty" query:"-"`
	StatusRegistry       string             `json:"contactstatus,omitempty" query:"-"`
	ParentKey            string             `json:"parentkey,omitempty" query:"-"`
	Name                 string             `json:"name,omitempty" query:"name" validate:"required,max=255"`
	Email                string             `json:"emailaddr,omitempty" query:"email" validate:"required,email"`
	Company              string             `json:"company,omitempty" query:"company" validate:"required,max=255"`
	Address              string             `json:"address1,omitempty" query:"address-line-1" validate:"required,max=64"`
	AddressLine2         string             `json:"address2,omitempty" query:"address-line-2,optional" validate:"-"`
	AddressLine3         string             `json:"address3,omitempty" query:"address-line-3,optional" validate:"-"`
	City                 string             `json:"city,omitempty" query:"city" validate:"required,max=64"`
	State                string             `json:"state,omitempty" query:"state,optional" validate:"omitempty,max=64"`
	CountryCode          string             `json:"country,omitempty" query:"country" validate:"required,iso3166_1_alpha2"`
	Zipcode              string             `json:"zip,omitempty" query:"zipcode" validate:"required,max=16"`
	PhoneCountryCode     string             `json:"telnocc,omitempty" query:"phone-cc" validate:"required,min=1,max=3"`
	Phone                string             `json:"telno,omitempty" query:"phone" validate:"required,min=4,max=12"`
	FaxCountryCode       string             `json:"faxnocc,omitempty" query:"fax-cc,optional" validate:"omitempty,min=1,max=3"`
	Fax                  string             `json:"faxno,omitempty" query:"fax,optional" validate:"omitempty,min=4,max=12"`
	ClassName            string             `json:"classname,omitempty" query:"-"`
	ClassKey             string             `json:"classkey,omitempty" query:"-"`
	EntityActionId       string             `json:"eaqid,omitempty" query:"-"`
	ActionCompleted      core.JSONUint16    `json:"actioncompleted,omitempty" query:"-"`
	ContactId            string             `json:"contactid,omitempty" query:"-"`
	EntityTypeId         string             `json:"entitytypeid,omitempty" query:"-"`
	Description          string             `json:"description,omitempty" query:"-"`
	TimeCreation         core.JSONTime      `json:"creationdt,omitempty" validate:"-" query:"-"`
	TimeCreationRegistry core.JSONTimestamp `json:"timestamp,omitempty" validate:"-" query:"-"`
	IsDesignatedAgent    core.JSONBool      `json:"designated-agent,omitempty" validate:"-" query:"-"`
	WhoisValidity        WHOISValidity      `json:"whoisValidity,omitempty" validate:"-" query:"-"`
}

type Action struct {
	Id                string `json:"eaqid,omitempty"`
	EntityId          string `json:"entityid,omitempty"`
	Type              string `json:"actiontype,omitempty"`
	Description       string `json:"actiontypedesc,omitempty"`
	Status            string `json:"actionstatus,omitempty"`
	StatusDescription string `json:"actionstatusdesc,omitempty"`
}

type ContactCriteria struct {
	CustomerId       string              `validate:"required,number" query:"customer-id"`
	ContactIds       []string            `validate:"omitempty,dive,number" query:"contact-id,optional"`
	Statuses         []core.EntityStatus `validate:"omitempty" query:"status,optional"`
	Name             string              `validate:"omitempty" query:"name,optional"`
	Email            string              `validate:"omitempty,email" query:"email,optional"`
	Company          string              `validate:"omitempty" query:"company,optional"`
	Type             ContactType         `validate:"omitempty" query:"type,optional"`
	IsIncludeInvalid bool                `validate:"omitempty" query:"include-invalid,optional"`
}

type ContactSearchResult struct {
	RequestedLimit  uint16
	RequestedOffset uint16
	TotalMatched    int
	Contacts        []ContactDetail
}

type WHOISValidity struct {
	IsValid     core.JSONBool `json:"valid,omitempty"`
	InvalidData []string      `json:"invalidData,omitempty"`
}

const (
	TypeContact   ContactType = "Contact"
	TypeAt        ContactType = "AtContact"
	TypeBr        ContactType = "BrContact"
	TypeBrOrg     ContactType = "BrOrgContact"
	TypeCa        ContactType = "CaContact"
	TypeCl        ContactType = "ClContact"
	TypeCn        ContactType = "CnContact"
	TypeCo        ContactType = "CoContact"
	TypeCoop      ContactType = "CoopContact"
	TypeDe        ContactType = "DeContact"
	TypeEs        ContactType = "EsContact"
	TypeEu        ContactType = "EuContact"
	TypeFr        ContactType = "FrContact"
	TypeMx        ContactType = "MxContact"
	TypeNl        ContactType = "NlContact"
	TypeNyc       ContactType = "NycContact"
	TypeUk        ContactType = "UkContact"
	TypeUKService ContactType = "UkServiceContact"
)

func (c *ContactCriteria) UrlValues() (url.Values, error) {
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

			if (len(fieldTag) <= 0 || fieldTag == "-") || (vField.IsZero() && strings.HasSuffix(fieldTag, ",optional")) {
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

func (c *ContactDetail) UrlValues() (*url.Values, error) {
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
		if len(tagFieldCurrent) <= 0 || tagFieldCurrent == "-" || vFieldCurrent.Kind() != reflect.String {
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

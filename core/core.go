package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

type EntityStatus string
type AuthType string

type core struct {
	resellerId   string
	apiKey       string
	isProduction bool
}

type JSONStatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Criteria struct {
	ResellerIDs       []string  `validate:"omitempty" query:"reseller-id,omitempty"`
	CustomerIDs       []string  `validate:"omitempty" query:"customer-id,omitempty"`
	TimeCreationStart time.Time `validate:"omitempty" query:"creation-date-start,omitempty"`
	TimeCreationEnd   time.Time `validate:"omitempty" query:"creation-date-end,omitempty"`
}

type Core interface {
	CallApi(method, namespace, apiName string, data url.Values) (*http.Response, error)
	IsProduction() bool
	PrintResponse(data []byte) error
}

const (
	StatusActive              EntityStatus = "Active"
	StatusInActive            EntityStatus = "InActive"
	StatusDeleted             EntityStatus = "Deleted"
	StatusArchived            EntityStatus = "Archived"
	StatusSuspended           EntityStatus = "Suspended"
	StatusVerificationPending EntityStatus = "Pending Verification"
	StatusVerificationFailed  EntityStatus = "Failed Verification"
	StatusRestorable          EntityStatus = "Pending Delete Restorable"
	StatusNotApplicable       EntityStatus = "Not Applicable"

	AuthSms          AuthType = "sms"
	AuthGoogle       AuthType = "gauth"
	AuthGoogleBackup AuthType = "gauthbackup"
)

var (
	host = map[bool]string{
		true:  "https://httpapi.com/api",
		false: "https://test.httpapi.com/api",
	}

	RgxEmail  = regexp.MustCompile("^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")
	RgxNumber = regexp.MustCompile("^[0-9]+$")

	ErrRcApiUnsupportedMethod = errors.New("unsupported http method")
	ErrRcOperationFailed      = errors.New("operation failed")
	ErrRcInvalidCredential    = errors.New("invalid credential")
)

func (c *core) IsProduction() bool {
	return c.isProduction
}

func (c Criteria) UrlValues() (url.Values, error) {
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

			if len(fieldTag) > 0 {
				if strings.HasSuffix(fieldTag, "omitempty") && vField.IsZero() {
					return
				}
				queryField := strings.TrimSuffix(fieldTag, ",omitempty")

				switch vField.Kind() {
				case reflect.Struct:
					if vField.Type().ConvertibleTo(reflect.TypeOf(time.Time{})) {
						unixTimestamp := vField.Interface().(time.Time).Unix()
						rwMutex.Lock()
						urlValues.Add(queryField, strconv.FormatInt(unixTimestamp, 64))
						rwMutex.Unlock()
					}
				case reflect.Slice:
					wgSlice := sync.WaitGroup{}
					for j := 0; j < vField.Len(); j++ {
						wgSlice.Add(1)
						go func(x int) {
							defer wgSlice.Done()
							vSlice := vField.Index(x)
							if vSlice.Type().Kind() == reflect.String {
								rwMutex.Lock()
								urlValues.Add(queryField, vSlice.String())
								rwMutex.Unlock()
							}
						}(j)
					}
					wgSlice.Wait()
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

func (c *core) PrintResponse(data []byte) error {
	var buffer bytes.Buffer
	if err := json.Indent(&buffer, data, "", "\t"); err != nil {
		return err
	}
	buffer.WriteTo(os.Stdout)
	return nil
}

func (c *core) CallApi(method, namespace, apiName string, data url.Values) (*http.Response, error) {
	urlPath := c.createUrlPath(namespace, apiName)
	switch method {
	case http.MethodGet:
		return http.Get(urlPath + "?" + c.createRequiredQueryString() + "&" + data.Encode())
	case http.MethodPost:
		data.Add("auth-userid", c.resellerId)
		data.Add("api-key", c.apiKey)
		return http.PostForm(urlPath, data)
	}
	return nil, ErrRcApiUnsupportedMethod
}

func (c *core) createUrlPath(namespace, apiName string) string {
	return host[c.isProduction] + "/" + namespace + "/" + apiName + ".json"
}

func (c *core) createRequiredQueryString() string {
	return "auth-userid=" + url.QueryEscape(c.resellerId) + "&api-key=" + url.QueryEscape(c.apiKey)
}

func New(resellerId, apiKey string, isProduction bool) Core {
	return &core{
		resellerId:   resellerId,
		apiKey:       apiKey,
		isProduction: isProduction,
	}
}

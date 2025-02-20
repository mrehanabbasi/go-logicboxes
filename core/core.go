// Package core contains core functionality of the client.
package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

type (
	EntityStatus string
	AuthType     string
)

type Config struct {
	ResellerID   string
	APIKey       string
	IsProduction bool
}

type core struct {
	cfg    Config
	client *http.Client
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
	CallAPI(ctx context.Context, method, namespace, apiName string, data url.Values) (*http.Response, error)
	IsProduction() bool
}

// Const for status.
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
	StatusNotAvailable        EntityStatus = "NA"

	AuthSMS          AuthType = "sms"
	AuthGoogle       AuthType = "gauth"
	AuthGoogleBackup AuthType = "gauthbackup"
)

var (
	host = map[bool]string{
		true:  "https://httpapi.com/api",
		false: "https://test.httpapi.com/api",
	}

	RgxEmail  = regexp.MustCompile("^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$") //nolint:lll
	RgxNumber = regexp.MustCompile(`^\d+$`)

	ErrRcAPIUnsupportedMethod = errors.New("unsupported http method")
	ErrRcOperationFailed      = errors.New("operation failed")
	ErrRcInvalidCredential    = errors.New("invalid credential")
)

func (c *core) IsProduction() bool {
	return c.cfg.IsProduction
}

// URLValues godoc
//
//nolint:gocognit
func (c Criteria) URLValues() (url.Values, error) {
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

			if fieldTag != "" {
				if strings.HasSuffix(fieldTag, "omitempty") && vField.IsZero() {
					return
				}
				queryField := strings.TrimSuffix(fieldTag, ",omitempty")

				switch vField.Kind() {
				case reflect.Struct:
					if vField.Type().ConvertibleTo(reflect.TypeOf(time.Time{})) {
						unixTimestamp := vField.Interface().(time.Time).Unix()
						rwMutex.Lock()
						urlValues.Add(queryField, strconv.FormatInt(unixTimestamp, 10))
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

func PrintResponse(data []byte) error {
	var buffer bytes.Buffer
	if err := json.Indent(&buffer, data, "", "\t"); err != nil {
		return err
	}
	if _, err := buffer.WriteTo(os.Stdout); err != nil {
		return err
	}

	return nil
}

func (c *core) CallAPI(ctx context.Context, method, namespace, apiName string, data url.Values) (*http.Response, error) {
	urlPath := host[c.cfg.IsProduction] + "/" + namespace + "/" + apiName + ".json"
	data.Add("auth-userid", c.cfg.ResellerID)
	data.Add("api-key", c.cfg.APIKey)

	if method != http.MethodGet && method != http.MethodPost {
		return nil, ErrRcAPIUnsupportedMethod
	}

	req, err := http.NewRequestWithContext(ctx, method, urlPath+"?"+data.Encode(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return c.client.Do(req)
}

func New(cfg Config, client *http.Client) Core {
	var c *http.Client
	if client == nil {
		c = http.DefaultClient
	}

	return &core{
		cfg:    cfg,
		client: c,
	}
}

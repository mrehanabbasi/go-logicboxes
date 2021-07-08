package domain

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

type OrdersCriteria struct {
	Limit             uint16        `validate:"required,min=10,max=500" query:"no-of-records"`
	Offset            uint8         `validate:"required,min=1" query:"page-no"`
	SortOrderBy       []SortOrder   `validate:"omitempty" query:"order-by,omitempty"`
	OrderIDs          []string      `validate:"omitempty" query:"order-id,omitempty"`
	ResellerIDs       []string      `validate:"omitempty" query:"reseller-id,omitempty"`
	CustomerIDs       []string      `validate:"omitempty" query:"customer-id,omitempty"`
	DomainKeys        []DomainKey   `validate:"omitempty" query:"product-key,omitempty"`
	DomainName        string        `validate:"omitempty" query:"domain-name,omitempty"`
	OrderStatuses     []OrderStatus `validate:"omitempty" query:"status,omitempty"`
	PrivacyStatus     PrivacyState  `validate:"omitempty" query:"privacy-enabled,omitempty"`
	ShowChildOrders   bool          `validate:"omitempty" query:"show-child-orders,omitempty"`
	TimeCreationStart time.Time     `validate:"omitempty" query:"creation-date-start,omitempty"`
	TimeCreationEnd   time.Time     `validate:"omitempty" query:"creation-date-end,omitempty"`
	TimeExpiryStart   time.Time     `validate:"omitempty" query:"expiry-date-start,omitempty"`
	TimeExpiryEnd     time.Time     `validate:"omitempty" query:"expiry-date-start,omitempty"`
}

func (c OrdersCriteria) UrlValues() (url.Values, error) {
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
				case reflect.Uint8, reflect.Uint16:
					rwMutex.Lock()
					urlValues.Add(queryField, fmt.Sprintf("%d", vField.Uint()))
					rwMutex.Unlock()
				case reflect.String:
					rwMutex.Lock()
					urlValues.Add(queryField, vField.String())
					rwMutex.Unlock()
				case reflect.Bool:
					rwMutex.Lock()
					urlValues.Add(queryField, fmt.Sprintf("%t", vField.Bool()))
					rwMutex.Unlock()
				case reflect.Struct:
					if vField.Type().ConvertibleTo(reflect.TypeOf(time.Time{})) {
						timeField := vField.Interface().(time.Time)
						rwMutex.Lock()
						urlValues.Add(queryField, fmt.Sprintf("%d", timeField.Unix()))
						rwMutex.Unlock()
					}
				case reflect.Slice:
					for j := 0; j < vField.Len(); j++ {
						vSlice := vField.Index(j)
						switch vSlice.Type().Kind() {
						case reflect.String:
							rwMutex.Lock()
							urlValues.Add(queryField, vSlice.String())
							rwMutex.Unlock()
						case reflect.Map:
							if vSlice.Type().ConvertibleTo(reflect.TypeOf(SortOrder{})) {
								vSortOrder := vSlice.Interface().(SortOrder)
								var wgSortOrder sync.WaitGroup
								for k, desc := range vSortOrder {
									wgSortOrder.Add(1)
									go func(key SortBy, value bool) {
										defer wgSortOrder.Done()
										vQuery := string(key)
										if value {
											vQuery += " desc"
										}
										rwMutex.Lock()
										urlValues.Add(queryField, vQuery)
										rwMutex.Unlock()
									}(k, desc)
								}
								wgSortOrder.Wait()
							}
						}
					}
				}
			}

		}(i)
	}

	wg.Wait()
	return urlValues, nil
}

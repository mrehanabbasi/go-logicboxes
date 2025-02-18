// Package pricing contains API methods for pricing.
package pricing

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mrehanabbasi/go-logicboxes/core"
)

type Pricing interface {
	GettingCustomerPricing(customerID string) (CustomerPrice, error)
	GettingResellerPricing(resellerID string) (ResellerPrice, error)
	GettingResellerCostPricing(resellerID string) (ResellerCostPrice, error)
	GettingPromoPrices() (PromoPrice, error)
}

func New(c core.Core) Pricing {
	return &pricing{c}
}

type pricing struct {
	core core.Core
}

func (p *pricing) GettingCustomerPricing(customerID string) (CustomerPrice, error) {
	data := make(url.Values)
	data.Add("customer-id", customerID)

	resp, err := p.core.CallAPI(http.MethodGet, "products", "customer-price", data)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result CustomerPrice
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *pricing) GettingResellerPricing(resellerID string) (ResellerPrice, error) {
	data := make(url.Values)
	data.Add("reseller-id", resellerID)

	resp, err := p.core.CallAPI(http.MethodGet, "products", "reseller-price", data)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result ResellerPrice
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *pricing) GettingResellerCostPricing(resellerID string) (ResellerCostPrice, error) {
	data := make(url.Values)
	data.Add("reseller-id", resellerID)

	resp, err := p.core.CallAPI(http.MethodGet, "products", "reseller-cost-price", data)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result ResellerCostPrice
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *pricing) GettingPromoPrices() (PromoPrice, error) {
	data := make(url.Values)

	resp, err := p.core.CallAPI(http.MethodGet, "products", "promo-details", data)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result PromoPrice
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

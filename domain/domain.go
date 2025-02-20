// Package domain contains APIs for domain management.
package domain

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mrehanabbasi/go-logicboxes/core"
)

type domain struct {
	core core.Core
}

type Domain interface {
	CheckAvailability(ctx context.Context, domainsWithoutTLD, tlds []string) (Availabilities, error)
	SuggestNames(ctx context.Context, keyword, tldOnly string, exactMatch, adult bool) (SuggestNames, error)
	Register(
		ctx context.Context,
		domainName string,
		years int,
		ns []string,
		customerID, regContactID, adminContactID, techContactID, billingContactID, invoiceOption string,
		purchasePrivacy, protectPrivacy, autoRenew bool,
		attrName, attrValue string,
		discountAmount float64,
		purchasePremiumDNS bool,
	) (*RegisterResponse, error)
	Transfer(
		ctx context.Context,
		domainName, authCode, customerID, regContactID, adminContactID, techContactID, billingContactID, invoiceOption string,
		purchasePrivacy, protectPrivacy, autoRenew bool,
		ns []string,
		attrName, attrValue string,
		purchasePremiumDNS bool,
	) (*RegisterResponse, error)
	ValidatingTransferRequest(ctx context.Context, domainName string) (bool, error)
	GetCustomerDefaultNameServers(ctx context.Context, customerID string) ([]string, error)
	GetOrderID(ctx context.Context, domainName string) (string, error)
	GetRegistrationOrderDetails(ctx context.Context, orderID string, options []string) (*OrderDetail, error)
	ModifyNameServers(ctx context.Context, orderID string, ns []string) (*NameServersResponse, error)
	AddChildNameServer(ctx context.Context, orderID, cns string, ips []string) (*NameServersResponse, error)
	ModifyChildNameServerHostName(ctx context.Context, orderID, oldCNS, newCNS string) (*NameServersResponse, error)
	ModifyChildNameServerIPAddress(ctx context.Context, orderID, cns, oldIP, newIP string) (*NameServersResponse, error)
	DeletingChildNameServerIPAddress(ctx context.Context, orderID, cns string, ips []string) (*NameServersResponse, error)
	ModifyContacts(
		ctx context.Context,
		orderID, regContactID, adminContactID, techContactID, billingContactID string,
		sixtyDayLockOptout, designatedAgent bool,
		attrName, attrValue string,
	) (*ModifyAuthCodeResponse, error)
	ModifyPrivacyProtectionStatus(
		ctx context.Context,
		orderID string,
		protectPrivacy bool,
		reason string,
	) (*ModifyPrivacyProtectionStatusResponse, error)
	ModifyAuthCode(ctx context.Context, orderID, authCode string) (*ModifyAuthCodeResponse, error)
	ApplyTheftProtectionLock(ctx context.Context, orderID string) (*TheftProtectionLockResponse, error)
	RemoveTheftProtectionLock(ctx context.Context, orderID string) (*TheftProtectionLockResponse, error)
	GetTheListOfLocksAppliedOnDomainName(ctx context.Context, orderID string) (*GetTheListOfLocksAppliedOnDomainNameResponse, error)
	CancelTransfer(ctx context.Context, orderID string) (*CancelTransferResponse, error)
	Suspend(ctx context.Context, orderID, reason string) (*TheftProtectionLockResponse, error)
	Unsuspend(ctx context.Context, orderID string) (*TheftProtectionLockResponse, error)
	Delete(ctx context.Context, orderID string) (*DeleteResponse, error)
}

func New(c core.Core) Domain {
	return &domain{c}
}

func (d *domain) CheckAvailability(ctx context.Context, domainName, tlds []string) (Availabilities, error) {
	if len(domainName) == 0 || len(tlds) == 0 {
		return Availabilities{}, errors.New("domainnames and tlds must not empty")
	}

	data := url.Values{}
	data["domain-name"] = append(data["domain-name"], domainName...)
	data["tlds"] = append(data["tlds"], tlds...)

	resp, err := d.core.CallAPI(ctx, http.MethodGet, "domains", "available", data)
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

	availabilities := Availabilities{}
	if err := json.Unmarshal(bytesResp, &availabilities); err != nil {
		return nil, err
	}

	return availabilities, nil
}

func (d *domain) SuggestNames(ctx context.Context, keyword, tldOnly string, exactMatch, adult bool) (SuggestNames, error) {
	data := make(url.Values)
	data.Add("keyword", keyword)
	data.Add("tld-only", tldOnly)
	data.Add("exact-match", strconv.FormatBool(exactMatch))
	data.Add("adult", strconv.FormatBool(adult))

	resp, err := d.core.CallAPI(ctx, http.MethodGet, "domains/v5", "suggest-names", data)
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
		err := json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	suggestNames := SuggestNames{}
	if err := json.Unmarshal(bytesResp, &suggestNames); err != nil {
		return nil, err
	}

	return suggestNames, nil
}

func (d *domain) Register(
	ctx context.Context,
	domainName string,
	years int,
	ns []string,
	customerID, regContactID, adminContactID, techContactID, billingContactID, invoiceOption string,
	purchasePrivacy, protectPrivacy, autoRenew bool,
	attrName, attrValue string,
	discountAmount float64,
	purchasePremiumDNS bool,
) (*RegisterResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("years", strconv.Itoa(years))
	data["ns"] = append(data["ns"], ns...)
	data.Add("customer-id", customerID)
	data.Add("reg-contact-id", regContactID)
	data.Add("admin-contact-id", adminContactID)
	data.Add("tech-contact-id", techContactID)
	data.Add("billing-contact-id", billingContactID)
	data.Add("invoice-option", invoiceOption)
	data.Add("purchase-privacy", strconv.FormatBool(purchasePrivacy))
	data.Add("protect-privacy", strconv.FormatBool(protectPrivacy))
	data.Add("auto-renew", strconv.FormatBool(autoRenew))
	data.Add("attr-name", attrName)
	data.Add("attr-value", attrValue)
	data.Add("discount-amount", strconv.FormatFloat(discountAmount, 'f', 2, 64))
	data.Add("purchase-premium-dns", strconv.FormatBool(purchasePremiumDNS))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "register", data)
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

	var result RegisterResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) Transfer(
	ctx context.Context,
	domainName, authCode, customerID, regContactID, adminContactID, techContactID, billingContactID, invoiceOption string,
	purchasePrivacy, protectPrivacy, autoRenew bool,
	ns []string,
	attrName, attrValue string,
	purchasePremiumDNS bool,
) (*RegisterResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("auth-code", authCode)
	data.Add("customer-id", customerID)
	data.Add("reg-contact-id", regContactID)
	data.Add("admin-contact-id", adminContactID)
	data.Add("tech-contact-id", techContactID)
	data.Add("billing-contact-id", billingContactID)
	data.Add("invoice-option", invoiceOption)
	data.Add("purchase-privacy", strconv.FormatBool(purchasePrivacy))
	data.Add("protect-privacy", strconv.FormatBool(protectPrivacy))
	data.Add("auto-renew", strconv.FormatBool(autoRenew))
	data["ns"] = append(data["ns"], ns...)
	data.Add("attr-name", attrName)
	data.Add("attr-value", attrValue)
	data.Add("purchase-premium-dns", strconv.FormatBool(purchasePremiumDNS))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "transfer", data)
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

	var result RegisterResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ValidatingTransferRequest(ctx context.Context, domainName string) (bool, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "validate-transfer", data)
	if err != nil {
		return false, err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return false, err
		}
		return false, errors.New(strings.ToLower(errResponse.Message))
	}

	var result bool
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (d *domain) Renew(
	ctx context.Context,
	orderID string,
	_, expDate int, // years & expDate
	purchasePrivacy, autoRenew bool,
	invoiceOption string,
	discountAmount float64,
	purchasePremiumDNS bool,
) error {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("exp-date", strconv.Itoa(expDate))
	data.Add("purchase-privacy", strconv.FormatBool(purchasePrivacy))
	data.Add("auto-renew", strconv.FormatBool(autoRenew))
	data.Add("invoice-option", invoiceOption)
	data.Add("discount-amount", strconv.FormatFloat(discountAmount, 'f', 2, 64))
	data.Add("purchase-premium-dns", strconv.FormatBool(purchasePremiumDNS))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "renew", data)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) SearchOrders(ctx context.Context, criteria OrderCriteria) error {
	urlValues, err := criteria.URLValues()
	if err != nil {
		return err
	}
	resp, err := d.core.CallAPI(ctx, http.MethodGet, "domains", "search", urlValues)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) GetCustomerDefaultNameServers(ctx context.Context, customerID string) ([]string, error) {
	data := make(url.Values)
	data.Add("customer-id", customerID)

	resp, err := d.core.CallAPI(ctx, http.MethodGet, "domains", "customer-default-ns", data)
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

	result := make([]string, 0)
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (d *domain) GetOrderID(ctx context.Context, domainName string) (string, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)

	resp, err := d.core.CallAPI(ctx, http.MethodGet, "domains", "orderid", data)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return "", err
		}
		return "", errors.New(strings.ToLower(errResponse.Message))
	}

	return string(bytesResp), nil
}

func (d *domain) GetRegistrationOrderDetails(ctx context.Context, orderID string, options []string) (*OrderDetail, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data["options"] = append(data["options"], options...)

	resp, err := d.core.CallAPI(ctx, http.MethodGet, "domains", "details", data)
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

	var orderDetail OrderDetail
	if err := json.Unmarshal(bytesResp, &orderDetail); err != nil {
		return nil, err
	}

	return &orderDetail, nil
}

func (d *domain) ModifyNameServers(ctx context.Context, orderID string, ns []string) (*NameServersResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data["ns"] = append(data["ns"], ns...)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "modify-ns", data)
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

	var result NameServersResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) AddChildNameServer(ctx context.Context, orderID, cns string, ips []string) (*NameServersResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("cns", cns)
	data["ip"] = append(data["ip"], ips...)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "add-cns", data)
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

	var result NameServersResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ModifyChildNameServerHostName(ctx context.Context, orderID, oldCNS, newCNS string) (*NameServersResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("old-cns", oldCNS)
	data.Add("new-cns", newCNS)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "modify-cns-name", data)
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

	var result NameServersResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ModifyChildNameServerIPAddress(ctx context.Context, orderID, cns, oldIP, newIP string) (*NameServersResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("cns", cns)
	data.Add("old-ip", oldIP)
	data.Add("new-ip", newIP)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "modify-cns-ip", data)
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

	var result NameServersResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) DeletingChildNameServerIPAddress(ctx context.Context, orderID, cns string, ips []string) (*NameServersResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("cns", cns)
	data["ip"] = append(data["ip"], ips...)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "delete-cns-ip", data)
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

	var result NameServersResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ModifyContacts(
	ctx context.Context,
	orderID, regContactID, adminContactID, techContactID, billingContactID string,
	sixtyDayLockOptout, designatedAgent bool,
	attrName, attrValue string,
) (*ModifyAuthCodeResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("reg-contact-id", regContactID)
	data.Add("admin-contact-id", adminContactID)
	data.Add("tech-contact-id", techContactID)
	data.Add("billing-contact-id", billingContactID)
	data.Add("sixty-day-lock-optout", strconv.FormatBool(sixtyDayLockOptout))
	data.Add("designated-agent", strconv.FormatBool(designatedAgent))
	data.Add("attr-name", attrName)
	data.Add("attr-value", attrValue)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "modify-contact", data)
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

	var result ModifyAuthCodeResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ModifyPrivacyProtectionStatus(
	ctx context.Context,
	orderID string,
	protectPrivacy bool,
	reason string,
) (*ModifyPrivacyProtectionStatusResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("protect-privacy", strconv.FormatBool(protectPrivacy))
	data.Add("reason", reason)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "modify-privacy-protection", data)
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

	var result ModifyPrivacyProtectionStatusResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ModifyAuthCode(ctx context.Context, orderID, authCode string) (*ModifyAuthCodeResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("auth-code", authCode)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "modify-auth-code", data)
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

	var result ModifyAuthCodeResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ApplyTheftProtectionLock(ctx context.Context, orderID string) (*TheftProtectionLockResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "enable-theft-protection", data)
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

	var result TheftProtectionLockResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) RemoveTheftProtectionLock(ctx context.Context, orderID string) (*TheftProtectionLockResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "disable-theft-protection", data)
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

	var result TheftProtectionLockResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) GetTheListOfLocksAppliedOnDomainName(
	ctx context.Context,
	orderID string,
) (*GetTheListOfLocksAppliedOnDomainNameResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallAPI(ctx, http.MethodGet, "domains", "locks", data)
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

	var result GetTheListOfLocksAppliedOnDomainNameResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ModifyTELWhoisPreference(ctx context.Context, orderID, whoisType, publish string) error {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("whois-type", whoisType)
	data.Add("publish", publish)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "tel/modify-whois-pref", data)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) ResendTransferApprovalMail(ctx context.Context, orderID string) error {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "resend-rfa", data)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) ReleaseUKDomainName(ctx context.Context, orderID, newTag string) error {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("new-tag", newTag)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "uk/release", data)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) CancelTransfer(ctx context.Context, orderID string) (*CancelTransferResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "cancel-transfer", data)
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

	var result CancelTransferResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) Suspend(ctx context.Context, orderID, reason string) (*TheftProtectionLockResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("reason", reason)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "orders", "suspend", data)
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

	var result TheftProtectionLockResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) Unsuspend(ctx context.Context, orderID string) (*TheftProtectionLockResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "orders", "unsuspend", data)
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

	var result TheftProtectionLockResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) Delete(ctx context.Context, orderID string) (*DeleteResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "delete", data)
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

	var result DeleteResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) Restore(ctx context.Context, orderID, invoiceOption string) error {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("invoice-option", invoiceOption)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "restore", data)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) RecheckingNSWithDERegistry(ctx context.Context, orderID string) error {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "de/recheck-ns", data)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) AssociatingOrDissociatingXXXMembershipTokenID(ctx context.Context, orderID, associationID string) error {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("association-id", associationID)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "domains", "dotxxx/association-details", data)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

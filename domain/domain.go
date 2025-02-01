package domain

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mrehanabbasi/go-resellerclub/core"
)

type domain struct {
	core core.Core
}

type Domain interface {
	CheckAvailability(domainsWithoutTLD, tlds []string) (DomainAvailabilities, error)
	SuggestNames(keyword, tldOnly string, exactMatch, adult bool) (SuggestNames, error)
	Register(domainName string, years int, ns []string, customerID, regContactID, adminContactID, techContactID, billingContactID, invoiceOption string, purchasePrivacy bool, protectPrivacy bool, autoRenew bool, attrName, attrValue string, discountAmount float64, purchasePremiumDNS bool) (*RegisterResponse, error)
	Transfer(domainName, authCode, customerID, regContactID, adminContactID, techContactID, billingContactID, invoiceOption string, purchasePrivacy, protectPrivacy, autoRenew bool, ns []string, attrName, attrValue string, purchasePremiumDNS bool) (*RegisterResponse, error)
	ValidatingTransferRequest(domainName string) (bool, error)
	GetCustomerDefaultNameServers(customerID string) ([]string, error)
	GetOrderID(domainName string) (string, error)
	GetRegistrationOrderDetails(orderID string, options []string) (*OrderDetail, error)
	ModifyNameServers(orderID string, ns []string) (*NameServersResponse, error)
	AddChildNameServer(orderID, cns string, ips []string) (*NameServersResponse, error)
	ModifyChildNameServerHostName(orderID, oldCNS, newCNS string) (*NameServersResponse, error)
	ModifyChildNameServerIPAddress(orderID, cns, oldIP, newIP string) (*NameServersResponse, error)
	DeletingChildNameServerIPAddress(orderID, cns string, ips []string) (*NameServersResponse, error)
	ModifyContacts(orderID, regContactID, adminContactID, techContactID, billingContactID string, sixtyDayLockOptout, designatedAgent bool, attrName, attrValue string) (*ModifyAuthCodeResponse, error)
	ModifyPrivacyProtectionStatus(orderID string, protectPrivacy bool, reason string) (*ModifyPrivacyProtectionStatusResponse, error)
	ModifyAuthCode(orderID, authCode string) (*ModifyAuthCodeResponse, error)
	ApplyTheftProtectionLock(orderID string) (*TheftProtectionLockResponse, error)
	RemoveTheftProtectionLock(orderID string) (*TheftProtectionLockResponse, error)
	GetTheListOfLocksAppliedOnDomainName(orderID string) (*GetTheListOfLocksAppliedOnDomainNameResponse, error)
	CancelTransfer(orderID string) (*CancelTransferResponse, error)
	Suspend(orderID, reason string) (*TheftProtectionLockResponse, error)
	Unsuspend(orderID string) (*TheftProtectionLockResponse, error)
	Delete(orderID string) (*DeleteResponse, error)
}

func New(c core.Core) Domain {
	return &domain{c}
}

func (d *domain) CheckAvailability(domainName, tlds []string) (DomainAvailabilities, error) {
	if len(domainName) <= 0 || len(tlds) <= 0 {
		return DomainAvailabilities{}, errors.New("domainnames and tlds must not empty")
	}

	data := url.Values{}
	data["domain-name"] = append(data["domain-name"], domainName...)
	data["tlds"] = append(data["tlds"], tlds...)

	resp, err := d.core.CallApi(http.MethodGet, "domains", "available", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	availabilities := DomainAvailabilities{}
	err = json.Unmarshal(bytesResp, &availabilities)
	if err != nil {
		return nil, err
	}

	return availabilities, nil
}

func (d *domain) SuggestNames(keyword, tldOnly string, exactMatch, adult bool) (SuggestNames, error) {
	data := make(url.Values)
	data.Add("keyword", keyword)
	data.Add("tld-only", tldOnly)
	data.Add("exact-match", strconv.FormatBool(exactMatch))
	data.Add("adult", strconv.FormatBool(adult))

	resp, err := d.core.CallApi(http.MethodGet, "domains/v5", "suggest-names", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	suggestNames := SuggestNames{}
	err = json.Unmarshal(bytesResp, &suggestNames)
	if err != nil {
		return nil, err
	}

	return suggestNames, nil
}

func (d *domain) Register(domainName string, years int, ns []string, customerID, regContactID, adminContactID, techContactID, billingContactID, invoiceOption string, purchasePrivacy bool, protectPrivacy bool, autoRenew bool, attrName, attrValue string, discountAmount float64, purchasePremiumDNS bool) (*RegisterResponse, error) {
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

	resp, err := d.core.CallApi(http.MethodPost, "domains", "register", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result RegisterResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) Transfer(domainName, authCode, customerID, regContactID, adminContactID, techContactID, billingContactID, invoiceOption string, purchasePrivacy, protectPrivacy, autoRenew bool, ns []string, attrName, attrValue string, purchasePremiumDNS bool) (*RegisterResponse, error) {
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

	resp, err := d.core.CallApi(http.MethodPost, "domains", "transfer", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result RegisterResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ValidatingTransferRequest(domainName string) (bool, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "validate-transfer", data)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return false, err
		}
		return false, errors.New(strings.ToLower(errResponse.Message))
	}

	var result bool
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (d *domain) Renew(orderID string, years, expDate int, purchasePrivacy, autoRenew bool, invoiceOption string, discountAmount float64, purchasePremiumDNS bool) error {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("exp-date", strconv.Itoa(expDate))
	data.Add("purchase-privacy", strconv.FormatBool(purchasePrivacy))
	data.Add("auto-renew", strconv.FormatBool(autoRenew))
	data.Add("invoice-option", invoiceOption)
	data.Add("discount-amount", strconv.FormatFloat(discountAmount, 'f', 2, 64))
	data.Add("purchase-premium-dns", strconv.FormatBool(purchasePremiumDNS))

	resp, err := d.core.CallApi(http.MethodPost, "domains", "renew", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) SearchOrders(criteria OrderCriteria) error {
	urlValues, err := criteria.UrlValues()
	if err != nil {
		return err
	}
	resp, err := d.core.CallApi(http.MethodGet, "domains", "search", urlValues)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) GetCustomerDefaultNameServers(customerID string) ([]string, error) {
	data := make(url.Values)
	data.Add("customer-id", customerID)

	resp, err := d.core.CallApi(http.MethodGet, "domains", "customer-default-ns", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	result := make([]string, 0)
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *domain) GetOrderID(domainName string) (string, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)

	resp, err := d.core.CallApi(http.MethodGet, "domains", "orderid", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return "", err
		}
		return "", errors.New(strings.ToLower(errResponse.Message))
	}

	return string(bytesResp), nil
}

func (d *domain) GetRegistrationOrderDetails(orderID string, options []string) (*OrderDetail, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data["options"] = append(data["options"], options...)

	resp, err := d.core.CallApi(http.MethodGet, "domains", "details", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var orderDetail OrderDetail
	err = json.Unmarshal(bytesResp, &orderDetail)
	if err != nil {
		return nil, err
	}

	return &orderDetail, nil
}

func (d *domain) ModifyNameServers(orderID string, ns []string) (*NameServersResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data["ns"] = append(data["ns"], ns...)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "modify-ns", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result NameServersResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) AddChildNameServer(orderID, cns string, ips []string) (*NameServersResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("cns", cns)
	data["ip"] = append(data["ip"], ips...)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "add-cns", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result NameServersResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ModifyChildNameServerHostName(orderID, oldCNS, newCNS string) (*NameServersResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("old-cns", oldCNS)
	data.Add("new-cns", newCNS)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "modify-cns-name", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result NameServersResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ModifyChildNameServerIPAddress(orderID, cns, oldIP, newIP string) (*NameServersResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("cns", cns)
	data.Add("old-ip", oldIP)
	data.Add("new-ip", newIP)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "modify-cns-ip", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result NameServersResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) DeletingChildNameServerIPAddress(orderID, cns string, ips []string) (*NameServersResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("cns", cns)
	data["ip"] = append(data["ip"], ips...)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "delete-cns-ip", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result NameServersResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ModifyContacts(orderID, regContactID, adminContactID, techContactID, billingContactID string, sixtyDayLockOptout, designatedAgent bool, attrName, attrValue string) (*ModifyAuthCodeResponse, error) {
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

	resp, err := d.core.CallApi(http.MethodPost, "domains", "modify-contact", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result ModifyAuthCodeResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ModifyPrivacyProtectionStatus(orderID string, protectPrivacy bool, reason string) (*ModifyPrivacyProtectionStatusResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("protect-privacy", strconv.FormatBool(protectPrivacy))
	data.Add("reason", reason)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "modify-privacy-protection", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result ModifyPrivacyProtectionStatusResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ModifyAuthCode(orderID, authCode string) (*ModifyAuthCodeResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("auth-code", authCode)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "modify-auth-code", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result ModifyAuthCodeResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ApplyTheftProtectionLock(orderID string) (*TheftProtectionLockResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "enable-theft-protection", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result TheftProtectionLockResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) RemoveTheftProtectionLock(orderID string) (*TheftProtectionLockResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "disable-theft-protection", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result TheftProtectionLockResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) GetTheListOfLocksAppliedOnDomainName(orderID string) (*GetTheListOfLocksAppliedOnDomainNameResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallApi(http.MethodGet, "domains", "locks", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result GetTheListOfLocksAppliedOnDomainNameResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) ModifyTELWhoisPreference(orderID, whoisType, publish string) error {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("whois-type", whoisType)
	data.Add("publish", publish)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "tel/modify-whois-pref", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) ResendTransferApprovalMail(orderID string) error {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "resend-rfa", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) ReleaseUKDomainName(orderID, newTag string) error {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("new-tag", newTag)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "uk/release", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) CancelTransfer(orderID string) (*CancelTransferResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "cancel-transfer", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result CancelTransferResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) Suspend(orderID, reason string) (*TheftProtectionLockResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("reason", reason)

	resp, err := d.core.CallApi(http.MethodPost, "orders", "suspend", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result TheftProtectionLockResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) Unsuspend(orderID string) (*TheftProtectionLockResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallApi(http.MethodPost, "orders", "unsuspend", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result TheftProtectionLockResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) Delete(orderID string) (*DeleteResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "delete", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(strings.ToLower(errResponse.Message))
	}

	var result DeleteResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domain) Restore(orderID, invoiceOption string) error {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("invoice-option", invoiceOption)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "restore", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) RecheckingNSWithDERegistry(orderID string) error {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "de/recheck-ns", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

func (d *domain) AssociatingOrDissociatingXXXMembershipTokenID(orderID, associationID string) error {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("association-id", associationID)

	resp, err := d.core.CallApi(http.MethodPost, "domains", "dotxxx/association-details", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return err
		}
		return errors.New(strings.ToLower(errResponse.Message))
	}

	return nil
}

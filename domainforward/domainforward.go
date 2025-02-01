package domainforward

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

type DomainForward interface {
	ActivatingDomainForwardingService(orderID, subDomainPrefix, forwardTo string, urlMasking bool, metaTags, noframes string, subDomainForwarding, pathForwarding bool) (*StdResponse, error)
	GettingDetailsDomainForwardingService(orderID string, includeSubdomain bool) (*DetailsDomainForward, error)
	ManagingDomainForwardingService(orderID, subDomainPrefix, forwardTo string, urlMasking bool, metaTags, noframes string, subDomainForwarding, pathForwarding bool) (*StdResponse, error)
	GettingDNSRecords(domainName string) ([]*DNSRecord, error)
	RemoveDomainForwardingForDomain(domainName string) (bool, error)
	DisableDomainForwardingForSubDomain(orderID, subDomainPrefix string) (bool, error)
}

func New(c core.Core) DomainForward {
	return &domainForward{c}
}

type domainForward struct {
	core core.Core
}

func (d *domainForward) ActivatingDomainForwardingService(orderID, subDomainPrefix, forwardTo string, urlMasking bool, metaTags, noframes string, subDomainForwarding, pathForwarding bool) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("sub-domain-prefix", subDomainPrefix)
	data.Add("forward-to", forwardTo)
	data.Add("url-masking", strconv.FormatBool(urlMasking))
	data.Add("meta-tags", metaTags)
	data.Add("noframes", noframes)
	data.Add("sub-domain-forwarding", strconv.FormatBool(subDomainForwarding))
	data.Add("path-forwarding", strconv.FormatBool(pathForwarding))

	resp, err := d.core.CallApi(http.MethodPost, "domainforward", "activate", data)
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

	var result StdResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domainForward) GettingDetailsDomainForwardingService(orderID string, includeSubdomain bool) (*DetailsDomainForward, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("include-subdomain", strconv.FormatBool(includeSubdomain))

	resp, err := d.core.CallApi(http.MethodGet, "domainforward", "details", data)
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

	var result DetailsDomainForward
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domainForward) ManagingDomainForwardingService(orderID, subDomainPrefix, forwardTo string, urlMasking bool, metaTags, noframes string, subDomainForwarding, pathForwarding bool) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("sub-domain-prefix", subDomainPrefix)
	data.Add("forward-to", forwardTo)
	data.Add("url-masking", strconv.FormatBool(urlMasking))
	data.Add("meta-tags", metaTags)
	data.Add("noframes", noframes)
	data.Add("sub-domain-forwarding", strconv.FormatBool(subDomainForwarding))
	data.Add("path-forwarding", strconv.FormatBool(pathForwarding))

	resp, err := d.core.CallApi(http.MethodPost, "domainforward", "manage", data)
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

	var result StdResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *domainForward) GettingDNSRecords(domainName string) ([]*DNSRecord, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)

	resp, err := d.core.CallApi(http.MethodGet, "domainforward", "dns-records", data)
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

	var result []*DNSRecord
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *domainForward) RemoveDomainForwardingForDomain(domainName string) (bool, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)

	resp, err := d.core.CallApi(http.MethodPost, "domainforward", "delete", data)
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

func (d *domainForward) DisableDomainForwardingForSubDomain(orderID, subDomainPrefix string) (bool, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)
	data.Add("sub-domain-prefix", subDomainPrefix)

	resp, err := d.core.CallApi(http.MethodPost, "domainforward", "sub-domain-record/delete", data)
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

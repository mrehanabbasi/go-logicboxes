// Package dns contains APIs for dns.
package dns

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mrehanabbasi/go-resellerclub/core"
)

type DNS interface {
	ActivatingDNSService(orderID string) (*ActivatingDNSServiceResponse, error)
	AddingIPv4AddressRecord(domainName, value, host string, ttl int) (*StdResponse, error)
	AddingIPv6AddressRecord(domainName, value, host string, ttl int) (*StdResponse, error)
	AddingCNAMERecord(domainName, value, host string, ttl int) (*StdResponse, error)
	AddingMXRecord(domainName, value, host string, ttl, priority int) (*StdResponse, error)
	AddingNSRecord(domainName, value, host string, ttl int) (*StdResponse, error)
	AddingTXTRecord(domainName, value, host string, ttl int) (*StdResponse, error)
	AddingSRVRecord(domainName, value, host string, ttl, priority, port, weight int) (*StdResponse, error)
	ModifyingIPv4AddressRecord(domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error)
	ModifyingIPv6AddressRecord(domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error)
	ModifyingCNAMERecord(domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error)
	ModifyingMXRecord(domainName, host, currentValue, newValue string, ttl, priority int) (*StdResponse, error)
	ModifyingNSRecord(domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error)
	ModifyingTXTRecord(domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error)
	ModifyingSRVRecord(domainName, host, currentValue, newValue string, ttl, priority, port, weight int) (*StdResponse, error)
	ModifyingSOARecord(domainName, responsiblePerson string, refresh, retry, expire, ttl int) (*StdResponse, error)
	SearchingDNSRecords(domainName string, typeRecord RecordType, noOfRecords, pageNo int, host, value string) (*SearchingDNSRecords, error)
	DeletingDNSRecord(host, value string) (*StdResponse, error)
	DeletingIPv4AddressRecord(domainName, host, value string) (*StdResponse, error)
	DeletingIPv6AddressRecord(domainName, host, value string) (*StdResponse, error)
	DeletingCNAMERecord(domainName, host, value string) (*StdResponse, error)
	DeletingMXRecord(domainName, host, value string) (*StdResponse, error)
	DeletingNSRecord(domainName, host, value string) (*StdResponse, error)
	DeletingTXTRecord(domainName, host, value string) (*StdResponse, error)
	DeletingSRVRecord(domainName, host, value string, port, weight int) (*StdResponse, error)
}

func New(c core.Core) DNS {
	return &dns{c}
}

type dns struct {
	core core.Core
}

func (d *dns) ActivatingDNSService(orderID string) (*ActivatingDNSServiceResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "activate", data)
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

	var result ActivatingDNSServiceResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) AddingIPv4AddressRecord(domainName, value, host string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/add-ipv4-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) AddingIPv6AddressRecord(domainName, value, host string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/add-ipv6-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) AddingCNAMERecord(domainName, value, host string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/add-cname-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) AddingMXRecord(domainName, value, host string, ttl, priority int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))
	data.Add("priority", strconv.Itoa(priority))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/add-mx-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) AddingNSRecord(domainName, value, host string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/add-ns-record", data)
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

	var result StdResponse
	err = json.Unmarshal(bytesResp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) AddingTXTRecord(domainName, value, host string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/manage/add-ns-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) AddingSRVRecord(domainName, value, host string, ttl, priority, port, weight int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))
	data.Add("priority", strconv.Itoa(priority))
	data.Add("port", strconv.Itoa(port))
	data.Add("weight", strconv.Itoa(weight))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/add-srv-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) ModifyingIPv4AddressRecord(domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/update-ipv4-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) ModifyingIPv6AddressRecord(domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/update-ipv6-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) ModifyingCNAMERecord(domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/update-cname-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) ModifyingMXRecord(domainName, host, currentValue, newValue string, ttl, priority int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))
	data.Add("priority", strconv.Itoa(priority))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/update-mx-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) ModifyingNSRecord(domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/update-ns-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) ModifyingTXTRecord(domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/update-txt-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) ModifyingSRVRecord(domainName, host, currentValue, newValue string, ttl, priority, port, weight int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))
	data.Add("priority", strconv.Itoa(priority))
	data.Add("port", strconv.Itoa(port))
	data.Add("weight", strconv.Itoa(weight))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/update-srv-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) ModifyingSOARecord(domainName, responsiblePerson string, refresh, retry, expire, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("responsible-person", responsiblePerson)
	data.Add("refresh", strconv.Itoa(refresh))
	data.Add("retry", strconv.Itoa(retry))
	data.Add("expire", strconv.Itoa(expire))
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/update-soa-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) SearchingDNSRecords(
	domainName string,
	typeRecord RecordType,
	noOfRecords, pageNo int,
	host, value string,
) (*SearchingDNSRecords, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("type", string(typeRecord))
	data.Add("no-of-records", strconv.Itoa(noOfRecords))
	data.Add("page-no", strconv.Itoa(pageNo))
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(http.MethodGet, "dns", "manage/search-records", data)
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

	var records SearchingDNSRecords
	var result map[string]interface{}
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	for k, v := range result {
		switch k {
		case "recsonpage":
			records.RecsOnPage = fmt.Sprintf("%v", v)
			continue
		case "recsindb":
			records.Recsindb = fmt.Sprintf("%v", v)
			continue
		}

		_, err := strconv.Atoi(k)
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		var record Record
		if err := json.Unmarshal(b, &record); err != nil {
			return nil, err
		}

		records.Records = append(records.Records, &record)
	}

	return &records, nil
}

func (d *dns) DeletingDNSRecord(host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(http.MethodGet, "dns", "manage/delete-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) DeletingIPv4AddressRecord(domainName, host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/delete-ipv4-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) DeletingIPv6AddressRecord(domainName, host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/delete-ipv6-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) DeletingCNAMERecord(domainName, host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/delete-cname-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) DeletingMXRecord(domainName, host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/delete-mx-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) DeletingNSRecord(domainName, host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/delete-ns-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) DeletingTXTRecord(domainName, host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/delete-txt-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *dns) DeletingSRVRecord(domainName, host, value string, port, weight int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)
	data.Add("port", strconv.Itoa(port))
	data.Add("weight", strconv.Itoa(weight))

	resp, err := d.core.CallAPI(http.MethodPost, "dns", "manage/delete-srv-record", data)
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

	var result StdResponse
	if err := json.Unmarshal(bytesResp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

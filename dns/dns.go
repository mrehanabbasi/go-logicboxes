// Package dns contains APIs for dns.
package dns

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mrehanabbasi/go-logicboxes/core"
)

type DNS interface {
	ActivatingDNSService(ctx context.Context, orderID string) (*ActivatingDNSServiceResponse, error)
	AddingIPv4AddressRecord(ctx context.Context, domainName, value, host string, ttl int) (*StdResponse, error)
	AddingIPv6AddressRecord(ctx context.Context, domainName, value, host string, ttl int) (*StdResponse, error)
	AddingCNAMERecord(ctx context.Context, domainName, value, host string, ttl int) (*StdResponse, error)
	AddingMXRecord(ctx context.Context, domainName, value, host string, ttl, priority int) (*StdResponse, error)
	AddingNSRecord(ctx context.Context, domainName, value, host string, ttl int) (*StdResponse, error)
	AddingTXTRecord(ctx context.Context, domainName, value, host string, ttl int) (*StdResponse, error)
	AddingSRVRecord(ctx context.Context, domainName, value, host string, ttl, priority, port, weight int) (*StdResponse, error)
	ModifyingIPv4AddressRecord(ctx context.Context, domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error)
	ModifyingIPv6AddressRecord(ctx context.Context, domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error)
	ModifyingCNAMERecord(ctx context.Context, domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error)
	ModifyingMXRecord(ctx context.Context, domainName, host, currentValue, newValue string, ttl, priority int) (*StdResponse, error)
	ModifyingNSRecord(ctx context.Context, domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error)
	ModifyingTXTRecord(ctx context.Context, domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error)
	ModifyingSRVRecord(
		ctx context.Context,
		domainName, host, currentValue, newValue string,
		ttl, priority, port, weight int,
	) (*StdResponse, error)
	ModifyingSOARecord(ctx context.Context, domainName, responsiblePerson string, refresh, retry, expire, ttl int) (*StdResponse, error)
	SearchingDNSRecords(
		ctx context.Context,
		domainName string,
		typeRecord RecordType,
		noOfRecords, pageNo int,
		host, value string,
	) (*SearchingDNSRecords, error)
	DeletingDNSRecord(ctx context.Context, host, value string) (*StdResponse, error)
	DeletingIPv4AddressRecord(ctx context.Context, domainName, host, value string) (*StdResponse, error)
	DeletingIPv6AddressRecord(ctx context.Context, domainName, host, value string) (*StdResponse, error)
	DeletingCNAMERecord(ctx context.Context, domainName, host, value string) (*StdResponse, error)
	DeletingMXRecord(ctx context.Context, domainName, host, value string) (*StdResponse, error)
	DeletingNSRecord(ctx context.Context, domainName, host, value string) (*StdResponse, error)
	DeletingTXTRecord(ctx context.Context, domainName, host, value string) (*StdResponse, error)
	DeletingSRVRecord(ctx context.Context, domainName, host, value string, port, weight int) (*StdResponse, error)
}

func New(c core.Core) DNS {
	return &dns{c}
}

type dns struct {
	core core.Core
}

func (d *dns) ActivatingDNSService(ctx context.Context, orderID string) (*ActivatingDNSServiceResponse, error) {
	data := make(url.Values)
	data.Add("order-id", orderID)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "activate", data)
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

func (d *dns) AddingIPv4AddressRecord(ctx context.Context, domainName, value, host string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/add-ipv4-record", data)
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

func (d *dns) AddingIPv6AddressRecord(ctx context.Context, domainName, value, host string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/add-ipv6-record", data)
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

func (d *dns) AddingCNAMERecord(ctx context.Context, domainName, value, host string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/add-cname-record", data)
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

func (d *dns) AddingMXRecord(ctx context.Context, domainName, value, host string, ttl, priority int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))
	data.Add("priority", strconv.Itoa(priority))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/add-mx-record", data)
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

func (d *dns) AddingNSRecord(ctx context.Context, domainName, value, host string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/add-ns-record", data)
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

func (d *dns) AddingTXTRecord(ctx context.Context, domainName, value, host string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/manage/add-ns-record", data)
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

func (d *dns) AddingSRVRecord(ctx context.Context, domainName, value, host string, ttl, priority, port, weight int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("value", value)
	data.Add("host", host)
	data.Add("ttl", strconv.Itoa(ttl))
	data.Add("priority", strconv.Itoa(priority))
	data.Add("port", strconv.Itoa(port))
	data.Add("weight", strconv.Itoa(weight))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/add-srv-record", data)
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

func (d *dns) ModifyingIPv4AddressRecord(
	ctx context.Context,
	domainName, host, currentValue, newValue string,
	ttl int,
) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/update-ipv4-record", data)
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

func (d *dns) ModifyingIPv6AddressRecord(
	ctx context.Context,
	domainName, host, currentValue, newValue string,
	ttl int,
) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/update-ipv6-record", data)
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

func (d *dns) ModifyingCNAMERecord(ctx context.Context, domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/update-cname-record", data)
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

func (d *dns) ModifyingMXRecord(
	ctx context.Context,
	domainName, host, currentValue, newValue string,
	ttl, priority int,
) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))
	data.Add("priority", strconv.Itoa(priority))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/update-mx-record", data)
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

func (d *dns) ModifyingNSRecord(ctx context.Context, domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/update-ns-record", data)
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

func (d *dns) ModifyingTXTRecord(ctx context.Context, domainName, host, currentValue, newValue string, ttl int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/update-txt-record", data)
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

func (d *dns) ModifyingSRVRecord(
	ctx context.Context,
	domainName, host, currentValue, newValue string,
	ttl, priority, port, weight int,
) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("current-value", currentValue)
	data.Add("new-value", newValue)
	data.Add("ttl", strconv.Itoa(ttl))
	data.Add("priority", strconv.Itoa(priority))
	data.Add("port", strconv.Itoa(port))
	data.Add("weight", strconv.Itoa(weight))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/update-srv-record", data)
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

func (d *dns) ModifyingSOARecord(
	ctx context.Context,
	domainName, responsiblePerson string,
	refresh, retry, expire, ttl int,
) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("responsible-person", responsiblePerson)
	data.Add("refresh", strconv.Itoa(refresh))
	data.Add("retry", strconv.Itoa(retry))
	data.Add("expire", strconv.Itoa(expire))
	data.Add("ttl", strconv.Itoa(ttl))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/update-soa-record", data)
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
	ctx context.Context,
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

	resp, err := d.core.CallAPI(ctx, http.MethodGet, "dns", "manage/search-records", data)
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

func (d *dns) DeletingDNSRecord(ctx context.Context, host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(ctx, http.MethodGet, "dns", "manage/delete-record", data)
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

func (d *dns) DeletingIPv4AddressRecord(ctx context.Context, domainName, host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/delete-ipv4-record", data)
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

func (d *dns) DeletingIPv6AddressRecord(ctx context.Context, domainName, host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/delete-ipv6-record", data)
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

func (d *dns) DeletingCNAMERecord(ctx context.Context, domainName, host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/delete-cname-record", data)
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

func (d *dns) DeletingMXRecord(ctx context.Context, domainName, host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/delete-mx-record", data)
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

func (d *dns) DeletingNSRecord(ctx context.Context, domainName, host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/delete-ns-record", data)
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

func (d *dns) DeletingTXTRecord(ctx context.Context, domainName, host, value string) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/delete-txt-record", data)
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

func (d *dns) DeletingSRVRecord(ctx context.Context, domainName, host, value string, port, weight int) (*StdResponse, error) {
	data := make(url.Values)
	data.Add("domain-name", domainName)
	data.Add("host", host)
	data.Add("value", value)
	data.Add("port", strconv.Itoa(port))
	data.Add("weight", strconv.Itoa(weight))

	resp, err := d.core.CallAPI(ctx, http.MethodPost, "dns", "manage/delete-srv-record", data)
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

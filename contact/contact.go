// Package contact contains APIs for contacts-related tasks.
package contact

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/mrehanabbasi/go-logicboxes/core"
)

type contact struct {
	core core.Core
}

type Contact interface {
	Add(ctx context.Context, details *Detail, attributes core.EntityAttributes) error
	Details(ctx context.Context, contactID string) (*Detail, error)
	Delete(ctx context.Context, contactID string) (*Action, error)
	Search(ctx context.Context, criteria Criteria, offset, limit uint16) (*SearchResult, error)
	SetDefault(
		ctx context.Context,
		customerID, registrantContactID, adminContactID, techContactID, billingContactID string,
		types []Type,
	) error
	Default(ctx context.Context, customerID string, types []Type) (map[string]Detail, error)
	ValidateRegistrant(ctx context.Context, contactID string, eligibilities []Eligibility) (RegistrantValidation, error)
	AddExtraDetails(ctx context.Context, contactID string, attributes core.EntityAttributes, domainKeys []core.DomainKey) error
	DotCAAgreement(ctx context.Context) (map[string]string, error)
	// AddDotCOOPSponsor(ctx context.Context, customerId string, details ContactDetail) (string, error)
	// DotCOOPSponsors(ctx context.Context, customerId string) error
}

func (c *contact) DotCAAgreement(ctx context.Context) (map[string]string, error) {
	resp, err := c.core.CallAPI(ctx, http.MethodGet, "contacts/dotca", "registrantagreement", url.Values{})
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

	ret := map[string]string{}
	if err := json.Unmarshal(bytesResp, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// func (c *contact) DotCOOPSponsors(ctx context.Context, customerId string) error {
// 	if !core.RgxNumber.MatchString(customerId) {
// 		return core.ErrRcInvalidCredential
// 	}

// 	resp, err := c.core.CallAPI(ctx, http.MethodGet, "contacts", "sponsors", url.Values{"customer-id": []string{customerId}})
// 	if err != nil {
// 		return err
// 	}
// 	defer func() { _ = resp.Body.Close() }()

// 	bytesResp, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		errResponse := core.JSONStatusResponse{}
// 		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
// 			return err
// 		}
// 		return errors.New(strings.ToLower(errResponse.Message))
// 	}

// 	return nil
// }

// func (c *contact) AddDotCOOPSponsor(ctx context.Context, customerId string, details ContactDetail) (string, error) {
// 	if !core.RgxNumber.MatchString(customerId) {
// 		return "", core.ErrRcInvalidCredential
// 	}

// 	if !core.RgxEmail.MatchString(details.Email) {
// 		return "", errors.New("invalid format for email")
// 	}

// 	data, err := extractSponsorData(details)
// 	if err != nil {
// 		return "", err
// 	}
// 	data.Add("customer-id", customerId)

// 	resp, err := c.core.CallAPI(ctx, http.MethodPost, "contacts/coop", "add-sponsor", *data)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer func() { _ = resp.Body.Close() }()

// 	bytesResp, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		errResponse := core.JSONStatusResponse{}
// 		if err := json.Unmarshal(bytesResp, &errResponse); err != nil {
// 			return "", err
// 		}
// 		return "", errors.New(strings.ToLower(errResponse.Message))
// 	}

// 	return string(bytesResp), nil
// }

func (c *contact) AddExtraDetails(
	ctx context.Context,
	contactID string,
	attributes core.EntityAttributes,
	domainKeys []core.DomainKey,
) error {
	if !core.RgxNumber.MatchString(contactID) {
		return core.ErrRcInvalidCredential
	}

	if attributes == nil || domainKeys == nil || len(domainKeys) == 0 {
		return errors.New("attributes and domain keys cannot be nil or empty")
	}

	data := url.Values{}
	data.Add("contact-id", contactID)
	attributes.CopyTo(&data)

	wg := sync.WaitGroup{}
	rwMutex := sync.RWMutex{}
	for _, k := range domainKeys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			rwMutex.Lock()
			data.Add("product-key", key)
			rwMutex.Unlock()
		}(string(k))
	}
	wg.Wait()

	resp, err := c.core.CallAPI(ctx, http.MethodPost, "contacts", "set-details", data)
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

	boolResult, err := strconv.ParseBool(string(bytesResp))
	if err != nil {
		return err
	}
	if !boolResult {
		return core.ErrRcOperationFailed
	}

	return nil
}

func (c *contact) ValidateRegistrant(ctx context.Context, contactID string, eligibilities []Eligibility) (RegistrantValidation, error) {
	if !core.RgxNumber.MatchString(contactID) {
		return nil, core.ErrRcInvalidCredential
	}

	if len(eligibilities) == 0 {
		return nil, errors.New("eligibilities must not empty")
	}

	data := url.Values{}
	data.Add("contact-id", contactID)

	wg := sync.WaitGroup{}
	rwMutex := sync.RWMutex{}
	for _, eligibility := range eligibilities {
		wg.Add(1)
		go func(e Eligibility) {
			defer wg.Done()
			rwMutex.Lock()
			data.Add("eligibility-criteria", string(e))
			rwMutex.Unlock()
		}(eligibility)
	}
	wg.Wait()

	resp, err := c.core.CallAPI(ctx, http.MethodGet, "contacts", "validate-registrant", data)
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

	validation := RegistrantValidation{}
	if err := json.Unmarshal(bytesResp, &validation); err != nil {
		return nil, err
	}

	return validation, nil
}

//nolint:funlen
func (c *contact) Default(ctx context.Context, customerID string, types []Type) (map[string]Detail, error) {
	if len(types) == 0 {
		return nil, errors.New("contact types must not empty")
	}
	if !core.RgxNumber.MatchString(customerID) {
		return nil, core.ErrRcInvalidCredential
	}

	data := url.Values{}
	data.Add("customer-id", customerID)
	for _, t := range types {
		data.Add("type", string(t))
	}

	resp, err := c.core.CallAPI(ctx, http.MethodPost, "contacts", "default", data)
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

	replacer := strings.NewReplacer("contact.", "", "entity.", "")
	strResp := replacer.Replace(string(bytesResp))
	bytesResp = []byte(strResp)

	exoSkeleton := map[string]core.JSONBytes{}
	if err := json.Unmarshal(bytesResp, &exoSkeleton); err != nil {
		return nil, err
	}
	if len(exoSkeleton) == 0 {
		return nil, errors.New("failed while extract exoskeleton")
	}

	contacts := map[string]core.JSONBytes{}
	for _, elem := range exoSkeleton {
		bytesResp = []byte(elem)
		if err := json.Unmarshal(bytesResp, &contacts); err != nil {
			return nil, err
		}
	}

	wg := sync.WaitGroup{}
	rwMutex := sync.RWMutex{}
	defaultContacts := map[string]Detail{}

	for k, v := range contacts {
		wg.Add(1)
		go func(key string, val core.JSONBytes) {
			defer wg.Done()
			bytesValue := []byte(val)
			switch key {
			case "registrant", "type", "tech", "billing", "admin":
				return
			default:
				ctc := Detail{}
				if err := json.Unmarshal(bytesValue, &ctc); err != nil {
					return
				}
				rwMutex.Lock()
				defaultContacts[strings.TrimSuffix(key, "ContactDetails")] = ctc
				rwMutex.Unlock()
			}
		}(k, v)
	}
	wg.Wait()

	return defaultContacts, nil
}

func (c *contact) SetDefault(
	ctx context.Context,
	customerID, regContactID, adminContactID, techContactID, billContactID string,
	types []Type,
) error {
	if len(types) == 0 {
		return errors.New("contact types must not empty")
	}
	if !core.RgxNumber.MatchString(customerID) || !core.RgxNumber.MatchString(regContactID) ||
		!core.RgxNumber.MatchString(adminContactID) || !core.RgxNumber.MatchString(techContactID) ||
		!core.RgxNumber.MatchString(billContactID) {
		return core.ErrRcInvalidCredential
	}

	data := url.Values{}
	data.Add("customer-id", customerID)
	data.Add("reg-contact-id", regContactID)
	data.Add("admin-contact-id", adminContactID)
	data.Add("tech-contact-id", techContactID)
	data.Add("billing-contact-id", billContactID)

	for _, t := range types {
		data.Add("type", string(t))
	}

	resp, err := c.core.CallAPI(ctx, http.MethodPost, "contacts", "modDefault", data)
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

func (c *contact) Search(ctx context.Context, criteria Criteria, offset, limit uint16) (*SearchResult, error) {
	if offset <= 0 || limit <= 0 {
		return nil, errors.New("offset or limit must greater than zero")
	}

	if err := validator.New().Struct(criteria); err != nil {
		return nil, err
	}

	data, err := criteria.URLValues()
	if err != nil {
		return nil, err
	}
	data.Add("no-of-records", strconv.FormatUint(uint64(limit), 10))
	data.Add("page-no", strconv.FormatUint(uint64(offset), 10))

	resp, err := c.core.CallAPI(ctx, http.MethodGet, "contacts", "search", data)
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

	replacer := strings.NewReplacer("entity.", "", "contact.", "")
	strResp := replacer.Replace(string(bytesResp))

	var buffer map[string]core.JSONBytes
	if err := json.Unmarshal([]byte(strResp), &buffer); err != nil {
		return nil, err
	}

	var dataBuffers []Detail
	var numMatched int

	for key, dataBytes := range buffer {
		switch {
		case key == "recsindb":
			numMatched, err = strconv.Atoi(string(dataBytes))
			if err != nil {
				numMatched = 0
			}
		case key == "result":
			if err := json.Unmarshal(dataBytes, &dataBuffers); err != nil {
				return nil, err
			}
		}
	}

	return &SearchResult{
		RequestedLimit:  limit,
		RequestedOffset: offset,
		Contacts:        dataBuffers,
		TotalMatched:    numMatched,
	}, nil
}

func (c *contact) Delete(ctx context.Context, contactID string) (*Action, error) {
	if !core.RgxNumber.MatchString(contactID) {
		return nil, core.ErrRcInvalidCredential
	}

	resp, err := c.core.CallAPI(ctx, http.MethodPost, "contacts", "delete", url.Values{"contact-id": {contactID}})
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

	ret := new(Action)
	if err := json.Unmarshal(bytesResp, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *contact) Details(ctx context.Context, contactID string) (*Detail, error) {
	if !core.RgxNumber.MatchString(contactID) {
		return nil, core.ErrRcInvalidCredential
	}

	data := url.Values{}
	data.Add("contact-id", contactID)

	resp, err := c.core.CallAPI(ctx, http.MethodGet, "contacts", "details", data)
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

	ret := new(Detail)
	if err := json.Unmarshal(bytesResp, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *contact) Add(ctx context.Context, details *Detail, attributes core.EntityAttributes) error {
	if details == nil {
		return errors.New("detail must not nil")
	}

	data, err := details.URLValues()
	if err != nil {
		return err
	}

	if attributes != nil {
		attributes.CopyTo(data)
	}

	resp, err := c.core.CallAPI(ctx, http.MethodPost, "contacts", "add", *data)
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

	details.ID = string(bytesResp)
	return nil
}

func New(c core.Core) Contact {
	return &contact{
		core: c,
	}
}

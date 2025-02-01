package contact

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/mrehanabbasi/go-resellerclub/core"
)

type contact struct {
	core core.Core
}

type Contact interface {
	Add(details *ContactDetail, attributes core.EntityAttributes) error
	Details(contactId string) (*ContactDetail, error)
	Delete(contactId string) (*Action, error)
	Search(criteria ContactCriteria, offset, limit uint16) (*ContactSearchResult, error)
	SetDefault(customerId, registrantContactID, adminContactID, techContactID, billingContactID string, types []ContactType) error
	Default(customerId string, types []ContactType) (map[string]ContactDetail, error)
	ValidateRegistrant(contactId string, eligibilities []Eligibility) (RegistrantValidation, error)
	AddExtraDetails(contactId string, attributes core.EntityAttributes, domainKeys []core.DomainKey) error
	DotCAAgreement() (map[string]string, error)
	// AddDotCOOPSponsor(customerId string, details ContactDetail) (string, error)
	// DotCOOPSponsors(customerId string) error
}

func (c *contact) DotCAAgreement() (map[string]string, error) {
	resp, err := c.core.CallApi(http.MethodGet, "contacts/dotca", "registrantagreement", url.Values{})
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

	ret := map[string]string{}
	if err = json.Unmarshal(bytesResp, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// func (c *contact) DotCOOPSponsors(customerId string) error {
// 	if !core.RgxNumber.MatchString(customerId) {
// 		return core.ErrRcInvalidCredential
// 	}

// 	resp, err := c.core.CallApi(http.MethodGet, "contacts", "sponsors", url.Values{"customer-id": []string{customerId}})
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	bytesResp, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		errResponse := core.JSONStatusResponse{}
// 		err = json.Unmarshal(bytesResp, &errResponse)
// 		if err != nil {
// 			return err
// 		}
// 		return errors.New(strings.ToLower(errResponse.Message))
// 	}

// 	return nil
// }

// func (c *contact) AddDotCOOPSponsor(customerId string, details ContactDetail) (string, error) {
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

// 	resp, err := c.core.CallApi(http.MethodPost, "contacts/coop", "add-sponsor", *data)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	bytesResp, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		errResponse := core.JSONStatusResponse{}
// 		err = json.Unmarshal(bytesResp, &errResponse)
// 		if err != nil {
// 			return "", err
// 		}
// 		return "", errors.New(strings.ToLower(errResponse.Message))
// 	}

// 	return string(bytesResp), nil
// }

func (c *contact) AddExtraDetails(contactId string, attributes core.EntityAttributes, domainKeys []core.DomainKey) error {
	if !core.RgxNumber.MatchString(contactId) {
		return core.ErrRcInvalidCredential
	}

	if attributes == nil || domainKeys == nil || len(domainKeys) <= 0 {
		return errors.New("attributes and domain keys cannot be nil or empty")
	}

	data := url.Values{}
	data.Add("contact-id", contactId)
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

	resp, err := c.core.CallApi(http.MethodPost, "contacts", "set-details", data)
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

	boolResult, err := strconv.ParseBool(string(bytesResp))
	if err != nil {
		return err
	}
	if !boolResult {
		return core.ErrRcOperationFailed
	}

	return nil
}

func (c *contact) ValidateRegistrant(contactId string, eligibilities []Eligibility) (RegistrantValidation, error) {
	if !core.RgxNumber.MatchString(contactId) {
		return nil, core.ErrRcInvalidCredential
	}

	if len(eligibilities) <= 0 {
		return nil, errors.New("eligibilities must not empty")
	}

	data := url.Values{}
	data.Add("contact-id", contactId)

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

	resp, err := c.core.CallApi(http.MethodGet, "contacts", "validate-registrant", data)
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

	validation := RegistrantValidation{}
	err = json.Unmarshal(bytesResp, &validation)
	if err != nil {
		return nil, err
	}

	return validation, nil
}

func (c *contact) Default(customerId string, types []ContactType) (map[string]ContactDetail, error) {
	if len(types) <= 0 {
		return nil, errors.New("contact types must not empty")
	}
	if !core.RgxNumber.MatchString(customerId) {
		return nil, core.ErrRcInvalidCredential
	}

	data := url.Values{}
	data.Add("customer-id", customerId)
	for _, t := range types {
		data.Add("type", string(t))
	}

	resp, err := c.core.CallApi(http.MethodPost, "contacts", "default", data)
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

	replacer := strings.NewReplacer("contact.", "", "entity.", "")
	strResp := replacer.Replace(string(bytesResp))
	bytesResp = []byte(strResp)

	exoSkeleton := map[string]core.JSONBytes{}
	err = json.Unmarshal(bytesResp, &exoSkeleton)
	if err != nil {
		return nil, err
	}
	if len(exoSkeleton) <= 0 {
		return nil, errors.New("failed while extract exoskeleton")
	}

	contacts := map[string]core.JSONBytes{}
	for _, elem := range exoSkeleton {
		bytesResp = []byte(elem)
		if err = json.Unmarshal(bytesResp, &contacts); err != nil {
			return nil, err
		}
	}

	wg := sync.WaitGroup{}
	rwMutex := sync.RWMutex{}
	defaultContacts := map[string]ContactDetail{}

	for k, v := range contacts {
		wg.Add(1)
		go func(key string, val core.JSONBytes) {
			defer wg.Done()
			bytesValue := []byte(val)
			switch key {
			case "registrant", "type", "tech", "billing", "admin":
				return
			default:
				ctc := ContactDetail{}
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

func (c *contact) SetDefault(customerId, regContactID, adminContactID, techContactID, billContactID string, types []ContactType) error {
	if len(types) <= 0 {
		return errors.New("contact types must not empty")
	}
	if !core.RgxNumber.MatchString(customerId) || !core.RgxNumber.MatchString(regContactID) || !core.RgxNumber.MatchString(adminContactID) || !core.RgxNumber.MatchString(techContactID) || !core.RgxNumber.MatchString(billContactID) {
		return core.ErrRcInvalidCredential
	}

	data := url.Values{}
	data.Add("customer-id", customerId)
	data.Add("reg-contact-id", regContactID)
	data.Add("admin-contact-id", adminContactID)
	data.Add("tech-contact-id", techContactID)
	data.Add("billing-contact-id", billContactID)

	for _, t := range types {
		data.Add("type", string(t))
	}

	resp, err := c.core.CallApi(http.MethodPost, "contacts", "modDefault", data)
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

func (c *contact) Search(criteria ContactCriteria, offset, limit uint16) (*ContactSearchResult, error) {
	if offset <= 0 || limit <= 0 {
		return nil, errors.New("offset or limit must greater than zero")
	}

	if err := validator.New().Struct(criteria); err != nil {
		return nil, err
	}

	data, err := criteria.UrlValues()
	if err != nil {
		return nil, err
	}
	data.Add("no-of-records", strconv.FormatUint(uint64(limit), 10))
	data.Add("page-no", strconv.FormatUint(uint64(offset), 10))

	resp, err := c.core.CallApi(http.MethodGet, "contacts", "search", data)
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

	replacer := strings.NewReplacer("entity.", "", "contact.", "")
	strResp := replacer.Replace(string(bytesResp))

	var buffer map[string]core.JSONBytes
	if err := json.Unmarshal([]byte(strResp), &buffer); err != nil {
		return nil, err
	}

	var dataBuffers []ContactDetail
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

	return &ContactSearchResult{
		RequestedLimit:  limit,
		RequestedOffset: offset,
		Contacts:        dataBuffers,
		TotalMatched:    numMatched,
	}, nil
}

func (c *contact) Delete(contactId string) (*Action, error) {
	if !core.RgxNumber.MatchString(contactId) {
		return nil, core.ErrRcInvalidCredential
	}

	resp, err := c.core.CallApi(http.MethodPost, "contacts", "delete", url.Values{"contact-id": {contactId}})
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

	ret := new(Action)
	if err = json.Unmarshal(bytesResp, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *contact) Details(contactId string) (*ContactDetail, error) {
	if !core.RgxNumber.MatchString(contactId) {
		return nil, core.ErrRcInvalidCredential
	}

	data := url.Values{}
	data.Add("contact-id", contactId)

	resp, err := c.core.CallApi(http.MethodGet, "contacts", "details", data)
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

	ret := new(ContactDetail)
	if err = json.Unmarshal(bytesResp, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *contact) Add(details *ContactDetail, attributes core.EntityAttributes) error {
	if details == nil {
		return errors.New("detail must not nil")
	}

	data, err := details.UrlValues()
	if err != nil {
		return err
	}

	if attributes != nil {
		attributes.CopyTo(data)
	}

	resp, err := c.core.CallApi(http.MethodPost, "contacts", "add", *data)
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

	details.Id = string(bytesResp)
	return nil
}

func New(c core.Core) Contact {
	return &contact{
		core: c,
	}
}

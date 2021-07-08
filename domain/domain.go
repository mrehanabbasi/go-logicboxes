package domain

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/xpartacvs/go-resellerclub/core"
)

type domain struct {
	core core.Core
}

type Domain interface {
	SearchRegistrationOrders(criteria OrdersCriteria) error
	CheckAvailability(domainsWithoutTLD, tlds []string) (DomainAvailabilities, error)
}

func New(c core.Core) Domain {
	return &domain{
		core: c,
	}
}

func (d *domain) CheckAvailability(domainsWithoutTLD, tlds []string) (DomainAvailabilities, error) {
	if len(domainsWithoutTLD) <= 0 || len(tlds) <= 0 {
		return DomainAvailabilities{}, errors.New("domainnames and tlds must not empty")
	}

	data := url.Values{}
	wg := sync.WaitGroup{}
	rwMutex := sync.RWMutex{}

	for _, v := range domainsWithoutTLD {
		wg.Add(1)
		go func(value string) {
			defer wg.Done()
			defer rwMutex.Unlock()
			rwMutex.Lock()
			data.Add("domain-name", value)
		}(v)
	}
	for _, v := range tlds {
		wg.Add(1)
		go func(value string) {
			defer wg.Done()
			defer rwMutex.Unlock()
			rwMutex.Lock()
			data.Add("tlds", value)
		}(v)
	}
	wg.Wait()

	resp, err := d.core.CallApi(http.MethodGet, "domains", "available", data)
	if err != nil {
		return DomainAvailabilities{}, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return DomainAvailabilities{}, err
	}

	if resp.StatusCode != http.StatusOK {
		errResponse := core.JSONStatusResponse{}
		err = json.Unmarshal(bytesResp, &errResponse)
		if err != nil {
			return DomainAvailabilities{}, err
		}
		return DomainAvailabilities{}, errors.New(strings.ToLower(errResponse.Message))
	}

	availabilities := DomainAvailabilities{}
	err = json.Unmarshal(bytesResp, &availabilities)
	if err != nil {
		return DomainAvailabilities{}, err
	}

	return availabilities, nil
}

func (d *domain) SearchRegistrationOrders(criteria OrdersCriteria) error {
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

	d.core.PrintResponse(bytesResp)

	return nil
}

package general

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/mrehanabbasi/go-resellerclub/core"
)

type stateMap map[string]string

type States interface {
	Name(id string) string
	Length() int
	ToMap() map[string]string
}

func fetchStateList(c core.Core, cc CountryISO) (States, error) {
	data := url.Values{}
	data.Add("country-code", string(cc))

	resp, err := c.CallAPI(http.MethodGet, "country", "state-list", data)
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

	keyPairs := map[string]string{}
	if err := json.Unmarshal(bytesResp, &keyPairs); err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	rwMutex := sync.RWMutex{}
	mapStates := stateMap{}
	for key, val := range keyPairs {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			rwMutex.Lock()
			mapStates[v] = k
			rwMutex.Unlock()
		}(key, val)
	}
	wg.Wait()

	return mapStates, nil
}

func (s stateMap) Name(id string) string {
	return s[id]
}

func (s stateMap) ToMap() map[string]string {
	return map[string]string(s)
}

func (s stateMap) Length() int {
	return len(s.ToMap())
}

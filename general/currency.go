package general

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/mrehanabbasi/go-logicboxes/core"
)

type (
	CurrencyISO string
	currencyDB  map[CurrencyISO]Currency
)

type Currency struct {
	ISO  CurrencyISO
	Unit uint16
	Name string
}

// Const for currency ISO values.
const (
	IsoTWD CurrencyISO = "TWD"
	IsoKES CurrencyISO = "KES"
	IsoLKR CurrencyISO = "LKR"
	IsoRWF CurrencyISO = "RWF"
	IsoAFN CurrencyISO = "AFN"
	IsoSDG CurrencyISO = "SDG"
	IsoARS CurrencyISO = "ARS"
	IsoGEL CurrencyISO = "GEL"
	IsoCRC CurrencyISO = "CRC"
	IsoFKP CurrencyISO = "FKP"
	IsoEEK CurrencyISO = "EEK"
	IsoHKD CurrencyISO = "HKD"
	IsoMDL CurrencyISO = "MDL"
	IsoIQD CurrencyISO = "IQD"
	IsoSCR CurrencyISO = "SCR"
	IsoVUV CurrencyISO = "VUV"
	IsoDKK CurrencyISO = "DKK"
	IsoIDR CurrencyISO = "IDR"
	IsoSOS CurrencyISO = "SOS"
	IsoAED CurrencyISO = "AED"
	IsoBWP CurrencyISO = "BWP"
	IsoLVL CurrencyISO = "LVL"
	IsoNIO CurrencyISO = "NIO"
	IsoADP CurrencyISO = "ADP"
	IsoFJD CurrencyISO = "FJD"
	IsoMOP CurrencyISO = "MOP"
	IsoRUB CurrencyISO = "RUB"
	IsoCDF CurrencyISO = "CDF"
	IsoDJF CurrencyISO = "DJF"
	IsoADF CurrencyISO = "ADF"
	IsoSBD CurrencyISO = "SBD"
	IsoECS CurrencyISO = "ECS"
	IsoPHP CurrencyISO = "PHP"
	IsoTHB CurrencyISO = "THB"
	IsoTTD CurrencyISO = "TTD"
	IsoSZL CurrencyISO = "SZL"
	IsoMNT CurrencyISO = "MNT"
	IsoSAR CurrencyISO = "SAR"
	IsoUAH CurrencyISO = "UAH"
	IsoHUF CurrencyISO = "HUF"
	IsoCOP CurrencyISO = "COP"
	IsoQAR CurrencyISO = "QAR"
	IsoUYU CurrencyISO = "UYU"
	IsoGBP CurrencyISO = "GBP"
	IsoBIF CurrencyISO = "BIF"
	IsoINR CurrencyISO = "INR"
	IsoLTL CurrencyISO = "LTL"
	IsoMZN CurrencyISO = "MZN"
	IsoKZT CurrencyISO = "KZT"
	IsoPGK CurrencyISO = "PGK"
	IsoOMR CurrencyISO = "OMR"
	IsoNGN CurrencyISO = "NGN"
	IsoAOA CurrencyISO = "AOA"
	IsoCNY CurrencyISO = "CNY"
	IsoGNF CurrencyISO = "GNF"
	IsoRSD CurrencyISO = "RSD"
	IsoHTG CurrencyISO = "HTG"
	IsoMAD CurrencyISO = "MAD"
	IsoTRY CurrencyISO = "TRY"
	IsoMMK CurrencyISO = "MMK"
	IsoMYR CurrencyISO = "MYR"
	IsoLSL CurrencyISO = "LSL"
	IsoBHD CurrencyISO = "BHD"
	IsoSLL CurrencyISO = "SLL"
	IsoBTN CurrencyISO = "BTN"
	IsoKMF CurrencyISO = "KMF"
	IsoANG CurrencyISO = "ANG"
	IsoCZK CurrencyISO = "CZK"
	IsoAZN CurrencyISO = "AZN"
	IsoKYD CurrencyISO = "KYD"
	IsoGMD CurrencyISO = "GMD"
	IsoVEF CurrencyISO = "VEF"
	IsoBGN CurrencyISO = "BGN"
	IsoCAD CurrencyISO = "CAD"
	IsoILS CurrencyISO = "ILS"
	IsoGYD CurrencyISO = "GYD"
	IsoMXN CurrencyISO = "MXN"
	IsoPEN CurrencyISO = "PEN"
	IsoLRD CurrencyISO = "LRD"
	IsoAMD CurrencyISO = "AMD"
	IsoBSD CurrencyISO = "BSD"
	IsoHRK CurrencyISO = "HRK"
	IsoCLP CurrencyISO = "CLP"
	IsoMKD CurrencyISO = "MKD"
	IsoALL CurrencyISO = "ALL"
	IsoMWK CurrencyISO = "MWK"
	IsoBRL CurrencyISO = "BRL"
	IsoKWD CurrencyISO = "KWD"
	IsoXCD CurrencyISO = "XCD"
	IsoNPR CurrencyISO = "NPR"
	IsoSVC CurrencyISO = "SVC"
	IsoJPY CurrencyISO = "JPY"
	IsoXOF CurrencyISO = "XOF"
	IsoMVR CurrencyISO = "MVR"
	IsoTOP CurrencyISO = "TOP"
	IsoRON CurrencyISO = "RON"
	IsoBDT CurrencyISO = "BDT"
	IsoAWG CurrencyISO = "AWG"
	IsoNOK CurrencyISO = "NOK"
	IsoMUR CurrencyISO = "MUR"
	IsoZAR CurrencyISO = "ZAR"
	IsoSHP CurrencyISO = "SHP"
	IsoZMW CurrencyISO = "ZMW"
	IsoVND CurrencyISO = "VND"
	IsoTZS CurrencyISO = "TZS"
	IsoGIP CurrencyISO = "GIP"
	IsoTND CurrencyISO = "TND"
	IsoUGX CurrencyISO = "UGX"
	IsoCVE CurrencyISO = "CVE"
	IsoJOD CurrencyISO = "JOD"
	IsoXAF CurrencyISO = "XAF"
	IsoLBP CurrencyISO = "LBP"
	IsoUGS CurrencyISO = "UGS"
	IsoSTD CurrencyISO = "STD"
	IsoWST CurrencyISO = "WST"
	IsoKHR CurrencyISO = "KHR"
	IsoDOP CurrencyISO = "DOP"
	IsoEUR CurrencyISO = "EUR"
	IsoTMT CurrencyISO = "TMT"
	IsoGHS CurrencyISO = "GHS"
	IsoSGD CurrencyISO = "SGD"
	IsoNZD CurrencyISO = "NZD"
	IsoUSD CurrencyISO = "USD"
	IsoBOB CurrencyISO = "BOB"
	IsoHNL CurrencyISO = "HNL"
	IsoPAB CurrencyISO = "PAB"
	IsoGTQ CurrencyISO = "GTQ"
	IsoAUD CurrencyISO = "AUD"
	IsoLAK CurrencyISO = "LAK"
	IsoNAD CurrencyISO = "NAD"
	IsoKGS CurrencyISO = "KGS"
	IsoBBD CurrencyISO = "BBD"
	IsoCHF CurrencyISO = "CHF"
	IsoMGA CurrencyISO = "MGA"
	IsoPYG CurrencyISO = "PYG"
	IsoYER CurrencyISO = "YER"
	IsoETB CurrencyISO = "ETB"
	IsoBND CurrencyISO = "BND"
	IsoEGP CurrencyISO = "EGP"
	IsoJMD CurrencyISO = "JMD"
	IsoPLN CurrencyISO = "PLN"
	IsoDZD CurrencyISO = "DZD"
	IsoISK CurrencyISO = "ISK"
	IsoLYD CurrencyISO = "LYD"
	IsoSRD CurrencyISO = "SRD"
	IsoBAM CurrencyISO = "BAM"
	IsoBZD CurrencyISO = "BZD"
	IsoKRW CurrencyISO = "KRW"
	IsoMRO CurrencyISO = "MRO"
	IsoZWD CurrencyISO = "ZWD"
	IsoSEK CurrencyISO = "SEK"
	IsoCSK CurrencyISO = "CSK"
	IsoBYR CurrencyISO = "BYR"
	IsoPKR CurrencyISO = "PKR"
	IsoBMD CurrencyISO = "BMD"
)

func fetchCurrencyDB(c core.Core) (currencyDB, error) {
	resp, err := c.CallAPI(http.MethodGet, "currency", "details", url.Values{})
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

	ret := make(map[string]map[string]string)
	if err := json.Unmarshal(bytesResp, &ret); err != nil {
		return nil, err
	}

	cdb := currencyDB{}
	rwMutex := sync.RWMutex{}
	wg := sync.WaitGroup{}
	for iso, cur := range ret {
		wg.Add(1)
		go func(key string, val map[string]string) {
			defer wg.Done()
			unit, err := strconv.Atoi(val["currencyunit"])
			if err != nil {
				return
			}
			iso := CurrencyISO(key)
			rwMutex.Lock()
			cdb[iso] = Currency{
				ISO:  iso,
				Unit: uint16(unit),
				Name: val["currencyname"],
			}
			rwMutex.Unlock()
		}(iso, cur)
	}
	wg.Wait()

	return cdb, nil
}

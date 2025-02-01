// Package general contains general util functions and variables.
package general

import (
	"github.com/mrehanabbasi/go-resellerclub/core"
)

type general struct {
	core       core.Core
	currencies currencyDB
	countries  countryDB
}

type General interface {
	CurrencyOf(iso CurrencyISO) Currency
	CountryName(iso CountryISO) string
	StatesOf(iso CountryISO) (States, error)
}

func (g *general) CountryName(iso CountryISO) string {
	return g.countries[iso]
}

func (g *general) CurrencyOf(iso CurrencyISO) Currency {
	return g.currencies[iso]
}

func (g *general) StatesOf(iso CountryISO) (States, error) {
	return fetchStateList(g.core, iso)
}

func New(c core.Core) (General, error) {
	curr, err := fetchCurrencyDB(c)
	if err != nil {
		return nil, err
	}
	cntrs, err := fetchCountryDB(c)
	if err != nil {
		return nil, err
	}
	return &general{
		core:       c,
		currencies: curr,
		countries:  cntrs,
	}, nil
}

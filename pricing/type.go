package pricing

import "github.com/mrehanabbasi/go-resellerclub/core"

type CustomerPrice map[string]map[string]map[string]float64

type ResellerPrice map[string]map[string]map[string]map[string]map[string]string

type ResellerCostPrice map[string]map[string]map[string]core.JSONFloat

type PromoPrice map[string]string

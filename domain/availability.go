package domain

import "github.com/mrehanabbasi/go-logicboxes/core"

type Registration struct {
	Key    core.DomainKey     `json:"classkey"`
	Status RegistrationStatus `json:"status"`
}

type Availabilities map[string]Registration

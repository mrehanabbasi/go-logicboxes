package domain

import (
	"net/http"
	"os"
	"testing"

	"github.com/mrehanabbasi/go-resellerclub/core"
	"github.com/stretchr/testify/require"
)

var d = New(core.New(core.Config{
	ResellerId:   os.Getenv("RESELLER_ID"),
	APIKey:       os.Getenv("API_KEY"),
	IsProduction: false,
}, http.DefaultClient))

var (
	domainName = os.Getenv("TEST_DOMAIN_NAME")
	orderID    = os.Getenv("TEST_ORDER_ID")
	cns        = os.Getenv("TEST_CNS")
	authCode   = os.Getenv("TEST_AUTH_CODE")
)

func TestSuggestNames(t *testing.T) {
	res, err := d.SuggestNames("domain", "", false, false)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestGetOrderID(t *testing.T) {
	res, err := d.GetOrderID(domainName)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestGetRegistrationOrderDetails(t *testing.T) {
	res, err := d.GetRegistrationOrderDetails(orderID, []string{"All"})
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestModifyNameServers(t *testing.T) {
	res, err := d.ModifyNameServers(orderID, []string{"ns1.domain.asia"})
	require.NoError(t, err)
	require.NotNil(t, res)

	res, err = d.ModifyNameServers(orderID, []string{"ns2.domain.asia"})
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestAddChildNameServer(t *testing.T) {
	res, err := d.AddChildNameServer(orderID, "new."+domainName, []string{"0.0.0.0", "1.1.1.1"})
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestModifyPrivacyProtectionStatus(t *testing.T) {
	res, err := d.ModifyPrivacyProtectionStatus(orderID, true, "some reason")
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestModifyAuthCode(t *testing.T) {
	res, err := d.ModifyAuthCode(orderID, authCode)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestApplyTheftProtectionLock(t *testing.T) {
	res, err := d.ApplyTheftProtectionLock(orderID)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestGetTheListOfLocksAppliedOnDomainName(t *testing.T) {
	res, err := d.GetTheListOfLocksAppliedOnDomainName(orderID)
	require.NoError(t, err)
	require.NotNil(t, res)
}

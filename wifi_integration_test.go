//go:build integration

package shelly_test

import (
	"testing"

	"github.com/DonRobo/shelly-go/components"
)

func TestWifiGetConfig(t *testing.T) {
	req := &components.WifiGetConfigRequest{}
	resp := req.NewTypedResponse()
	GetCallWithVerify(t, req, resp)
}

func TestWifiGetStatus(t *testing.T) {
	req := &components.WifiGetStatusRequest{}
	resp := req.NewTypedResponse()
	GetCallWithVerify(t, req, resp)
}

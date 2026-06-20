//go:build integration

package shelly_test

import (
	"testing"

	"github.com/DonRobo/go-shelly-lite"
)

func TestSwitchGetConfig(t *testing.T) {
	req := &shelly.SwitchGetConfigRequest{
		ID: 0,
	}
	resp := req.NewTypedResponse()
	GetCallWithVerify(t, req, resp)
}

func TestSwitchGetStatus(t *testing.T) {
	req := &shelly.SwitchGetStatusRequest{
		ID: 0,
	}
	resp := req.NewTypedResponse()
	GetCallWithVerify(t, req, resp)
}

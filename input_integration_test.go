//go:build integration

package shelly_test

import (
	"testing"

	"github.com/DonRobo/shelly-go/components"
)

func TestInputGetConfig(t *testing.T) {
	req := &components.InputGetConfigRequest{
		ID: 0,
	}
	resp := req.NewTypedResponse()
	GetCallWithVerify(t, req, resp)
}

func TestInputGetStatus(t *testing.T) {
	req := &components.InputGetStatusRequest{
		ID: 0,
	}
	resp := req.NewTypedResponse()
	GetCallWithVerify(t, req, resp)
}

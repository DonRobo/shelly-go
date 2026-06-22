//go:build integration

package shelly_test

import (
	"testing"

	"github.com/DonRobo/shelly-go/components"
)

func TestSysGetConfig(t *testing.T) {
	req := &components.SysGetConfigRequest{}
	resp := req.NewTypedResponse()
	GetCallWithVerify(t, req, resp)
}

func TestSysGetStatus(t *testing.T) {
	req := &components.SysGetStatusRequest{}
	resp := req.NewTypedResponse()
	GetCallWithVerify(t, req, resp)
}

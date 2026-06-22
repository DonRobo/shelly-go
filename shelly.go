package shelly

import (
	"encoding/json"

	"github.com/DonRobo/shelly-go/rpc"
	"resty.dev/v3"
)

type ShellyGetDeviceInfoRequest struct {
	// Ident is a flag specifying if extra identifying information should be displayed.
	Ident bool
}

func (r *ShellyGetDeviceInfoRequest) Method() string {
	return "Shelly.GetDeviceInfo"
}

func (r *ShellyGetDeviceInfoRequest) NewTypedResponse() *ShellyGetDeviceInfoResponse {
	return &ShellyGetDeviceInfoResponse{}
}

func (r *ShellyGetDeviceInfoRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *ShellyGetDeviceInfoRequest) Do(
	client *resty.Client,
) (
	*ShellyGetDeviceInfoResponse,
	*rpc.Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := rpc.Do(client, r, resp)
	return resp, raw, err
}

type ShellyGetDeviceInfoResponse struct {
	// ID of the device.
	ID string `json:"id"`

	// MAC of the device.
	MAC string `json:"mac"`

	// Model of the device
	Model string `json:"model"`

	// Gen is the generation of the device
	Gen json.Number `json:"gen"`

	// FW_ID is the firmware id of the device.
	FW_ID string `json:"fw_id"`

	// Ver is the version of the device firmware.
	Ver string `json:"ver"`

	// App is the application name.
	App string `json:"app"`

	// Profile is the name of the device profile (only applicable for multi-profile devices)
	Profile string `json:"profile"`

	// AuthEn is true if authentication is enabled.
	AuthEn bool `json:"auth_en"`

	// Name of the domain (null if authentication is not enabled)
	AuthDomain *string `json:"auth_domain"`

	// Present only when false. If true, device is shown in 'Discovered devices'. If false, the device is hidden.
	Discoverable *bool `json:"discoverable"`

	// Key is cloud key of the device (see note below), present only when the ident parameter is set to true.
	Key string `json:"key"`

	// Batch used to provision the device, present only when the ident parameter is set to true.
	Batch string `json:"batch"`

	// FW_SBits are shelly internal flags, present only when the ident parameter is set to true.
	FW_SBits string
}

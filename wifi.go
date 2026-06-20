package shelly

import "resty.dev/v3"

type WifiStatus struct {
	// StaIP is the IP of the device in the network (null if disconnected).
	StaIP *string `json:"sta_ip,omitempty"`

	// Status of the connection. Range of values: disconnected, connecting, connected, got ip
	Status string `json:"status,omitempty"`

	// SSID of the network (null if disconnected)
	SSID *string `json:"ssid,omitempty"`

	// RSSI is the strength of the signal in dBms.
	RRSI *float64 `json:"rssi,omitempty"`

	// APClientCount is the number of clients connected to the access point. Present only when
	// AP is enabled and range extender functionality is present and enabled.
	APClientCount *int `json:"ap_client_count,omitempty"`
}

type WifiGetStatusRequest struct{}

func (r *WifiGetStatusRequest) Method() string {
	return "Wifi.GetStatus"
}

func (r *WifiGetStatusRequest) NewTypedResponse() *WifiStatus {
	return &WifiStatus{}
}

func (r *WifiGetStatusRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *WifiGetStatusRequest) Do(
	client *resty.Client,
) (
	*WifiStatus,
	*Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := Do(client, r, resp)
	return resp, raw, err
}

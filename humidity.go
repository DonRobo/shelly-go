package shelly

import "resty.dev/v3"

// HumidityGetStatusRequst contains parameters for the Humidity.GetStatus RPC request.
type HumidityGetStatusRequest struct {
	// ID of the humidity component instance.
	ID int `json:"id"`
}

func (r *HumidityGetStatusRequest) Method() string {
	return "Humidity.GetStatus"
}

func (r *HumidityGetStatusRequest) NewTypedResponse() *HumidityStatus {
	return &HumidityStatus{}
}

func (r *HumidityGetStatusRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *HumidityGetStatusRequest) Do(
	client *resty.Client,
) (
	*HumidityStatus,
	*Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := Do(client, r, resp)
	return resp, raw, err
}

// HumidityStatus describes the status of humidity component instances.
type HumidityStatus struct {
	// ID of the humidity component instance.
	ID int `json:"id"`

	// RH is the relative humidity in % (null if valid value could not be obtained)
	RH *float64 `json:"rh,omitempty"`

	// Errors is a list of error events related to humidity.
	Errors []string `json:"errors,omitempty"`
}

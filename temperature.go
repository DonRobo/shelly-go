package shelly

import "resty.dev/v3"

// TemperatureGetStatusRequst contains parameters for the Temperature.GetStatus RPC request.
type TemperatureGetStatusRequest struct {
	// ID of the temperature component instance.
	ID int `json:"id"`
}

func (r *TemperatureGetStatusRequest) Method() string {
	return "Temperature.GetStatus"
}

func (r *TemperatureGetStatusRequest) NewTypedResponse() *TemperatureStatus {
	return &TemperatureStatus{}
}

func (r *TemperatureGetStatusRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *TemperatureGetStatusRequest) Do(
	client *resty.Client,
) (
	*TemperatureStatus,
	*Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := Do(client, r, resp)
	return resp, raw, err
}

// TemperatureStatus describes the status of temperature component instances.
type TemperatureStatus struct {
	// ID of the temperature component instance.
	ID int `json:"id"`

	// TC is the temperature in Celsius (null if valid value could not be obtained)
	TC *float64 `json:"tC,omitempty"`

	// TF is the temperature in Fahrenheit  (null if valid value could not be obtained)
	TF *float64 `json:"tF,omitempty"`

	// Errors is a list of error shown only if at least one error is present. May contain
	// out_of_range, read when there is problem reading sensor.
	Errors []string `json:"errors,omitempty"`
}

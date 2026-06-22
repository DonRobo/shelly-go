package components

import (
	"github.com/DonRobo/shelly-go/rpc"

	"resty.dev/v3"
)

type SwitchGetStatusRequest struct {
	// ID of the switch component instance.
	ID int `json:"id"`
}

func (r *SwitchGetStatusRequest) Method() string {
	return "Switch.GetStatus"
}

func (r *SwitchGetStatusRequest) Do(
	client *resty.Client,
) (
	*SwitchStatus,
	*rpc.Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := rpc.Do(client, r, resp)
	return resp, raw, err
}

func (r *SwitchGetStatusRequest) NewTypedResponse() *SwitchStatus {
	return &SwitchStatus{}
}

func (r *SwitchGetStatusRequest) NewResponse() any {
	return r.NewTypedResponse()
}

type SwitchSetRequest struct {
	// ID of the switch component instance.
	ID int `json:"id"`

	// On is true for switch on, false otherwise. Required
	On bool `json:"on"`

	// ToggleAfter is the number of seconds afterwhich the switch will flip-back.s
	ToggleAfter *float64 `json:"toggle_after,omitempty"`
}

func (r *SwitchSetRequest) Method() string {
	return "Switch.Set"
}

func (r *SwitchSetRequest) Do(
	client *resty.Client,
) (
	*SwitchActionResponse,
	*rpc.Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := rpc.Do(client, r, resp)
	return resp, raw, err
}

func (r *SwitchSetRequest) NewTypedResponse() *SwitchActionResponse {
	return &SwitchActionResponse{}
}

func (r *SwitchSetRequest) NewResponse() any {
	return r.NewTypedResponse()
}

type SwitchToggleRequest struct {
	// ID of the switch component instance.
	ID int `json:"id"`
}

func (r *SwitchToggleRequest) Method() string {
	return "Switch.Toggle"
}

func (r *SwitchToggleRequest) Do(
	client *resty.Client,
) (
	*SwitchActionResponse,
	*rpc.Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := rpc.Do(client, r, resp)
	return resp, raw, err
}

func (r *SwitchToggleRequest) NewTypedResponse() *SwitchActionResponse {
	return &SwitchActionResponse{}
}

func (r *SwitchToggleRequest) NewResponse() any {
	return r.NewTypedResponse()
}

type SwitchStatus struct {
	// ID of the switch component instance.
	ID int `json:"id"`

	// Source of the last command, for example: init, WS_in, http, ...
	Source *string `json:"source,omitempty"`

	// Output is true if the output channel is currently on, false otherwise.
	Output *bool `json:"output,omitempty"`

	// TimerStartedAt is the unix timestamp, start time of the timer (in UTC)
	// (shown if the timer is triggered)
	TimerStartedAt *float64 `json:"timer_started_at,omitempty"`

	// TimerDuration is the number of seconds for the timer (shown if the timer
	// is triggered).
	TimerDuration *float64 `json:"timer_duration,omitempty"`

	// APower is the last measured instantaneous active power (in Watts)
	// delivered to the attached load (shown if applicable).
	APower *float64 `json:"apower,omitempty"`

	// Voltage last measured in Volts (shown if applicable).
	Voltage *float64 `json:"voltage,omitempty"`

	// Current last measured in Amperes (shown if applicable).
	Current *float64 `json:"current,omitempty"`

	// PF is the last measured power factor (shown if applicable).
	PF *float64 `json:"pf,omitempty"`

	// Freq is the last measured network frequency in Hz (shown if applicable).
	Freq *float64 `json:"freq,omitempty"`

	// AEnergy contains information about the active energy counter (shown if
	// applicable)
	AEnergy *EnergyCounters `json:"aenergy,omitempty"`

	// RetAEnergy contains information about the returned active energy counter
	// (shown if applicable)
	RetAEnergy *EnergyCounters `json:"ret_aenergy,omitempty"`

	// Temperature describes the internal temperature of the relay.
	Temperature *Temperature `json:"temperature,omitempty"`

	// Errors lists error conditions occurred. May contain overtemp, overpower,
	// overvoltage, undervoltage, (shown if at least one error is present).
	Errors []string `json:"errors,omitempty"`
}

// EnergyCounters describes active energy counters.
type EnergyCounters struct {
	// Total energy consumed in Watt-hours.
	Total float64 `json:"total"`

	// ByMinute is the energy consumption by minute (in Milliwatt-hours) for
	// the last three minutes (the lower the index of the element in the array,
	// the closer to the current moment the minute)
	ByMinute []float64 `json:"by_minute"`

	// MinuteTS is the Unix timestamp of the first second of the last minute (in UTC)
	MinuteTS float64 `json:"minute_ts,omitempty"`
}

// Temperature describes a temperature measurement.
type Temperature struct {
	// C is the temperature in Celsius (null if temperature is out of the
	// measurement range)
	C *float64 `json:"tC,omitempty"`
	// F is the temperature in Fahrenheit (null if temperature is out of the
	// measurement range)
	F *float64 `json:"tF,omitempty"`
}

type SwitchActionResponse struct {
	// WasOn is true if the switch was on before the method was executed,
	// false otherwise.
	WasOn bool `json:"was_on"`
}

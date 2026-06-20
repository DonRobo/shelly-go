package shelly

import "resty.dev/v3"

// LightGetStatusRequst contains parameters for the Light.GetStatus RPC request.
type LightGetStatusRequest struct {
	// ID of the light component instance.
	ID int `json:"id"`
}

func (r *LightGetStatusRequest) Method() string {
	return "Light.GetStatus"
}

func (r *LightGetStatusRequest) NewTypedResponse() *LightStatus {
	return &LightStatus{}
}

func (r *LightGetStatusRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *LightGetStatusRequest) Do(
	client *resty.Client,
) (
	*LightStatus,
	*Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := Do(client, r, resp)
	return resp, raw, err
}

// LightSetRequest is the parameters for the Light.Set RPC, which enables or disables a light.
type LightSetRequest struct {
	// ID of the light component instance.
	ID int `json:"id"`

	// On is true for light on, false otherwise. (optional). On or Brightness must be provided.
	On *bool `json:"on,omitempty"`

	// Brightness level (optional). On or Brightness must be provided.
	Brightness *float64 `json:"brightness,omitempty"`

	// TransitionDuration in seconds - time between change from current brightness level
	// to desired brightness level in request. (optional)
	TransitionDuration *float64 `json:"transition_duration,omitempty"`

	// ToggleAfter is the number of seconds afterwhich the light will flip-back. (optional)
	ToggleAfter *float64 `json:"toggle_after,omitempty"`
}

func (r *LightSetRequest) Method() string {
	return "Light.Set"
}

func (r *LightSetRequest) NewTypedResponse() *LightSetResponse {
	return &LightSetResponse{}
}

func (r *LightSetRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *LightSetRequest) Do(
	client *resty.Client,
) (
	*LightSetResponse,
	*Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := Do(client, r, resp)
	return resp, raw, err
}

// LightSetResponse is the response body for the Light.Set RPC.
type LightSetResponse struct{}

// LightToggleRequest contains parameters for the Light.Toggle RPC request.
type LightToggleRequest struct {
	// ID of the light component instance.
	ID int `json:"id"`
}

func (r *LightToggleRequest) Method() string {
	return "Light.Toggle"
}

func (r *LightToggleRequest) NewTypedResponse() *LightToggleResponse {
	return &LightToggleResponse{}
}

func (r *LightToggleRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *LightToggleRequest) Do(
	client *resty.Client,
) (
	*LightToggleResponse,
	*Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := Do(client, r, resp)
	return resp, raw, err
}

// LightToggleResponse is the body for the Light.Toggle RPC response.
type LightToggleResponse struct{}

// LightStatus describes the status of light component instances.
type LightStatus struct {
	// ID of the light component instance.
	ID int `json:"id"`

	// Source of the last command, for example: init, WS_in, http, ...
	Source *string `json:"source,omitempty"`

	// Output is true if the output channel is currently on, false otherwise.
	Output *bool `json:"output,omitempty"`

	// Brightness level (in percent)
	Brightness *float64 `json:"brightness,omitempty"`

	// TimerStartedAt is the unix timestamp, start time of the timer (in UTC)
	// (shown if the timer is triggered)
	TimerStartedAt *float64 `json:"timer_started_at,omitempty"`

	// TimerDuration is the number of seconds for the timer (shown if the timer
	// is triggered).
	TimerDuration *float64 `json:"timer_duration,omitempty"`

	// Transition provides information about the transition (shown if transition is triggered).
	Transition *LightTransitionStatus `json:"transition,omitempty"`

	// Temperature describes the internal temperature of the relay.
	Temperature *Temperature `json:"temperature,omitempty"`
}

// LightTransitionStatus provides information about the transition (shown if transition
// is triggered).
type LightTransitionStatus struct {
	// Target describes the desired result of the transition.
	Target LightTransitionTargetStatus `json:"target,omitempty"`

	// StartedAt is the unix timestamp start time of the transition (in UTC).
	StartedAt *float64 `json:"started_at,omitempty"`

	// Duration of the transition in seconds.
	Duration *float64 `json:"duration,omitempty"`
}

// LightTransitionTargetStatus describes the desired result of the transition.
type LightTransitionTargetStatus struct {
	// Output is true if the output channel becomes on, false otherwise
	Output bool `json:"output"`

	// Brightness level (in percent).
	Brightness *float64 `json:"brightness,omitempty"`
}

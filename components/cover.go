package components

import (
	"github.com/DonRobo/shelly-go/rpc"

	"resty.dev/v3"
)

type CoverGetStatusRequest struct {
	// ID of the cover component instance.
	ID int `json:"id"`
}

func (r *CoverGetStatusRequest) Method() string {
	return "Cover.GetStatus"
}

func (r *CoverGetStatusRequest) NewTypedResponse() *CoverStatus {
	return &CoverStatus{}
}

func (r *CoverGetStatusRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *CoverGetStatusRequest) Do(
	client *resty.Client,
) (
	*CoverStatus,
	*rpc.Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := rpc.Do(client, r, resp)
	return resp, raw, err
}

// CoverStatus describes the current state of the Cover.
type CoverStatus struct {
	// ID of the cover component instance.
	ID int `json:"id"`

	// Source of the last command, for example: init, WS_in, http, ...
	Source *string `json:"source,omitempty"`

	// State describes the current state of the cover device. One of open (Cover is
	// fully open), closed (Cover is fully closed), opening (Cover is actively opening),
	// closing (Cover is actively closing), stopped (Cover is not moving, and is neither
	// fully open nor fully closed, or the open/close state is unknown), calibrating
	// (Cover is performing a calibration procedure).
	State *string `json:"state,omitempty"`

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

	// CurrentPos is the current position current position in percent from 0 (fully
	// closed) to 100 (fully open); null if the position is unknown. Only present if
	// Cover is calibrated.
	CurrentPos *float64 `json:"current_pos,omitempty"`

	// TargetPos is only present if Cover is calibrated and is actively moving to a
	// requested position in either open or closed directions. Represents the target
	// position in percent from 0 (fully closed) to 100 (fully open); null if target
	// position has been reached or the movement was canceled.
	TargetPos *float64 `json:"target_pos,omitempty"`

	// MoveTimeout is the timeout in seconds until the cover stops regardless of completion.
	// Only present if Cover is actively moving in either open or closed directions.
	MoveTimeout *float64 `json:"move_timeout,omitempty"`

	// MoveStartedAt represents the time at which the movement has begun. Only present if
	// Cover is actively moving in either open or closed directions.
	MoveStartedAt *float64 `json:"move_started_at,omitempty"`

	// PosControl is false if Cover is not calibrated and only discrete open/close is
	// possible; true if Cover is calibrated and can be commanded to go to arbitrary
	// positions between fully open and fully closed.
	PosControl *bool `json:"pos_control,omitempty"`

	// LastDirection is the direction of the last movement: open/close or null when unknown.
	LastDirection *string `json:"last_direction,omitempty"`

	// Temperature describes the internal temperature of the cover instance. Only present if
	// a temperature monitor is associated with the Cover instance
	Temperature *Temperature `json:"temperature,omitempty"`

	// Errors lists error conditions occurred. May contain overtemp, overpower,
	// overvoltage, undervoltage, (shown if at least one error is present).
	Errors []string `json:"errors,omitempty"`
}

// CoverCalibrateRequest causes the device to enter calibration mode. See:
// - https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Cover#covercalibrate
// - https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Cover#calibration-kb
type CoverCalibrateRequest struct {
	// ID of the cover component instance.
	ID int `json:"id"`
}

func (r *CoverCalibrateRequest) Method() string {
	return "Cover.Calibrate"
}

func (r *CoverCalibrateRequest) NewTypedResponse() *CoverCalibrateRespose {
	return &CoverCalibrateRespose{}
}

func (r *CoverCalibrateRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *CoverCalibrateRequest) Do(
	client *resty.Client,
) (
	*CoverCalibrateRespose,
	*rpc.Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := rpc.Do(client, r, resp)
	return resp, raw, err
}

// CoverCalibrateRespose is the RPC response for Cover.Calibrate.
type CoverCalibrateRespose struct{}

// CoverOpenRequest causes the device to open the cover instance.
type CoverOpenRequest struct {
	// ID of the cover component instance.
	ID int `json:"id"`

	// Duration (seconds) if provided, Cover will move in the open direction for the
	// specified time. duration must be in the range [0.1..maxtime_open].
	// If duration is not provided, Cover will fully open, unless it times out because
	// of maxtime_open first.
	Duration *float64 `json:"duration,omitempty"`
}

func (r *CoverOpenRequest) Method() string {
	return "Cover.Open"
}

func (r *CoverOpenRequest) NewTypedResponse() *CoverOpenResponse {
	return &CoverOpenResponse{}
}

func (r *CoverOpenRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *CoverOpenRequest) Do(
	client *resty.Client,
) (
	*CoverOpenResponse,
	*rpc.Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := rpc.Do(client, r, resp)
	return resp, raw, err
}

// CoverOpenResponse is the RPC response for Cover.Open.
type CoverOpenResponse struct{}

// CoverCloseRequest causes the device to close the cover instance.
type CoverCloseRequest struct {
	// ID of the cover component instance.
	ID int `json:"id"`

	// Duration (seconds) if provided, Cover will move in the close direction for the
	// specified time. duration must be in the range [0.1..maxtime_open].
	// If duration is not provided, Cover will fully close, unless it times out because
	// of maxtime_close first.
	Duration *float64 `json:"duration,omitempty"`
}

func (r *CoverCloseRequest) Method() string {
	return "Cover.Close"
}

func (r *CoverCloseRequest) NewTypedResponse() *CoverCloseResponse {
	return &CoverCloseResponse{}
}

func (r *CoverCloseRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *CoverCloseRequest) Do(
	client *resty.Client,
) (
	*CoverCloseResponse,
	*rpc.Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := rpc.Do(client, r, resp)
	return resp, raw, err
}

// CoverCloseResponse is the RPC response for Cover.Close.
type CoverCloseResponse struct{}

// CoverStopRequest causes the device to stop in progress actions for the cover instance.
type CoverStopRequest struct {
	// ID of the cover component instance.
	ID int `json:"id"`
}

func (r *CoverStopRequest) Method() string {
	return "Cover.Stop"
}

func (r *CoverStopRequest) NewTypedResponse() *CoverStopResponse {
	return &CoverStopResponse{}
}

func (r *CoverStopRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *CoverStopRequest) Do(
	client *resty.Client,
) (
	*CoverStopResponse,
	*rpc.Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := rpc.Do(client, r, resp)
	return resp, raw, err
}

// CoverStopResponse is the RPC response for Cover.Stop.
type CoverStopResponse struct{}

// CoverGoToPositionRequest causes the device to travel to the specified position.
type CoverGoToPositionRequest struct {
	// ID of the cover component instance.
	ID int `json:"id"`

	// Pos represents target position in %, allowed range [0..100].
	// Pos is mutually exclusive with Rel. Rel or Pos is required, but both may not be set.
	Pos *float64 `json:"pos,omitempty"`

	// Rel represents a relative move in %, allowed range [-100..100] Cover will move
	// to a target_position = current_position + rel. If the value of rel is so big that
	// it results in overshoot (i.e. target_position is beyond fully open / fully closed),
	// target_position will be silently capped to fully open / fully closed.
	// Rel is mutually exclusive with Pos. Rel or Pos is required, but both may not be set.
	Rel *float64 `json:"rel,omitempty"`
}

func (r *CoverGoToPositionRequest) Method() string {
	return "Cover.GoToPosition"
}

func (r *CoverGoToPositionRequest) NewTypedResponse() *CoverGoToPositionResponse {
	return &CoverGoToPositionResponse{}
}

func (r *CoverGoToPositionRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *CoverGoToPositionRequest) Do(
	client *resty.Client,
) (
	*CoverGoToPositionResponse,
	*rpc.Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := rpc.Do(client, r, resp)
	return resp, raw, err
}

// CoverGoToPositionResponse is the RPC response for Cover.GoToPosition.
type CoverGoToPositionResponse struct{}

// CoverResetCountersRequest resets counters for the cover.
type CoverResetCountersRequest struct {
	// ID of the cover component instance.
	ID int `json:"id"`

	// Type describes which counters should be reset.
	Type []string `json:"type,omitempty"`
}

func (r *CoverResetCountersRequest) Method() string {
	return "Cover.ResetCounters"
}

func (r *CoverResetCountersRequest) NewTypedResponse() *CoverResetCountersResponse {
	return &CoverResetCountersResponse{}
}

func (r *CoverResetCountersRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *CoverResetCountersRequest) Do(
	client *resty.Client,
) (
	*CoverResetCountersResponse,
	*rpc.Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := rpc.Do(client, r, resp)
	return resp, raw, err
}

// CoverResetCountersResponse is the RPC response for Cover.ResetCounters.
type CoverResetCountersResponse struct {
	// AEnergy contains information about the active energy counter prior to reset.
	AEnergy *EnergyCounters `json:"aenergy,omitempty"`
}

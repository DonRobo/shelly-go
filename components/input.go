package components

import (
	"encoding/json"
	"fmt"

	"github.com/DonRobo/shelly-go/rpc"

	"resty.dev/v3"
)

type InputGetStatusRequest struct {
	// ID of the switch component instance.
	ID int `json:"id"`
}

func (r *InputGetStatusRequest) Method() string {
	return "Input.GetStatus"
}

func (r *InputGetStatusRequest) NewTypedResponse() *InputStatus {
	return &InputStatus{}
}

func (r *InputGetStatusRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *InputGetStatusRequest) Do(
	client *resty.Client,
) (
	*InputStatus,
	*rpc.Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := rpc.Do(client, r, resp)
	return resp, raw, err
}

type InputStatus struct {
	// ID of the input component instance.
	ID int `json:"id"`

	// State of the input (null if the input instance is stateless, i.e. for type button)
	// (only for type switch, button).
	State *bool `json:"state,omitempty"`

	// Percent is the analog value in percent (null if the valid value could not be obtained)
	// (only for type "analog").
	Percent *float64 `json:"percent,omitempty"`

	// XPercent is percent transformed with config.xpercent.expr. Present only when both
	// config.xpercent.expr and config.xpercent.unit are set to non-empty values. null if
	// config.xpercent.expr can not be evaluated.
	// (only for type "analog").
	XPercent *float64 `json:"xpercent,omitempty"`

	// Errors is shown only if at least one error is present. May contain out_of_range, read.
	Errors []string `json:"errors,omitempty"`
}

// InputXPercent is value transformation config for status.percent.
type InputXPercent struct {
	// Expr is a JS expression containing x, where x is the raw value to be transformed
	// (status.percent), for example "x+1". Accepted range: null or [0..100] chars. Both
	// null and "" mean value transformation is disabled.
	Expr *string `json:"expr,omitempty"`

	// Unit of the transformed value (status.xpercent), for example, "m/s".
	// Accepted range: null or [0..20] chars. Both null and "" mean value transformation
	// is disabled.
	Unit *string `json:"unit,omitempty"`
}

type InputCheckExpressionRequest struct {
	// Expr is the JS expression to evaluate.
	Expr string `json:"expr,omitempty"`

	// Inputs on which to apply expr. Elements are allowed to be null
	Inputs []*float64 `json:"inputs,omitempty"`
}

func (r *InputCheckExpressionRequest) Method() string {
	return "Input.CheckExpression"
}

func (r *InputCheckExpressionRequest) NewTypedResponse() *InputCheckExpressionResponse {
	return &InputCheckExpressionResponse{}
}

func (r *InputCheckExpressionRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *InputCheckExpressionRequest) Do(
	client *resty.Client,
) (
	*InputCheckExpressionResponse,
	*rpc.Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := rpc.Do(client, r, resp)
	return resp, raw, err
}

type InputCheckExpressionResponse struct {
	Results []InputCheckExpressionResult
}

type InputCheckExpressionResult struct {
	Input *float64

	Output *float64

	Error *string
}

func (r *InputCheckExpressionResult) UnmarshalJSON(b []byte) error {
	var got []interface{}
	if err := json.Unmarshal(b, &got); err != nil {
		return err
	}
	if len(got) >= 1 {
		if n, ok := got[0].(json.Number); ok {
			f64, err := n.Float64()
			if err != nil {
				return fmt.Errorf("parsing input result: %w", err)
			}
			r.Input = rpc.Float64Ptr(f64)
		} else if n, ok := got[0].(float64); ok {
			r.Input = rpc.Float64Ptr(n)
		}

	}
	if len(got) >= 2 {
		if n, ok := got[1].(json.Number); ok {
			f64, err := n.Float64()
			if err != nil {
				return fmt.Errorf("parsing output result: %w", err)
			}
			r.Output = rpc.Float64Ptr(f64)
		} else if n, ok := got[1].(float64); ok {
			r.Output = rpc.Float64Ptr(n)
		}
	}
	if len(got) >= 3 {
		if s, ok := got[2].(string); ok {
			r.Error = rpc.StrPtr(s)
		}
	}

	return nil
}

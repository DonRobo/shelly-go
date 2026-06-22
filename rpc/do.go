package rpc

import (
	"encoding/json"
	"fmt"

	"resty.dev/v3"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

type Frame struct {
	ID int64 `json:"id,omitempty"`

	// Request
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`

	// Response
	Result json.RawMessage `json:"result,omitempty"`
	Error  *Error          `json:"error,omitempty"`
}

func Do(
	client *resty.Client,
	req RPCRequestBody,
	resp any,
) (*Frame, error) {
	args, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling shelly rpc request: %w", err)
	}
	command := &Frame{
		ID:     1,
		Method: req.Method(),
		Params: json.RawMessage(args),
	}
	restyResp, err := client.R().
		SetBody(command).
		SetContentType("application/json").
		SetResult(&Frame{}).
		Post("/rpc")
	rawResp := restyResp.Result().(*Frame)
	if err != nil {
		return rawResp, fmt.Errorf("making shelly rpc request: %w", err)
	}
	if rawResp.Error != nil {
		return rawResp, &BadStatusWithMessageError{Status: ShellyErrorCode(rawResp.Error.Code), Msg: rawResp.Error.Message}
	}
	if err := json.Unmarshal(rawResp.Result, resp); err != nil {
		return rawResp, fmt.Errorf("failed to unmarshal response body: %w from %s", err, string(rawResp.Result))
	}
	return rawResp, err
}

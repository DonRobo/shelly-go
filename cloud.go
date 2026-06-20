package shelly

import "resty.dev/v3"

type CloudStatus struct {
	Connected bool `json:"connected"`
}

type CloudGetStatusRequest struct{}

func (r *CloudGetStatusRequest) Method() string {
	return "Cloud.GetStatus"
}

func (r *CloudGetStatusRequest) NewTypedResponse() *RPCEmptyResponse {
	return &RPCEmptyResponse{}
}

func (r *CloudGetStatusRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *CloudGetStatusRequest) Do(
	client *resty.Client,
) (
	*RPCEmptyResponse,
	*Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := Do(client, r, resp)
	return resp, raw, err
}

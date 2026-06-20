package shelly

import "resty.dev/v3"

type EthStatus struct {
	// IP of the device in the network.
	IP *string `json:"ip"`
}

type EthGetStatusRequest struct{}

func (r *EthGetStatusRequest) Method() string {
	return "Eth.GetStatus"
}

func (r *EthGetStatusRequest) NewTypedResponse() *EthStatus {
	return &EthStatus{}
}

func (r *EthGetStatusRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *EthGetStatusRequest) Do(
	client *resty.Client,
) (
	*EthStatus,
	*Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := Do(client, r, resp)
	return resp, raw, err
}

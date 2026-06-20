package shelly

import "resty.dev/v3"

type BLEStatus struct{}

type BLEGetStatusRequest struct{}

func (r *BLEGetStatusRequest) Method() string {
	return "BLE.GetStatus"
}

func (r *BLEGetStatusRequest) NewTypedResponse() *BLEStatus {
	return &BLEStatus{}
}

func (r *BLEGetStatusRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *BLEGetStatusRequest) Do(
	client *resty.Client,
) (
	*BLEStatus,
	*Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := Do(client, r, resp)
	return resp, raw, err
}

//go:build integration

package shelly_test

import (
	"testing"

	shelly "github.com/DonRobo/shelly-go"

	"resty.dev/v3"
)

func TestRestyRpc(t *testing.T) {
	// Test POST request to Shelly RPC API using resty v3
	client := resty.New()
	defer client.Close()

	url := "http://192.168.1.169/rpc"
	body := map[string]interface{}{
		"id":     0,
		"method": "Shelly.GetStatus",
		"params": map[string]interface{}{},
	}

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)

	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	t.Logf("Status: %s", res.Status())
	t.Logf("Body: %s", res.String())
}

func TestShellyGoLiteWithResty(t *testing.T) {
	shellyIp := "http://192.168.1.169"
	client := resty.New()
	client.SetBaseURL(shellyIp)
	defer client.Close()

	req := &shelly.ShellyGetStatusRequest{}
	statusResp, _, err := req.Do(client)
	if err != nil {
		t.Fatalf("querying device status: %v", err)
	}
	t.Logf("Device status: %+v", statusResp)
}

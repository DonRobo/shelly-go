package shelly

import (
	"errors"
	"strings"

	"resty.dev/v3"
)

// MQTT_SSL_CA is a type to differentiate between not-set (empty string), null (no TLS), and string
// values.
type MQTT_SSL_CA string

const (
	// MQTT_SSL_CA_NULL will disable the TLS on the MQTT connection
	MQTT_SSL_CA_NULL MQTT_SSL_CA = "null"

	// MQTT_SSL_CA_NOT_SET will not send a value for the `ssl_ca` property, leaving it unchanged.
	MQTT_SSL_CA_NOT_SET MQTT_SSL_CA = ""

	// MQTT_SSL_CA_NO_VERIFY will enable TLS but CA skip verification of the server certificate.
	MQTT_SSL_CA_NO_VERIFY MQTT_SSL_CA = "*"

	// MQTT_SSL_CA_DEFAULT_CA will enable TLS with server verification against the default CA bundle.
	MQTT_SSL_CA_DEFAULT_CA MQTT_SSL_CA = "ca.pem"

	// MQTT_SSL_CA_USER_CA will enable TLS with server verification against the user-provided CA.
	// See `Shelly.PutUserCA`.
	MQTT_SSL_CA_USER_CA MQTT_SSL_CA = "user_ca.pem"
)

func (ca *MQTT_SSL_CA) UnmarshalJSON(b []byte) error {
	// NOTE if the balue is unset, this UnmarshallJSON method will not be called which is why
	// MQTT_SSL_CA_NOT_SET is absent.
	s := strings.TrimSpace(string(b))
	switch s {
	case string(MQTT_SSL_CA_NULL):
		*ca = MQTT_SSL_CA_NULL
		return nil
	case `"` + string(MQTT_SSL_CA_NO_VERIFY) + `"`:
		*ca = MQTT_SSL_CA_NO_VERIFY
		return nil
	case `"` + string(MQTT_SSL_CA_DEFAULT_CA) + `"`:
		*ca = MQTT_SSL_CA_DEFAULT_CA
		return nil
	case `"` + string(MQTT_SSL_CA_USER_CA) + `"`:
		*ca = MQTT_SSL_CA_USER_CA
		return nil
	default:
		return errors.New("unknown value for MQTTConfig.SSL_CA")
	}
}

func (ca *MQTT_SSL_CA) MarshalJSON() ([]byte, error) {
	if ca == nil || *ca == MQTT_SSL_CA_NULL {
		return []byte("null"), nil
	}
	return []byte(`"` + *ca + `"`), nil
}

type MQTTGetStatusRequest struct{}

// Method returns the method name.
func (r *MQTTGetStatusRequest) Method() string {
	return "MQTT.GetStatus"
}

func (r *MQTTGetStatusRequest) NewTypedResponse() *MQTTStatus {
	return &MQTTStatus{}
}

func (r *MQTTGetStatusRequest) NewResponse() any {
	return r.NewTypedResponse()
}

func (r *MQTTGetStatusRequest) Do(
	client *resty.Client,
) (
	*MQTTStatus,
	*Frame,
	error,
) {
	resp := r.NewTypedResponse()
	raw, err := Do(client, r, resp)
	return resp, raw, err
}

type MQTTStatus struct {
	Connected bool `json:"connected"`
}

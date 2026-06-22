//go:build integration

package shelly

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestShellyGetStatusResponseUnmarshall validates that the generated aggregate
// and component status structs unmarshal real-device Shelly.GetStatus payloads
// without losing structure or typing. The fixtures are verbatim captures from
// three physical devices; the assertions spot-check the routing (singletons vs
// keyed slices) and the richer nested status shapes the codegen now produces
// (typed aenergy/temperature sub-structs, typed arrays, expanded objects).
func TestShellyGetStatusResponseUnmarshall(t *testing.T) {
	tcs := []struct {
		name  string
		input string
	}{
		{
			name: "pro pm 4",
			input: `{
				"ble": {},
				"cloud": {
				  "connected": true
				},
				"eth": {
				  "ip": null
				},
				"input:0": {
				  "id": 0,
				  "state": false
				},
				"input:1": {
				  "id": 1,
				  "state": false
				},
				"input:2": {
				  "id": 2,
				  "state": false
				},
				"input:3": {
				  "id": 3,
				  "state": false
				},
				"mqtt": {
				  "connected": false
				},
				"switch:0": {
				  "id": 0,
				  "source": "timer",
				  "output": false,
				  "apower": 0,
				  "voltage": 120.8,
				  "freq": 60,
				  "current": 0,
				  "pf": 0,
				  "aenergy": {
					"total": 1342.238,
					"by_minute": [
					  0,
					  0,
					  0
					],
					"minute_ts": 1703811193
				  },
				  "ret_aenergy": {
					"total": 0,
					"by_minute": [
					  0,
					  0,
					  0
					],
					"minute_ts": 1703811193
				  },
				  "temperature": {
					"tC": 41.3,
					"tF": 106.3
				  }
				},
				"switch:1": {
				  "id": 1,
				  "source": "HTTP_in",
				  "output": true,
				  "apower": 83.9,
				  "voltage": 120.8,
				  "freq": 60,
				  "current": 1.143,
				  "pf": 0.61,
				  "aenergy": {
					"total": 102650.773,
					"by_minute": [
					  344.204,
					  1475.177,
					  1474.888
					],
					"minute_ts": 1703811193
				  },
				  "ret_aenergy": {
					"total": 0,
					"by_minute": [
					  0,
					  0,
					  0
					],
					"minute_ts": 1703811193
				  },
				  "temperature": {
					"tC": 41.3,
					"tF": 106.3
				  }
				},
				"switch:2": {
				  "id": 2,
				  "source": "HTTP_in",
				  "output": true,
				  "apower": 210.3,
				  "voltage": 120.9,
				  "freq": 60,
				  "current": 1.741,
				  "pf": 1,
				  "aenergy": {
					"total": 69346.948,
					"by_minute": [
					  840.825,
					  3605.178,
					  3624.834
					],
					"minute_ts": 1703811193
				  },
				  "ret_aenergy": {
					"total": 0,
					"by_minute": [
					  0,
					  0,
					  0
					],
					"minute_ts": 1703811193
				  },
				  "temperature": {
					"tC": 41.3,
					"tF": 106.3
				  }
				},
				"switch:3": {
				  "id": 3,
				  "source": "init",
				  "output": false,
				  "apower": 0,
				  "voltage": 120.9,
				  "freq": 60,
				  "current": 0,
				  "pf": 0,
				  "aenergy": {
					"total": 13.264,
					"by_minute": [
					  0,
					  0,
					  0
					],
					"minute_ts": 1703811193
				  },
				  "ret_aenergy": {
					"total": 0,
					"by_minute": [
					  0,
					  0,
					  0
					],
					"minute_ts": 1703811193
				  },
				  "temperature": {
					"tC": 41.3,
					"tF": 106.3
				  }
				},
				"sys": {
				  "mac": "C8F09E87D088",
				  "restart_required": false,
				  "time": "19:53",
				  "unixtime": 1703811195,
				  "uptime": 97431,
				  "ram_size": 241028,
				  "ram_free": 100452,
				  "fs_size": 524288,
				  "fs_free": 196608,
				  "cfg_rev": 26,
				  "kvs_rev": 1,
				  "schedule_rev": 0,
				  "webhook_rev": 0,
				  "available_updates": {},
				  "reset_reason": 3
				},
				"ui": {},
				"wifi": {
				  "sta_ip": "192.168.1.24",
				  "status": "got ip",
				  "ssid": "PickleTown",
				  "rssi": -36
				},
				"ws": {
				  "connected": false
				}
			  }`,
		},
		{
			name: "pro 3",
			input: `{
				"ble": {},
				"cloud": {
				  "connected": true
				},
				"eth": {
				  "ip": null
				},
				"input:0": {
				  "id": 0,
				  "state": false
				},
				"input:1": {
				  "id": 1,
				  "state": false
				},
				"input:2": {
				  "id": 2,
				  "state": false
				},
				"mqtt": {
				  "connected": false
				},
				"switch:0": {
				  "id": 0,
				  "source": "init",
				  "output": false,
				  "temperature": {
					"tC": 35.7,
					"tF": 96.2
				  }
				},
				"switch:1": {
				  "id": 1,
				  "source": "timer",
				  "output": false,
				  "temperature": {
					"tC": 35.7,
					"tF": 96.2
				  }
				},
				"switch:2": {
				  "id": 2,
				  "source": "timer",
				  "output": false,
				  "temperature": {
					"tC": 35.7,
					"tF": 96.2
				  }
				},
				"sys": {
				  "mac": "C8F09E883630",
				  "restart_required": false,
				  "time": "19:52",
				  "unixtime": 1703811156,
				  "uptime": 98059,
				  "ram_size": 243420,
				  "ram_free": 104384,
				  "fs_size": 524288,
				  "fs_free": 212992,
				  "cfg_rev": 16,
				  "kvs_rev": 0,
				  "schedule_rev": 0,
				  "webhook_rev": 0,
				  "available_updates": {},
				  "reset_reason": 3
				},
				"wifi": {
				  "sta_ip": "192.168.1.23",
				  "status": "got ip",
				  "ssid": "PickleTown",
				  "rssi": -22
				},
				"ws": {
				  "connected": false
				}
			  }`,
		},
		{
			name: "shelly plus ht",
			input: `{
				"ble": {},
				"cloud": {
					"connected": true
				},
				"devicepower:0": {
					"id": 0,
					"battery": {
						"V": 0.43,
						"percent": 0
					},
					"external": {
						"present": true
					}
				},
				"ht_ui": {},
				"humidity:0": {
					"id": 0,
					"rh": 59.4
				},
				"mqtt": {
					"connected": true
				},
				"sys": {
					"mac": "C049EF8BB8F8",
					"restart_required": false,
					"time": "09:29",
					"unixtime": 1733063380,
					"uptime": 6,
					"ram_size": 247452,
					"ram_free": 159420,
					"fs_size": 458752,
					"fs_free": 176128,
					"cfg_rev": 15,
					"kvs_rev": 0,
					"webhook_rev": 2,
					"available_updates": {},
					"wakeup_reason": {
						"boot": "deepsleep_wake",
						"cause": "status_update"
					},
					"wakeup_period": 600,
					"reset_reason": 8
				},
				"temperature:0": {
					"id": 0,
					"tC": 2.7,
					"tF": 36.8
				},
				"wifi": {
					"sta_ip": "192.168.1.199",
					"status": "got ip",
					"ssid": "PickleTown_Garage",
					"rssi": -35
				},
				"ws": {
					"connected": false
				}
			}`,
		},
	}

	by := make(map[string]ShellyGetStatusResponse, len(tcs))
	for _, tc := range tcs {
		var got ShellyGetStatusResponse
		require.NoError(t, json.Unmarshal([]byte(tc.input), &got), tc.name)
		by[tc.name] = got
	}

	t.Run("pro pm 4", func(t *testing.T) {
		r := by["pro pm 4"]

		// Singletons routed to pointer fields.
		require.NotNil(t, r.BLE)
		require.NotNil(t, r.Sys)
		require.NotNil(t, r.Cloud)
		require.NotNil(t, r.MQTT)
		require.NotNil(t, r.Wifi)
		require.NotNil(t, r.Eth)
		require.NotNil(t, r.Ws)
		assert.True(t, *r.Cloud.Connected)
		assert.False(t, *r.MQTT.Connected)
		assert.False(t, *r.Ws.Connected)
		assert.Nil(t, r.Eth.IP)
		assert.Equal(t, "PickleTown", *r.Wifi.SSID)
		assert.Equal(t, float64(-36), *r.Wifi.Rssi)
		assert.Equal(t, "C8F09E87D088", r.Sys.Mac)

		// Keyed components routed to slices, uncapped (4 switches + 4 inputs).
		require.Len(t, r.Inputs, 4)
		require.Len(t, r.Switches, 4)
		for i, in := range r.Inputs {
			assert.Equal(t, i, in.ID)
			assert.False(t, *in.State)
		}

		// Rich nested status: typed aenergy (typed by_minute array) + temperature.
		s1 := r.Switches[1]
		assert.Equal(t, 1, s1.ID)
		assert.Equal(t, "HTTP_in", *s1.Source)
		assert.True(t, *s1.Output)
		assert.Equal(t, 83.9, *s1.Apower)
		require.NotNil(t, s1.Aenergy)
		assert.Equal(t, 102650.773, *s1.Aenergy.Total)
		assert.Equal(t, []float64{344.204, 1475.177, 1474.888}, s1.Aenergy.ByMinute)
		assert.Equal(t, float64(1703811193), *s1.Aenergy.MinuteTs)
		require.NotNil(t, s1.RetAenergy)
		assert.Equal(t, float64(0), *s1.RetAenergy.Total)
		require.NotNil(t, s1.Temperature)
		assert.Equal(t, 41.3, *s1.Temperature.TC)
		assert.Equal(t, 106.3, *s1.Temperature.TF)
	})

	t.Run("pro 3", func(t *testing.T) {
		r := by["pro 3"]

		require.NotNil(t, r.Sys)
		require.NotNil(t, r.Cloud)
		require.NotNil(t, r.Wifi)
		assert.Equal(t, float64(-22), *r.Wifi.Rssi)
		assert.Equal(t, "C8F09E883630", r.Sys.Mac)

		require.Len(t, r.Inputs, 3)
		require.Len(t, r.Switches, 3)

		s0 := r.Switches[0]
		assert.Equal(t, "init", *s0.Source)
		assert.False(t, *s0.Output)
		require.NotNil(t, s0.Temperature)
		assert.Equal(t, 35.7, *s0.Temperature.TC)
		assert.Equal(t, 96.2, *s0.Temperature.TF)
		// No power section on this device -> those pointers stay nil.
		assert.Nil(t, s0.Apower)
		assert.Nil(t, s0.Aenergy)
	})

	t.Run("shelly plus ht", func(t *testing.T) {
		r := by["shelly plus ht"]

		require.NotNil(t, r.Sys)
		require.NotNil(t, r.Cloud)
		require.NotNil(t, r.MQTT)
		require.NotNil(t, r.Wifi)
		assert.True(t, *r.Cloud.Connected)
		assert.True(t, *r.MQTT.Connected)
		assert.Equal(t, "PickleTown_Garage", *r.Wifi.SSID)

		// Battery-device sys extras.
		require.NotNil(t, r.Sys.WakeUpReason)
		assert.Equal(t, "deepsleep_wake", r.Sys.WakeUpReason.Boot)
		assert.Equal(t, "status_update", r.Sys.WakeUpReason.Cause)
		assert.Equal(t, float64(600), r.Sys.WakeUpPeriod)

		require.Len(t, r.DevicePowers, 1)
		require.Len(t, r.Humidities, 1)
		require.Len(t, r.Temperatures, 1)

		dp := r.DevicePowers[0]
		assert.Equal(t, 0, dp.ID)
		require.NotNil(t, dp.Battery)
		assert.Equal(t, float64(0), *dp.Battery.Percent)
		// external was a bare "object" in the docs; expandStatusObject gives it a
		// typed Present field.
		require.NotNil(t, dp.External)
		assert.True(t, *dp.External.Present)

		assert.Equal(t, 59.4, *r.Humidities[0].Rh)
		assert.Equal(t, 2.7, *r.Temperatures[0].TC)
		assert.Equal(t, 36.8, *r.Temperatures[0].TF)
	})
}

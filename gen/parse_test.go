package gen

import (
	"testing"
)

func TestGoName(t *testing.T) {
	cases := map[string]string{
		"id":            "ID",
		"in_mode":       "InMode",
		"initial_state": "InitialState",
		"input_id":      "InputID",
		"mac":           "MAC",
		"ssid":          "SSID",
		"power_thr":     "PowerThr",
		"wifi":          "WiFi",
		"ap":            "AP",
	}
	for in, want := range cases {
		if got := goName(in); got != want {
			t.Errorf("goName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNormaliseType(t *testing.T) {
	cases := []struct {
		in       string
		wantBase string
		wantNull bool
	}{
		{"number", "number", false},
		{"string", "string", false},
		{"string or null", "string", true},
		{"boolean", "boolean", false},
		{"object", "object", false},
		{"array of objects", "array", false},
		{"Cover operation in both directions", "unknown", false},
	}
	for _, c := range cases {
		base, null := normaliseType(c.in)
		if base != c.wantBase || null != c.wantNull {
			t.Errorf("normaliseType(%q) = (%q,%v), want (%q,%v)", c.in, base, null, c.wantBase, c.wantNull)
		}
	}
}

func TestEnumValues(t *testing.T) {
	got := enumValues("Mode of the associated input. Range of values: momentary, follow, flip, detached, cycle (if applicable)")
	want := []string{"momentary", "follow", "flip", "detached", "cycle"}
	if len(got) != len(want) {
		t.Fatalf("enumValues len = %d (%v), want %d (%v)", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("enumValues[%d] = %q, want %q", i, got[i], want[i])
		}
	}
	if enumValues("just a description with no range") != nil {
		t.Errorf("expected nil enum for plain description")
	}
}

func TestValidKey(t *testing.T) {
	valid := []string{"id", "in_mode", "counts.power_thr", "auto_off_delay"}
	invalid := []string{"when slat control disabled", "Some Prose", "", "single value"}
	for _, k := range valid {
		if !validKey.MatchString(k) {
			t.Errorf("validKey(%q) = false, want true", k)
		}
	}
	for _, k := range invalid {
		if validKey.MatchString(k) {
			t.Errorf("validKey(%q) = true, want false", k)
		}
	}
}

// TestParsePageMinimal exercises method detection and config-table parsing on a
// minimal Docusaurus-shaped fragment, including a wide table whose enum-value
// continuation rows must be skipped.
func TestParsePageMinimal(t *testing.T) {
	doc := `<html><body>
	<h3 id="demogetconfig">Demo.GetConfig</h3>
	<h3 id="demosetconfig">Demo.SetConfig</h3>
	<h2 id="configuration">Configuration</h2>
	<table>
	  <tr><th>Property</th><th>Type</th><th>Description</th></tr>
	  <tr><td>id</td><td>number</td><td>Id of the instance</td></tr>
	  <tr><td>name</td><td>string or null</td><td>Name of the instance</td></tr>
	  <tr><td>mode</td><td>string</td><td>Range of values: a, b, c</td></tr>
	  <tr><td>in both directions</td><td>does X</td><td>continuation row</td></tr>
	</table>
	</body></html>`

	comp, err := parsePage("Demo", []byte(doc))
	if err != nil {
		t.Fatal(err)
	}
	if !comp.HasGetConfig || !comp.HasSetConfig {
		t.Fatalf("method detection failed: get=%v set=%v", comp.HasGetConfig, comp.HasSetConfig)
	}
	if !comp.Keyed {
		t.Errorf("expected Keyed=true (has id field)")
	}
	if len(comp.Fields) != 3 {
		t.Fatalf("got %d fields, want 3 (continuation row must be skipped): %+v", len(comp.Fields), comp.Fields)
	}
	if !comp.Fields[1].Nullable {
		t.Errorf("name should be nullable")
	}
	if len(comp.Fields[2].Enum) != 3 {
		t.Errorf("mode enum = %v, want [a b c]", comp.Fields[2].Enum)
	}
}

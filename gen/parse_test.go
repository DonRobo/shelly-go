package gen

import (
	"strings"
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

func TestParseRange(t *testing.T) {
	cases := []struct {
		desc     string
		min, max float64
		ok       bool
	}{
		{"Accepted range: [0.5..30]s", 0.5, 30, true},
		{"Range [1,5] where 5 is fastest", 1, 5, true},
		{"range [1 - 2147483647]", 1, 2147483647, true},
		{"value, [0 to 100] percent", 0, 100, true},
		{"Offset, [-50, 50]", -50, 50, true},
		{"range: null or [0..N]", 0, 0, false}, // non-numeric upper bound
		{"Range of values: off, on, restore_last", 0, 0, false},
		{"no range at all", 0, 0, false},
	}
	for _, c := range cases {
		min, max := parseRange(c.desc)
		if !c.ok {
			if min != nil || max != nil {
				t.Errorf("parseRange(%q) = (%v,%v), want nil", c.desc, min, max)
			}
			continue
		}
		if min == nil || max == nil || *min != c.min || *max != c.max {
			t.Errorf("parseRange(%q) = (%v,%v), want (%v,%v)", c.desc, min, max, c.min, c.max)
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

// A scalar field may carry a nested "Value | Description" table that merely
// enumerates its allowed values; that is not a sub-property table, so the field
// must stay a leaf. A real object field uses a "Property | Type | Description"
// table and must still recurse into its children.
func TestParseValueTableNotObject(t *testing.T) {
	doc := `<html><body>
	<h3 id="demogetconfig">Demo.GetConfig</h3>
	<h2 id="configuration">Configuration</h2>
	<table>
	  <tr><th>Property</th><th>Type</th><th>Description</th></tr>
	  <tr><td>in_mode</td><td>string</td><td>The mode, one of the:
	    <table><tbody><tr><th>Value</th><th>Description</th></tr>
	      <tr><td>single</td><td>one input</td></tr>
	      <tr><td>dual</td><td>two inputs</td></tr>
	    </tbody></table></td></tr>
	  <tr><td>motor</td><td>object</td><td>Motor settings:
	    <table><tbody><tr><th>Property</th><th>Type</th><th>Description</th></tr>
	      <tr><td>idle_power_thr</td><td>number</td><td>threshold</td></tr>
	    </tbody></table></td></tr>
	</table>
	</body></html>`

	comp, err := parsePage("Demo", []byte(doc))
	if err != nil {
		t.Fatal(err)
	}
	byKey := map[string]*Field{}
	for _, f := range comp.Fields {
		byKey[f.Key] = f
	}
	in, ok := byKey["in_mode"]
	if !ok || in.Type != "string" {
		t.Fatalf("in_mode should be a string leaf; fields=%v", byKey)
	}
	if want := []string{"single", "dual"}; len(in.Enum) != 2 || in.Enum[0] != want[0] || in.Enum[1] != want[1] {
		t.Errorf("in_mode enum = %v, want %v (captured from value table)", in.Enum, want)
	}
	if !strings.HasSuffix(in.Description, "one of the: single, dual.") {
		t.Errorf("in_mode description not completed from value table: %q", in.Description)
	}
	if _, ok := byKey["motor.idle_power_thr"]; !ok {
		t.Errorf("motor.idle_power_thr should be recursed from object table; fields=%v", byKey)
	}
	if _, ok := byKey["motor"]; ok {
		t.Errorf("motor object should not be emitted as a leaf; fields=%v", byKey)
	}
}

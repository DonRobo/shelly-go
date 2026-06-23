package gen

import "sort"

// ApplyExceptions patches the parsed Spec for cases the docs express in a way
// the generic parser cannot handle correctly on its own.
//
// Design rule: keep the generic fetch/parse code free of special cases. Every
// per-component or per-field correction lives here, in one place, so it is easy
// to find and iterate on as the docs evolve or the fidelity gate surfaces a
// real difference against the hand-written structs.
func ApplyExceptions(spec *Spec) {
	for _, fix := range fixes {
		fix(spec)
	}
}

// fixes is the ordered list of corrections applied after parsing. Entries are
// added as the generated output is diffed against the existing hand-written
// client or validated against real devices.
var fixes = []func(*Spec){
	// The established client (and the devices) use identifier/RPC casing that
	// differs from the docs page titles. Keep the proven casing so generated
	// types and method strings match what already works.
	setGoPrefix("WiFi", "Wifi"),  // RPC: "Wifi.GetConfig"
	setGoPrefix("Mqtt", "MQTT"),  // RPC: "MQTT.GetConfig"

	// Ui is a real device component (idle screen brightness on devices with a
	// display) but Shelly does not document it — the docs page 404s, so the
	// docs-driven parser never sees it. Inject it from a shape discovered by
	// introspecting a live device (Ui.GetConfig -> {"idle_brightness": 0}).
	// Hardcoding keeps the build hermetic and uncapped by on-hand hardware;
	// drop this fix if Shelly ever documents Ui.
	addComponent(&Component{
		Name:         "Ui",
		HasGetConfig: true,
		HasSetConfig: true,
		Keyed:        false,
		Fields: []*Field{{
			Key:         "idle_brightness",
			Type:        "integer",
			Description: "Brightness of the device's screen while idle, in percent (0-100). Not documented by Shelly; sourced via device introspection.",
		}},
	}),

	// A few status fields are documented as a bare "object" with no sub-table, so
	// the parser can only produce json.RawMessage. Inject the known shapes (from
	// the docs prose / a live device) so they generate as typed nested structs.
	// Cover's aenergy mirrors the shape Switch documents in full.
	expandStatusObject("Cover", "aenergy", []*Field{
		{Key: "total", Type: "number", Description: "Total energy consumed in Watt-hours."},
		{Key: "by_minute", Type: "array", Elem: "number", Description: "Energy consumption by minute (in Milliwatt-hours) for the last three minutes."},
		{Key: "minute_ts", Type: "integer", Description: "Unix timestamp of the first second of the last minute."},
	}),
	expandStatusObject("Cover", "temperature", []*Field{
		{Key: "tC", Type: "number", Description: "Temperature in Celsius (null if temperature is out of the measurement range)."},
		{Key: "tF", Type: "number", Description: "Temperature in Fahrenheit (null if temperature is out of the measurement range)."},
	}),
	expandStatusObject("DevicePower", "external", []*Field{
		{Key: "present", Type: "boolean", Description: "Whether an external power source is connected."},
	}),
	// Sys documents wakeup_reason as an object but lists its sub-fields only in
	// prose (no nested table), so the parser drops the row entirely. Inject the
	// documented shape; expandStatusObject appends it when the parent is absent.
	expandStatusObject("Sys", "wakeup_reason", []*Field{
		{Key: "boot", Type: "string", Description: "Boot type, one of: poweron, software_restart, deepsleep_wake, internal, unknown."},
		{Key: "cause", Type: "string", Description: "Boot cause, one of: button, usb, periodic, status_update, alarm, alarm_test, undefined."},
	}),
	// reset_reason and safe_mode are undocumented but appear on real devices;
	// keep the coverage the hand-written SysStatus carried before it moved to
	// codegen. Sourced from observed responses; drop if Shelly documents them.
	addStatusField("Sys", &Field{Key: "reset_reason", Type: "integer", Description: "Numeric reset reason. Not documented by Shelly; appears in firmware responses."}),
	addStatusField("Sys", &Field{Key: "safe_mode", Type: "boolean", Description: "True if the device is operating in Safe Mode; present only in that mode. Not documented by Shelly."}),

	// Cover.swap_inputs is documented, but its Type cell is left blank in the
	// docs, so the parser cannot infer a type and drops the row. It is a boolean
	// (swap the two inputs' open/close functions); inject it with the type the
	// docs omit. Drop this fix if Shelly fills in the Type cell.
	addField("Cover", &Field{Key: "swap_inputs", Type: "boolean", Description: "Defines whether the functions of the two inputs are swapped. Only present if there are two inputs associated with the Cover instance. Documented without a type by Shelly."}),
}

// expandStatusObject returns a fix that represents a status object the docs
// leave unbroken-down with typed sub-fields, so the emitter builds a proper
// nested struct instead of json.RawMessage. Sub-field keys are relative; the
// parent key is prefixed. An opaque "object" leaf is replaced in place (keeping
// field order stable); if the parser dropped the row entirely — as it does when
// the docs describe the children only in prose (Sys.wakeup_reason) — the
// sub-fields are appended instead.
func expandStatusObject(comp, key string, sub []*Field) func(*Spec) {
	return func(s *Spec) {
		c := s.component(comp)
		if c == nil {
			return
		}
		expanded := make([]*Field, 0, len(sub))
		for _, sf := range sub {
			expanded = append(expanded, &Field{
				Key:         key + "." + sf.Key,
				Type:        sf.Type,
				Elem:        sf.Elem,
				Description: sf.Description,
			})
		}
		out := make([]*Field, 0, len(c.StatusFields)+len(sub))
		found := false
		for _, f := range c.StatusFields {
			if f.Key == key && f.Type == "object" {
				out = append(out, expanded...)
				found = true
				continue
			}
			out = append(out, f)
		}
		if !found {
			out = append(out, expanded...)
		}
		c.StatusFields = out
	}
}

// addStatusField returns a fix that appends a status leaf the docs omit but real
// devices return, so it generates as a typed field instead of being silently
// dropped on unmarshal. No-op if the docs already document a field with that key.
func addStatusField(comp string, f *Field) func(*Spec) {
	return func(s *Spec) {
		c := s.component(comp)
		if c == nil {
			return
		}
		for _, ex := range c.StatusFields {
			if ex.Key == f.Key {
				return
			}
		}
		c.StatusFields = append(c.StatusFields, f)
	}
}

// addField returns a fix that appends a config leaf the docs document without a
// usable type (so the parser drops it). No-op if a field with that key already
// parsed, so the fix self-disables once the docs are corrected.
func addField(comp string, f *Field) func(*Spec) {
	return func(s *Spec) {
		c := s.component(comp)
		if c == nil {
			return
		}
		for _, ex := range c.Fields {
			if ex.Key == f.Key {
				return
			}
		}
		c.Fields = append(c.Fields, f)
	}
}

// addComponent returns a fix that appends a component the docs do not express,
// keeping Components sorted by Name. No-op if a component of that name is already
// present (e.g. the docs started documenting it).
func addComponent(c *Component) func(*Spec) {
	return func(s *Spec) {
		if s.component(c.Name) != nil {
			return
		}
		s.Components = append(s.Components, c)
		sort.Slice(s.Components, func(i, j int) bool {
			return s.Components[i].Name < s.Components[j].Name
		})
	}
}

// setGoPrefix returns a fix that overrides a component's identifier/method
// prefix.
func setGoPrefix(docName, prefix string) func(*Spec) {
	return func(s *Spec) {
		if c := s.component(docName); c != nil {
			c.GoPrefix = prefix
		}
	}
}

// component returns the named component, or nil if absent. Helper for fixes.
func (s *Spec) component(name string) *Component {
	for _, c := range s.Components {
		if c.Name == name {
			return c
		}
	}
	return nil
}

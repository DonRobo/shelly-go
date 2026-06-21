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

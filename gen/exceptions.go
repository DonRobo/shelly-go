package gen

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

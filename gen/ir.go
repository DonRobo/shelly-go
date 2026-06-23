// Package gen contains the documentation-driven code generator for shelly-go.
//
// The pipeline has three stages, run as one command (see ./cmd/gen):
//
//	fetch  — download the Shelly Gen2+ API docs HTML into a persisted cache
//	parse  — turn the cached HTML into the machine-readable Spec (the IR)
//	emit   — generate Go from the IR plus a hardcoded exceptions/fixes layer
//
// The IR (Spec) is the single source of truth shared by the library client
// generator and the Terraform provider generator. It is intentionally simple:
// the messy, doc-specific reality is normalised here so downstream generators
// stay dumb.
package gen

// Spec is the parsed, normalised representation of the Shelly component API.
type Spec struct {
	// SourceURL records where the docs were fetched from, for provenance.
	SourceURL string `json:"sourceUrl"`
	// Components is sorted by Name for stable, reviewable diffs.
	Components []*Component `json:"components"`
}

// Component is one Shelly component/service (Switch, Input, Sys, ...).
type Component struct {
	// Name is the documentation name, e.g. "Switch". Used for the cache file
	// and the generated file name.
	Name string `json:"name"`
	// GoPrefix overrides the identifier/RPC-method prefix when the established
	// client uses a different casing than the docs title (set via exceptions,
	// e.g. docs "WiFi" -> "Wifi", docs "Mqtt" -> "MQTT"). Empty means use Name.
	GoPrefix string `json:"goPrefix,omitempty"`
	// HasGetConfig / HasSetConfig record which config RPCs the docs expose.
	// A component needs both to become a Terraform *_config resource.
	HasGetConfig bool `json:"hasGetConfig"`
	HasSetConfig bool `json:"hasSetConfig"`
	// HasGetStatus records whether the docs expose a GetStatus RPC, i.e. whether
	// a typed <Name>Status + <Name>GetStatusRequest can be generated.
	HasGetStatus bool `json:"hasGetStatus"`
	// HasGetDeviceInfo records whether the docs expose a GetDeviceInfo RPC (only
	// the Shelly service does). DeviceInfoFields then holds its response shape.
	HasGetDeviceInfo bool `json:"hasGetDeviceInfo,omitempty"`
	// DeviceInfoFields are the response properties of the GetDeviceInfo RPC, in
	// document order. Empty unless HasGetDeviceInfo.
	DeviceInfoFields []*Field `json:"deviceInfoFields,omitempty"`
	// Keyed is true when component instances are addressed by a numeric id
	// (Switch:0, Input:1, ...). Singletons like Sys/WiFi/Cloud are not keyed.
	Keyed bool `json:"keyed"`
	// Fields are the configuration properties, in document order.
	Fields []*Field `json:"fields,omitempty"`
	// StatusFields are the status properties (the "Status" table), in document
	// order. Same shape as Fields; the emitter turns them into <Name>Status.
	StatusFields []*Field `json:"statusFields,omitempty"`
}

// Prefix is the identifier and RPC-method prefix for the component.
func (c *Component) Prefix() string {
	if c.GoPrefix != "" {
		return c.GoPrefix
	}
	return c.Name
}

// Field is one configuration property from a component's "configuration" table.
type Field struct {
	// Key is the JSON key. It may be dotted for nested objects (e.g.
	// "counts.enable"); the emitter is responsible for building nested structs.
	Key string `json:"key"`
	// Type is the normalised base type: string, number, integer, boolean,
	// object, array, or unknown.
	Type string `json:"type"`
	// Elem is the normalised element base type for Type == "array" when the docs
	// specify it ("array of type number" -> "number"). Empty when unknown; the
	// emitter then falls back to json.RawMessage for the array.
	Elem string `json:"elem,omitempty"`
	// Nullable is true when the docs describe the type as "<t> or null".
	Nullable bool `json:"nullable,omitempty"`
	// Enum holds the allowed values when the description specifies a
	// "Range of values: ...". Empty otherwise.
	Enum []string `json:"enum,omitempty"`
	// Min and Max hold the documented numeric bounds when the description states
	// an accepted range (e.g. "[0.5..30]"). Both are nil unless a bound parsed.
	// Consumed by the provider to emit plan-time range validators; the library
	// itself does not enforce them.
	Min *float64 `json:"min,omitempty"`
	Max *float64 `json:"max,omitempty"`
	// Description is the cleaned doc text for the property.
	Description string `json:"description,omitempty"`
}

package gen

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// EmitClient writes one <comp>_gen.go file per component into outDir. Each file
// defines the component's typed Config and/or Status structs (with nested
// structs for dotted keys), plus GetConfig/SetConfig/GetStatus request wrappers
// in the same shape as the hand-written client.
//
// Config and status are skipped independently: a component whose <Name>Config or
// <Name>Status type already exists in hand-written (non-generated) source keeps
// that part hand-written, so generation only fills gaps and never collides.
// Deleting a hand-written type later causes generation to take it over — the
// migration path to a fully generated client.
func EmitClient(spec *Spec, componentsDir, shellyDir string) error {
	if err := os.MkdirAll(componentsDir, 0o755); err != nil {
		return err
	}
	// Remove previously generated files first, so a change to the skip-set
	// never leaves a stale *_gen.go behind.
	if err := removeGenerated(componentsDir); err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(shellyDir, "shelly_gen.go")); err != nil && !os.IsNotExist(err) {
		return err
	}
	handConfig, handStatus, err := handWrittenComponents(componentsDir)
	if err != nil {
		return err
	}
	for _, c := range spec.Components {
		if c.Name == "Shelly" {
			continue // the Shelly service is the aggregate, emitted below — not a component
		}
		key := strings.ToLower(c.Name)
		genConfig := (c.HasGetConfig || c.HasSetConfig) && !handConfig[key]
		genStatus := c.HasGetStatus && !handStatus[key]
		if !genConfig && !genStatus {
			continue // nothing to generate, or both parts hand-written
		}
		src, err := emitComponent(c, genConfig, genStatus)
		if err != nil {
			return fmt.Errorf("emit %s: %w", c.Name, err)
		}
		dst := filepath.Join(componentsDir, key+"_gen.go")
		if err := os.WriteFile(dst, src, 0o644); err != nil { //nolint:gosec // generated source.
			return err
		}
	}

	// The Shelly.GetConfig/GetStatus aggregates (one field per component) are
	// generated from the full component list — complete and uncapped, unlike a
	// hand-curated list. They live in the root shelly package and reference the
	// component types. Gap-filled like the rest: skipped while a hand-written
	// aggregate still exists.
	if !handAggregate(shellyDir) {
		agg, err := emitAggregates(spec)
		if err != nil {
			return fmt.Errorf("emit aggregate: %w", err)
		}
		if err := os.WriteFile(filepath.Join(shellyDir, "shelly_gen.go"), agg, 0o644); err != nil { //nolint:gosec // generated source.
			return err
		}
	}
	return nil
}

// handAggregate reports whether a hand-written ShellyGet*Response aggregate
// still exists in the shelly package dir (so generation does not collide).
func handAggregate(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_gen.go") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		if strings.Contains(string(b), "type ShellyGetStatusResponse struct") ||
			strings.Contains(string(b), "type ShellyGetConfigResponse struct") {
			return true
		}
	}
	return false
}

func removeGenerated(dir string) error {
	matches, err := filepath.Glob(filepath.Join(dir, "*_gen.go"))
	if err != nil {
		return err
	}
	for _, m := range matches {
		if err := os.Remove(m); err != nil {
			return err
		}
	}
	return nil
}

// configIdentifier captures the base component name from a hand-written config
// type (FooConfig) or config request (FooGetConfigRequest / FooSetConfigRequest).
var configIdentifier = regexp.MustCompile(`type\s+(\w+?)(?:Config|GetConfigRequest|SetConfigRequest)\s+struct`)

// statusIdentifier captures the base component name from a hand-written status
// type (FooStatus) or status request (FooGetStatusRequest). Nested sub-status
// types (FooBarStatus) also match and add harmless extra keys that never
// correspond to a real component name.
var statusIdentifier = regexp.MustCompile(`type\s+(\w+?)(?:Status|GetStatusRequest)\s+struct`)

// handWrittenComponents returns two sets of component base names (lower-cased):
// those with a hand-written config type/request, and those with a hand-written
// status type/request. Config and status are tracked separately so each part is
// skipped independently. Lower-casing makes detection robust to casing
// differences between the docs and the hand-written client (WiFi vs Wifi).
func handWrittenComponents(dir string) (config, status map[string]bool, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	config = map[string]bool{}
	status = map[string]bool{}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_gen.go") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, nil, err
		}
		for _, m := range configIdentifier.FindAllSubmatch(b, -1) {
			base := strings.ToLower(string(m[1]))
			if base != "" && base != "set" { // ignore the shared SetConfigResponse helper
				config[base] = true
			}
		}
		for _, m := range statusIdentifier.FindAllSubmatch(b, -1) {
			base := strings.ToLower(string(m[1]))
			if base != "" {
				status[base] = true
			}
		}
	}
	return config, status, nil
}

// emitComponent renders and gofmt-formats a single component's generated file,
// containing the config part, the status part, or both.
func emitComponent(c *Component, genConfig, genStatus bool) ([]byte, error) {
	var b strings.Builder
	b.WriteString("// Code generated by cmd/gen from the Shelly API docs. DO NOT EDIT.\n\n")
	b.WriteString("package components\n\n")
	b.WriteString("import (\n\t\"github.com/DonRobo/shelly-go/rpc\"\n\t\"resty.dev/v3\"\n)\n\n")

	needsJSON := false
	if genConfig {
		// Nested struct types first (depth-first), then the root Config struct.
		emitStructs(&b, c.Prefix()+"Config", buildTree(c.Fields), &needsJSON)
		if c.HasGetConfig {
			emitGetConfig(&b, c)
		}
		if c.HasSetConfig {
			emitSetConfig(&b, c)
		}
	}
	if genStatus {
		emitStructs(&b, c.Prefix()+"Status", buildTree(c.StatusFields), &needsJSON)
		emitGetStatus(&b, c)
	}

	src := b.String()
	if needsJSON {
		src = strings.Replace(src, "import (\n", "import (\n\t\"encoding/json\"\n", 1)
	}

	formatted, err := format.Source([]byte(src))
	if err != nil {
		// Return the unformatted source in the error to aid debugging.
		return nil, fmt.Errorf("%w\n--- source ---\n%s", err, src)
	}
	return formatted, nil
}

// --- field tree (handles dotted nested keys) ---------------------------------

type node struct {
	name     string // path segment (json key fragment)
	field    *Field // non-nil for leaves
	order    []string
	children map[string]*node
}

func newNode(name string) *node {
	return &node{name: name, children: map[string]*node{}}
}

func (n *node) child(name string) *node {
	c, ok := n.children[name]
	if !ok {
		c = newNode(name)
		n.children[name] = c
		n.order = append(n.order, name)
	}
	return c
}

// buildTree groups fields by their dotted key segments into a tree.
func buildTree(fields []*Field) *node {
	root := newNode("")
	for _, f := range fields {
		segs := strings.Split(f.Key, ".")
		cur := root
		for i, s := range segs {
			cur = cur.child(s)
			if i == len(segs)-1 {
				cur.field = f
			}
		}
	}
	return root
}

// emitStructs writes nested struct types (post-order) and then the named struct
// for the given node's children.
func emitStructs(b *strings.Builder, typeName string, n *node, needsJSON *bool) {
	// First emit child object types, so they are declared before use.
	for _, name := range n.order {
		c := n.children[name]
		if len(c.children) > 0 {
			emitStructs(b, typeName+goName(name), c, needsJSON)
		}
	}

	fmt.Fprintf(b, "// %s is generated from the Shelly API documentation.\n", typeName)
	fmt.Fprintf(b, "type %s struct {\n", typeName)
	for _, name := range n.order {
		c := n.children[name]
		field := goName(name)
		if len(c.children) > 0 {
			// Nested object.
			fmt.Fprintf(b, "\t%s *%s `json:\"%s,omitempty\"`\n", field, typeName+field, name)
			continue
		}
		writeComment(b, field, c.field.Description)
		goType := goLeafType(name, c.field, needsJSON)
		fmt.Fprintf(b, "\t%s %s `json:\"%s\"`\n\n", field, goType, jsonTag(name, c.field))
	}
	b.WriteString("}\n\n")
}

// goLeafType maps an IR leaf field to a Go type. A keyed "id" field is a bare
// int (matching the hand-written client); everything else is a pointer so unset
// config values are omitted. The Shelly.GetDeviceInfo response also has an "id",
// but it is a string device identifier, so the documented type wins there.
func goLeafType(key string, f *Field, needsJSON *bool) string {
	if key == "id" {
		if f.Type == "string" {
			return "string"
		}
		return "int"
	}
	switch f.Type {
	case "string":
		return "*string"
	case "number":
		return "*float64"
	case "integer":
		return "*int"
	case "boolean":
		return "*bool"
	case "array":
		// A slice of scalars when the docs name the element type; nil-able as-is.
		switch f.Elem {
		case "string":
			return "[]string"
		case "number":
			return "[]float64"
		case "integer":
			return "[]int"
		case "boolean":
			return "[]bool"
		default:
			*needsJSON = true
			return "json.RawMessage"
		}
	default: // object, unknown
		*needsJSON = true
		return "json.RawMessage"
	}
}

func jsonTag(key string, f *Field) string {
	if key == "id" {
		return key
	}
	return key + ",omitempty"
}

// --- request wrappers --------------------------------------------------------

func emitGetConfig(b *strings.Builder, c *Component) {
	req := c.Prefix() + "GetConfigRequest"
	cfg := c.Prefix() + "Config"
	fmt.Fprintf(b, "// %s requests the configuration of the %s component.\n", req, c.Name)
	fmt.Fprintf(b, "type %s struct {\n", req)
	if c.Keyed {
		fmt.Fprintf(b, "\t// ID of the %s component instance.\n\tID int `json:\"id\"`\n", c.Name)
	}
	b.WriteString("}\n\n")
	fmt.Fprintf(b, "func (r *%s) Method() string { return %q }\n\n", req, c.Prefix()+".GetConfig")
	fmt.Fprintf(b, "func (r *%s) NewTypedResponse() *%s { return &%s{} }\n\n", req, cfg, cfg)
	fmt.Fprintf(b, "func (r *%s) NewResponse() any { return r.NewTypedResponse() }\n\n", req)
	fmt.Fprintf(b, "func (r *%s) Do(client *resty.Client) (*%s, *rpc.Frame, error) {\n", req, cfg)
	b.WriteString("\tresp := r.NewTypedResponse()\n\traw, err := rpc.Do(client, r, resp)\n\treturn resp, raw, err\n}\n\n")
}

func emitSetConfig(b *strings.Builder, c *Component) {
	req := c.Prefix() + "SetConfigRequest"
	cfg := c.Prefix() + "Config"
	fmt.Fprintf(b, "// %s updates the configuration of the %s component.\n", req, c.Name)
	fmt.Fprintf(b, "type %s struct {\n", req)
	if c.Keyed {
		fmt.Fprintf(b, "\t// ID of the %s component instance.\n\tID int `json:\"id\"`\n\n", c.Name)
	}
	fmt.Fprintf(b, "\t// Config that the method takes.\n\tConfig %s `json:\"config\"`\n", cfg)
	b.WriteString("}\n\n")
	fmt.Fprintf(b, "func (r *%s) Method() string { return %q }\n\n", req, c.Prefix()+".SetConfig")
	fmt.Fprintf(b, "func (r *%s) NewTypedResponse() *rpc.SetConfigResponse { return &rpc.SetConfigResponse{} }\n\n", req)
	fmt.Fprintf(b, "func (r *%s) NewResponse() any { return r.NewTypedResponse() }\n\n", req)
	fmt.Fprintf(b, "func (r *%s) Do(client *resty.Client) (*rpc.SetConfigResponse, *rpc.Frame, error) {\n", req)
	b.WriteString("\tresp := r.NewTypedResponse()\n\traw, err := rpc.Do(client, r, resp)\n\treturn resp, raw, err\n}\n\n")
}

func emitGetStatus(b *strings.Builder, c *Component) {
	req := c.Prefix() + "GetStatusRequest"
	st := c.Prefix() + "Status"
	fmt.Fprintf(b, "// %s requests the status of the %s component.\n", req, c.Name)
	fmt.Fprintf(b, "type %s struct {\n", req)
	if c.Keyed {
		fmt.Fprintf(b, "\t// ID of the %s component instance.\n\tID int `json:\"id\"`\n", c.Name)
	}
	b.WriteString("}\n\n")
	fmt.Fprintf(b, "func (r *%s) Method() string { return %q }\n\n", req, c.Prefix()+".GetStatus")
	fmt.Fprintf(b, "func (r *%s) NewTypedResponse() *%s { return &%s{} }\n\n", req, st, st)
	fmt.Fprintf(b, "func (r *%s) NewResponse() any { return r.NewTypedResponse() }\n\n", req)
	fmt.Fprintf(b, "func (r *%s) Do(client *resty.Client) (*%s, *rpc.Frame, error) {\n", req, st)
	b.WriteString("\tresp := r.NewTypedResponse()\n\traw, err := rpc.Do(client, r, resp)\n\treturn resp, raw, err\n}\n\n")
}

// --- aggregate (Shelly.GetConfig / Shelly.GetStatus) -------------------------

// emitAggregates renders shelly_gen.go: the ShellyGetConfigResponse and
// ShellyGetStatusResponse aggregates plus their request wrappers.
func emitAggregates(spec *Spec) ([]byte, error) {
	var b strings.Builder
	b.WriteString("// Code generated by cmd/gen from the Shelly API docs. DO NOT EDIT.\n\n")
	b.WriteString("package shelly\n\n")
	b.WriteString("import (\n\t\"encoding/json\"\n\t\"fmt\"\n\n\t\"github.com/DonRobo/shelly-go/components\"\n\t\"github.com/DonRobo/shelly-go/rpc\"\n\t\"resty.dev/v3\"\n)\n\n")

	emitAggregate(&b, spec, "Status")
	emitAggregate(&b, spec, "Config")

	// The Shelly service also exposes GetDeviceInfo, a one-off method whose flat
	// response the docs document directly. Emit it alongside the aggregates so the
	// whole Shelly service is generated.
	for _, c := range spec.Components {
		if c.Name == "Shelly" && c.HasGetDeviceInfo {
			emitDeviceInfo(&b, c)
		}
	}

	formatted, err := format.Source([]byte(b.String()))
	if err != nil {
		return nil, fmt.Errorf("%w\n--- source ---\n%s", err, b.String())
	}
	return formatted, nil
}

// emitAggregate writes one aggregate response (kind is "Config" or "Status"):
// a struct with a field per component — a pointer for singletons, a slice for
// keyed components — plus a custom UnmarshalJSON that routes the device's
// component keys ("sys", "switch:0", ...) to the typed fields, and the request
// wrapper. The Shelly service itself is excluded (it is this aggregate).
func emitAggregate(b *strings.Builder, spec *Spec, kind string) {
	resp := "ShellyGet" + kind + "Response"
	req := "ShellyGet" + kind + "Request"
	method := "Shelly.Get" + kind

	var singles, keyed []*Component
	for _, c := range spec.Components {
		if c.Name == "Shelly" {
			continue
		}
		has := c.HasGetStatus
		if kind == "Config" {
			has = c.HasGetConfig
		}
		if !has {
			continue
		}
		if c.Keyed {
			keyed = append(keyed, c)
		} else {
			singles = append(singles, c)
		}
	}

	fmt.Fprintf(b, "// %s is the aggregate %s of every component, returned by %s.\n",
		resp, strings.ToLower(kind), method)
	fmt.Fprintf(b, "type %s struct {\n", resp)
	for _, c := range singles {
		fmt.Fprintf(b, "\t%s *components.%s%s `json:\"%s,omitempty\"`\n", c.Prefix(), c.Prefix(), kind, strings.ToLower(c.Name))
	}
	for _, c := range keyed {
		fmt.Fprintf(b, "\t%s []*components.%s%s `json:\"%s,omitempty\"`\n",
			pluralize(c.Prefix()), c.Prefix(), kind, strings.ToLower(pluralize(c.Name)))
	}
	b.WriteString("}\n\n")

	fmt.Fprintf(b, "func (r *%s) UnmarshalJSON(b []byte) error {\n", resp)
	b.WriteString("\traw := make(map[string]json.RawMessage)\n")
	b.WriteString("\tif err := json.Unmarshal(b, &raw); err != nil {\n\t\treturn err\n\t}\n")
	for _, c := range singles {
		fmt.Fprintf(b, "\tif v, ok := raw[%q]; ok {\n\t\tvar s components.%s%s\n\t\tif err := json.Unmarshal(v, &s); err != nil {\n\t\t\treturn err\n\t\t}\n\t\tr.%s = &s\n\t}\n",
			strings.ToLower(c.Name), c.Prefix(), kind, c.Prefix())
	}
	for _, c := range keyed {
		field := pluralize(c.Prefix())
		fmt.Fprintf(b, "\tfor i := 0; ; i++ {\n\t\tv, ok := raw[fmt.Sprintf(%q, i)]\n\t\tif !ok {\n\t\t\tbreak\n\t\t}\n\t\tvar s components.%s%s\n\t\tif err := json.Unmarshal(v, &s); err != nil {\n\t\t\treturn err\n\t\t}\n\t\tr.%s = append(r.%s, &s)\n\t}\n",
			strings.ToLower(c.Name)+":%d", c.Prefix(), kind, field, field)
	}
	b.WriteString("\treturn nil\n}\n\n")

	fmt.Fprintf(b, "// %s requests the aggregate %s of every component.\n", req, strings.ToLower(kind))
	fmt.Fprintf(b, "type %s struct{}\n\n", req)
	fmt.Fprintf(b, "func (r *%s) Method() string { return %q }\n\n", req, method)
	fmt.Fprintf(b, "func (r *%s) NewTypedResponse() *%s { return &%s{} }\n\n", req, resp, resp)
	fmt.Fprintf(b, "func (r *%s) NewResponse() any { return r.NewTypedResponse() }\n\n", req)
	fmt.Fprintf(b, "func (r *%s) Do(client *resty.Client) (*%s, *rpc.Frame, error) {\n", req, resp)
	b.WriteString("\tresp := r.NewTypedResponse()\n\traw, err := rpc.Do(client, r, resp)\n\treturn resp, raw, err\n}\n\n")
}

// emitDeviceInfo renders the ShellyGetDeviceInfoResponse struct (from the docs'
// flat response table) and its request wrapper. The request carries the single
// documented `ident` parameter.
func emitDeviceInfo(b *strings.Builder, c *Component) {
	resp := "ShellyGetDeviceInfoResponse"
	req := "ShellyGetDeviceInfoRequest"

	needsJSON := false
	emitStructs(b, resp, buildTree(c.DeviceInfoFields), &needsJSON)

	fmt.Fprintf(b, "// %s requests static device information via Shelly.GetDeviceInfo.\n", req)
	fmt.Fprintf(b, "type %s struct {\n", req)
	b.WriteString("\t// Ident includes extra identifying fields (key, batch, fw_sbits) when true.\n")
	b.WriteString("\tIdent *bool `json:\"ident,omitempty\"`\n")
	b.WriteString("}\n\n")
	fmt.Fprintf(b, "func (r *%s) Method() string { return %q }\n\n", req, "Shelly.GetDeviceInfo")
	fmt.Fprintf(b, "func (r *%s) NewTypedResponse() *%s { return &%s{} }\n\n", req, resp, resp)
	fmt.Fprintf(b, "func (r *%s) NewResponse() any { return r.NewTypedResponse() }\n\n", req)
	fmt.Fprintf(b, "func (r *%s) Do(client *resty.Client) (*%s, *rpc.Frame, error) {\n", req, resp)
	b.WriteString("\tresp := r.NewTypedResponse()\n\traw, err := rpc.Do(client, r, resp)\n\treturn resp, raw, err\n}\n\n")
}

// pluralize returns the plural of a Go identifier for aggregate slice fields
// (Switch -> Switches, Humidity -> Humidities, Cover -> Covers).
func pluralize(s string) string {
	switch {
	case strings.HasSuffix(s, "s"), strings.HasSuffix(s, "x"), strings.HasSuffix(s, "z"),
		strings.HasSuffix(s, "ch"), strings.HasSuffix(s, "sh"):
		return s + "es"
	case len(s) >= 2 && strings.HasSuffix(s, "y") && !isVowel(s[len(s)-2]):
		return s[:len(s)-1] + "ies"
	default:
		return s + "s"
	}
}

func isVowel(b byte) bool {
	switch b {
	case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
		return true
	}
	return false
}

// --- naming + comments -------------------------------------------------------

// initialisms that should be fully upper-cased (or specially cased) in Go names.
var initialisms = map[string]string{
	"id": "ID", "ip": "IP", "mac": "MAC", "ssid": "SSID", "bssid": "BSSID",
	"ap": "AP", "ble": "BLE", "mqtt": "MQTT", "http": "HTTP", "https": "HTTPS",
	"url": "URL", "uri": "URI", "tls": "TLS", "ssl": "SSL", "ca": "CA",
	"rpc": "RPC", "ui": "UI", "em": "EM", "pm": "PM", "dali": "DALI",
	"rgb": "RGB", "rgbw": "RGBW", "cct": "CCT", "utc": "UTC", "dns": "DNS",
	"ntp": "NTP", "json": "JSON", "xml": "XML", "gps": "GPS", "cpu": "CPU",
	"uuid": "UUID", "wifi": "WiFi", "eth": "Eth", "fw": "FW", "hw": "HW",
}

// GoName converts a snake_case JSON key to the exported Go identifier used in
// the generated client. Exported so the Terraform provider generator derives
// the same field names without duplicating the rules.
func GoName(key string) string { return goName(key) }

// goName converts a snake_case json key to an exported Go identifier.
func goName(key string) string {
	parts := strings.Split(key, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		if init, ok := initialisms[strings.ToLower(p)]; ok {
			parts[i] = init
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

// writeComment emits a doc comment, word-wrapped to a sensible width.
func writeComment(b *strings.Builder, name, desc string) {
	desc = strings.TrimSpace(desc)
	if desc == "" {
		fmt.Fprintf(b, "\t// %s is generated from the Shelly API documentation.\n", name)
		return
	}
	words := strings.Fields(name + " " + lowerFirstSentence(desc))
	line := "\t//"
	for _, w := range words {
		if len(line)+1+len(w) > 84 {
			b.WriteString(line + "\n")
			line = "\t//"
		}
		line += " " + w
	}
	b.WriteString(line + "\n")
}

// lowerFirstSentence makes the description read naturally after the field name
// ("Name name of the switch" -> "Name name of the switch instance").
func lowerFirstSentence(desc string) string {
	if desc == "" {
		return desc
	}
	return strings.ToLower(desc[:1]) + desc[1:]
}

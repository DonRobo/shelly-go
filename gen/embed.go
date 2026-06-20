package gen

import (
	_ "embed"
	"encoding/json"
)

//go:embed spec.json
var specJSON []byte

// Load returns the committed Spec parsed from the embedded spec.json.
//
// The Terraform provider generator uses this to drive resource generation
// without re-fetching or re-parsing the documentation: the IR is the shared,
// versioned contract between the two repos.
func Load() (*Spec, error) {
	var s Spec
	if err := json.Unmarshal(specJSON, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

package shelly

// Regenerate the typed config client from the Shelly API documentation.
//
// The generator fetches the docs into a local cache under gen/cache (gitignored;
// only hitting the network for pages it doesn't already have), parses them into
// the committed gen/spec.json, and emits <component>_config_gen.go files for
// every component that isn't already implemented by hand.
//
//go:generate go run ./cmd/gen

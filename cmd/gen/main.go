// Command gen runs the documentation-driven code generator: fetch (cached) ->
// parse -> emit. Run it via `go generate ./...` or directly:
//
//	go run ./cmd/gen            # use committed cache, regenerate IR
//	go run ./cmd/gen -refresh   # re-download the docs first
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/DonRobo/shelly-go/gen"
)

func main() {
	cacheDir := flag.String("cache", "gen/cache", "directory for persisted docs HTML")
	specOut := flag.String("spec", "gen/spec.json", "output path for the parsed IR")
	clientOut := flag.String("client", ".", "directory to write generated client code into")
	refresh := flag.Bool("refresh", false, "re-download docs even if cached")
	flag.Parse()

	names, err := gen.Fetch(*cacheDir, *refresh)
	check(err)
	fmt.Fprintf(os.Stderr, "fetched %d component pages\n", len(names))

	spec, err := gen.Parse(*cacheDir, names)
	check(err)
	gen.ApplyExceptions(spec)

	b, err := json.MarshalIndent(spec, "", "  ")
	check(err)
	check(os.WriteFile(*specOut, append(b, '\n'), 0o644))
	fmt.Fprintf(os.Stderr, "wrote %s (%d components)\n", *specOut, len(spec.Components))

	check(gen.EmitClient(spec, *clientOut))
	fmt.Fprintf(os.Stderr, "wrote generated client to %s\n", *clientOut)
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

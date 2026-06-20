package gen

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

// DocsBase is the root of the Shelly Gen2+ Components & Services documentation.
const DocsBase = "https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/"

// indexFile is the cache filename for the Components & Services index page.
const indexFile = "_index.html"

// componentHref matches nav links like href="/gen2/ComponentsAndServices/Switch"
// and captures the component name. The index uses no trailing slash on these.
var componentHref = regexp.MustCompile(`href="/gen2/ComponentsAndServices/([A-Za-z0-9]+)"`)

// Fetch ensures the index and every component page exists under cacheDir,
// downloading only what is missing (or everything, when refresh is true). It
// returns the sorted list of component names discovered in the index.
//
// The cache (gen/cache) is gitignored: it is re-fetched on demand and exists
// only to avoid re-downloading on every run. The committed source of truth is
// gen/spec.json, parsed from this cache.
func Fetch(cacheDir string, refresh bool) ([]string, error) {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, err
	}

	indexHTML, err := cachedGet(filepath.Join(cacheDir, indexFile), DocsBase, refresh)
	if err != nil {
		return nil, fmt.Errorf("fetch index: %w", err)
	}

	names := componentNames(indexHTML)
	for _, name := range names {
		dst := filepath.Join(cacheDir, name+".html")
		if _, err := cachedGet(dst, DocsBase+name+"/", refresh); err != nil {
			return nil, fmt.Errorf("fetch %s: %w", name, err)
		}
	}
	return names, nil
}

// componentNames extracts the unique, sorted component names from the index nav.
func componentNames(indexHTML []byte) []string {
	seen := map[string]bool{}
	for _, m := range componentHref.FindAllSubmatch(indexHTML, -1) {
		seen[string(m[1])] = true
	}
	out := make([]string, 0, len(seen))
	for n := range seen {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// cachedGet returns the cached file when present (and refresh is false),
// otherwise downloads url, writes it to path, and returns the bytes.
func cachedGet(path, url string, refresh bool) ([]byte, error) {
	if !refresh {
		if b, err := os.ReadFile(path); err == nil {
			return b, nil
		}
	}

	// Polite pause so a cold cache fill doesn't hammer the docs server.
	time.Sleep(250 * time.Millisecond)

	resp, err := http.Get(url) //nolint:gosec // URL is a fixed docs host + parsed component name.
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: unexpected status %s", url, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, body, 0o644); err != nil { //nolint:gosec // docs cache, not sensitive.
		return nil, err
	}
	return body, nil
}

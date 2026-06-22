package gen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

// Parse reads every <name>.html from cacheDir and builds the Spec. names is the
// component list returned by Fetch.
func Parse(cacheDir string, names []string) (*Spec, error) {
	spec := &Spec{SourceURL: DocsBase}
	for _, name := range names {
		b, err := os.ReadFile(filepath.Join(cacheDir, name+".html"))
		if err != nil {
			return nil, fmt.Errorf("read cache for %s: %w", name, err)
		}
		comp, err := parsePage(name, b)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", name, err)
		}
		spec.Components = append(spec.Components, comp)
	}
	sort.Slice(spec.Components, func(i, j int) bool {
		return spec.Components[i].Name < spec.Components[j].Name
	})
	return spec, nil
}

// parsePage extracts a single component's methods and config fields from its
// documentation HTML.
func parsePage(name string, htmlBytes []byte) (*Component, error) {
	doc, err := html.Parse(bytes.NewReader(htmlBytes))
	if err != nil {
		return nil, err
	}

	// Pre-order DFS gives nodes in document order, which lets us pair a heading
	// with the table that follows it.
	var nodes []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		nodes = append(nodes, n)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	comp := &Component{Name: name}

	// Method detection: the docs render headings with ids like "switchgetconfig".
	lname := strings.ToLower(name)
	for _, n := range nodes {
		if n.Type != html.ElementNode {
			continue
		}
		switch attrVal(n, "id") {
		case lname + "getconfig":
			comp.HasGetConfig = true
		case lname + "setconfig":
			comp.HasSetConfig = true
		case lname + "getstatus":
			comp.HasGetStatus = true
		}
	}

	// Config fields live in the table following the heading id="configuration".
	if t := tableAfterHeading(nodes, "configuration"); t != nil {
		comp.Fields = parseRows(t, "")
	}
	// Status fields live in the table following the heading id="status".
	if t := tableAfterHeading(nodes, "status"); t != nil {
		comp.StatusFields = parseRows(t, "")
	}
	comp.Keyed = inferKeyed(comp)

	return comp, nil
}

// tableAfterHeading finds the heading with the given id and returns the first
// <table> element that appears after it in document order, or nil.
func tableAfterHeading(nodes []*html.Node, headingID string) *html.Node {
	start := -1
	for i, n := range nodes {
		if n.Type == html.ElementNode && isHeading(n) && attrVal(n, "id") == headingID {
			start = i
			break
		}
	}
	if start < 0 {
		return nil
	}
	for i := start + 1; i < len(nodes); i++ {
		if nodes[i].Type == html.ElementNode && nodes[i].Data == "table" {
			return nodes[i]
		}
	}
	return nil
}

// validKey matches a real Shelly config property key: lowercase, dotted for
// nesting (e.g. "in_mode", "counts.power_thr"). This rejects continuation rows
// some tables use to explain enum values, whose first cell is prose.
var validKey = regexp.MustCompile(`^[a-z][a-z0-9_]*(\.[a-z0-9_]+)*$`)

// parseRows turns a Property|Type|Description table into Fields. prefix is the
// dotted path of the enclosing object ("" at the top level).
//
// Object properties document their children in a nested <table> inside the
// description cell; those are parsed recursively into dotted keys (ap.ssid,
// device.name), which the emitter turns into nested structs.
//
// Some component tables (e.g. Cover) use a wide, multi-column layout where the
// possible values of an enum property are rendered as extra rows. Those rows
// have either an invalid key (contains spaces) or a prose "type" cell, so we
// skip any row that doesn't have a valid key AND a recognised type.
func parseRows(table *html.Node, prefix string) []*Field {
	var fields []*Field
	for _, tr := range directRows(table) {
		cells := directCells(tr)
		if len(cells) < 2 {
			continue // header row (uses <th>) or malformed
		}
		key := strings.TrimSpace(textOf(cells[0]))
		if !validKey.MatchString(key) {
			continue // continuation/explanation row, not a property
		}
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		// An object property with a nested table documents its children; recurse
		// and don't emit the parent as a leaf.
		if len(cells) >= 3 {
			if nested := firstDescendantTable(cells[2]); nested != nil {
				fields = append(fields, parseRows(nested, fullKey)...)
				continue
			}
		}

		base, nullable := normaliseType(textOf(cells[1]))
		if base == "unknown" {
			continue // prose in the type cell -> not a real property row
		}
		desc := ""
		if len(cells) >= 3 {
			desc = cleanText(textExcludingTables(cells[2]))
		}
		fields = append(fields, &Field{
			Key:         fullKey,
			Type:        base,
			Nullable:    nullable,
			Enum:        enumValues(desc),
			Description: desc,
		})
	}
	return fields
}

var rangeOfValues = regexp.MustCompile(`(?i)Range of values?:\s*(.+)`)

// enumValues pulls a "Range of values: a, b, c" list out of a description.
// It is deliberately conservative: real enum values are single tokens, so
// multi-word trailing prose is dropped (and refined further in exceptions.go).
func enumValues(desc string) []string {
	m := rangeOfValues.FindStringSubmatch(desc)
	if m == nil {
		return nil
	}
	// Stop at the first sentence end or parenthetical qualifier.
	tail := m[1]
	if i := strings.IndexAny(tail, ".("); i >= 0 {
		tail = tail[:i]
	}
	var out []string
	for _, tok := range strings.Split(tail, ",") {
		tok = strings.TrimSpace(tok)
		if tok == "" || strings.ContainsAny(tok, " \t") {
			continue // not a bare enum token
		}
		out = append(out, tok)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// normaliseType maps a doc type string to a base type and a nullable flag.
func normaliseType(typeStr string) (base string, nullable bool) {
	s := strings.ToLower(strings.TrimSpace(typeStr))
	nullable = strings.Contains(s, "null")
	switch {
	case strings.Contains(s, "array"):
		return "array", nullable
	case strings.Contains(s, "object"):
		return "object", nullable
	case strings.Contains(s, "bool"):
		return "boolean", nullable
	case strings.Contains(s, "integer"):
		return "integer", nullable
	case strings.Contains(s, "number") || strings.Contains(s, "float"):
		return "number", nullable
	case strings.Contains(s, "string"):
		return "string", nullable
	default:
		return "unknown", nullable
	}
}

// inferKeyed reports whether instances are addressed by a numeric id. The docs
// describe such components with an "id" property in their config table.
func inferKeyed(comp *Component) bool {
	for _, f := range comp.Fields {
		if f.Key == "id" {
			return true
		}
	}
	return false
}

// --- small HTML helpers ------------------------------------------------------

func isHeading(n *html.Node) bool {
	switch n.Data {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return true
	}
	return false
}

func attrVal(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// directRows returns the data <tr> rows of a table: the direct children of its
// <tbody> (falling back to the table itself). Rows of nested tables are not
// included, so nesting is handled explicitly via recursion.
func directRows(table *html.Node) []*html.Node {
	scope := table
	for c := table.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "tbody" {
			scope = c
			break
		}
	}
	var rows []*html.Node
	for c := scope.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "tr" {
			rows = append(rows, c)
		}
	}
	return rows
}

// directCells returns the direct <td> children of a row.
func directCells(tr *html.Node) []*html.Node {
	var cells []*html.Node
	for c := tr.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "td" {
			cells = append(cells, c)
		}
	}
	return cells
}

// firstDescendantTable returns the outermost <table> nested inside n, or nil.
func firstDescendantTable(n *html.Node) *html.Node {
	tables := descendants(n, "table")
	if len(tables) == 0 {
		return nil
	}
	return tables[0]
}

// textExcludingTables returns the text of n, skipping any nested <table>
// subtrees (so an object property's description doesn't absorb its children).
func textExcludingTables(n *html.Node) string {
	var sb strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "table" {
			return
		}
		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return sb.String()
}

// descendants returns all element descendants of n with the given tag, in
// document order.
func descendants(n *html.Node, tag string) []*html.Node {
	var out []*html.Node
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == tag {
				out = append(out, c)
			}
			walk(c)
		}
	}
	walk(n)
	return out
}

// textOf returns the concatenated text content of a node.
func textOf(n *html.Node) string {
	var sb strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return sb.String()
}

var wsRun = regexp.MustCompile(`\s+`)

// cleanText collapses whitespace and decodes the few entities the docs emit.
func cleanText(s string) string {
	s = wsRun.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

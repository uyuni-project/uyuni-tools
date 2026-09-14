// SPDX-FileCopyrightText: 2026 Jayprakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"path"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestParseCustomColumns(t *testing.T) {
	spec := "ID:.id,NAME:.name,Nested Field:.status.state"
	expected := []ColumnDef{
		{Header: "ID", Field: "id"},
		{Header: "NAME", Field: "name"},
		{Header: "Nested Field", Field: "status.state"},
	}

	result := parseCustomColumns(spec)
	testutils.AssertEquals(t, "wrong parsed columns", expected, result)
}

type status struct {
	State string `json:"state"`
}

type testItem struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status status `json:"status"`
}

// The items and columns the PrintOutput tests are all rendering.
var (
	printerItems = []testItem{
		{ID: 1, Name: "server-01", Status: status{State: "running"}},
		{ID: 2, Name: "server-02", Status: status{State: "stopped"}},
	}

	printerColumns = []ColumnDef{
		{Header: "ID", Field: "ID"},
		{Header: "NAME", Field: "Name"},
		{Header: "STATUS", Field: "Status.State"},
	}
)

// The outputs expected by more than one format, since a format and its file variant render the same.
var (
	expectedTable = `ID    NAME         STATUS
1     server-01    running
2     server-02    stopped
`

	expectedIDNameTable = `ID    NAME
1     server-01
2     server-02
`

	expectedNames = "server-01 server-02 "
)

// namesTemplate is a Go template rendering the name of every item.
const namesTemplate = `{{range .items}}{{.name}} {{end}}`

// assertPrintOutput fails the test if the format does not render the printerItems as expected.
func assertPrintOutput(t *testing.T, format string, expected string) {
	t.Helper()

	var buf bytes.Buffer
	testutils.AssertNoError(t, "unexpected error", PrintOutput(format, printerItems, printerColumns, &buf))
	testutils.AssertEquals(t, "output mismatch for "+format, expected, buf.String())
}

func TestFieldValue(t *testing.T) {
	item := testItem{
		ID:   123,
		Name: "web-server",
		Status: status{
			State: "running",
		},
	}

	tests := []struct {
		description string
		path        string
		expected    any
		found       bool
	}{
		{"Field with an int value", "ID", 123, true},
		{"Field with a string value", "Name", "web-server", true},
		{"Nested field", "Status.State", "running", true},
		{"Non existent field", "Missing", nil, false},
		{"Non existent nested field", "Status.missing", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result, ok := fieldValue(item, tt.path)
			testutils.AssertEquals(t, "wrong field presence", tt.found, ok)
			testutils.AssertEquals(t, "wrong returned value", tt.expected, result)
		})
	}
}

func TestPrintTable(t *testing.T) {
	tests := []struct {
		description string
		items       []testItem
		expected    string
	}{
		{
			description: "Columns aligned on the header",
			items:       printerItems,
			expected:    expectedTable,
		},
		{
			description: "Columns widened by the longest value",
			items: []testItem{
				{ID: 1, Name: "a-very-long-server-name", Status: status{State: "running"}},
				{ID: 2, Name: "s2", Status: status{State: "stopped"}},
			},
			expected: `ID    NAME                       STATUS
1     a-very-long-server-name    running
2     s2                         stopped
`,
		},
		{
			description: "Headers only when there is no item",
			items:       []testItem{},
			expected: `ID    NAME    STATUS
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			var buf bytes.Buffer
			testutils.AssertNoError(t, "unexpected error", printTable(tt.items, printerColumns, &buf))
			testutils.AssertEquals(t, "printTable output mismatch", tt.expected, buf.String())
		})
	}
}

func TestPrintJSONPath(t *testing.T) {
	items := []testItem{
		{ID: 1, Name: "alpha"},
		{ID: 2, Name: "beta"},
	}

	var buf bytes.Buffer
	testutils.AssertNoError(t, "unexpected error", printJSONPath(items, "{.items[*].name}", &buf))
	testutils.AssertEquals(t, "wrong value", "alpha beta", strings.TrimSpace(buf.String()))
}

func TestPrintJSONPathLargeID(t *testing.T) {
	items := []testItem{{ID: 1000010000, Name: "uyuni.demo.local"}}

	var buf bytes.Buffer
	testutils.AssertNoError(t, "unexpected error", printJSONPath(items, "{.items[*].id}", &buf))
	testutils.AssertEquals(t, "wrong ID returned", "1000010000", strings.TrimSpace(buf.String()))
}

func TestPrintJSONPathEmptyTemplate(t *testing.T) {
	err := printJSONPath([]testItem{}, "  ", &bytes.Buffer{})
	testutils.AssertError(t, "no template given", err)
}

func TestPrintGoTemplate(t *testing.T) {
	items := []testItem{
		{ID: 1, Name: "alpha"},
		{ID: 2, Name: "beta"},
	}

	var buf bytes.Buffer
	tmpl := `{{range .items}}{{.name}} {{end}}`
	testutils.AssertNoError(t, "unexpected error", printGoTemplate(items, tmpl, &buf))
	testutils.AssertEquals(t, "wrong output", "alpha beta", strings.TrimSpace(buf.String()))
}

func TestPrintGoTemplateLargeID(t *testing.T) {
	items := []testItem{{ID: 1000010000, Name: "uyuni.demo.local"}}

	var buf bytes.Buffer
	testutils.AssertNoError(t, "unexpected error", printGoTemplate(items, `{{range .items}}{{.id}}{{end}}`, &buf))
	testutils.AssertEquals(t, "wrong output", "1000010000", strings.TrimSpace(buf.String()))
}

func TestPrintGoTemplateEmptyTemplate(t *testing.T) {
	testutils.AssertError(t, "no template given", printGoTemplate([]testItem{}, "  ", &bytes.Buffer{}))
}

func TestPrintOutputJSON(t *testing.T) {
	assertPrintOutput(t, "json", `[
  {
    "id": 1,
    "name": "server-01",
    "status": {
      "state": "running"
    }
  },
  {
    "id": 2,
    "name": "server-02",
    "status": {
      "state": "stopped"
    }
  }
]
`)
}

func TestPrintOutputYAML(t *testing.T) {
	assertPrintOutput(t, "yaml", `- id: 1
  name: server-01
  status:
    state: running
- id: 2
  name: server-02
  status:
    state: stopped
`)
}

func TestPrintOutputCustomColumns(t *testing.T) {
	assertPrintOutput(t, "custom-columns=ID:ID,NAME:Name", expectedIDNameTable)
}

func TestPrintOutputCustomColumnsFile(t *testing.T) {
	file := path.Join(t.TempDir(), "columns")
	testutils.WriteFile(t, file, "ID:ID\nNAME:Name\n")

	assertPrintOutput(t, "custom-columns-file="+file, expectedIDNameTable)
}

func TestPrintOutputInvalidCustomColumns(t *testing.T) {
	err := PrintOutput("custom-columns=", printerItems, printerColumns, &bytes.Buffer{})

	testutils.AssertError(t, "custom-columns format specified but no valid columns given", err)
}

func TestPrintOutputJSONPath(t *testing.T) {
	assertPrintOutput(t, "jsonpath={.items[*].name}", "server-01 server-02")
}

func TestPrintOutputGoTemplate(t *testing.T) {
	assertPrintOutput(t, "go-template="+namesTemplate, expectedNames)
}

func TestPrintOutputGoTemplateFile(t *testing.T) {
	file := path.Join(t.TempDir(), "template")
	testutils.WriteFile(t, file, namesTemplate)

	assertPrintOutput(t, "go-template-file="+file, expectedNames)
}

func TestPrintOutputTable(t *testing.T) {
	assertPrintOutput(t, "table", expectedTable)
}

func TestPrintOutputUnknownFormat(t *testing.T) {
	err := PrintOutput("not-a-format", printerItems, printerColumns, &bytes.Buffer{})

	testutils.AssertError(t, `unsupported output format "not-a-format"`, err)
}

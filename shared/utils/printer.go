// SPDX-FileCopyrightText: 2026 Jayprakash
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"gopkg.in/yaml.v2"
)

type ColumnDef struct {
	Header string
	Field  string
}

func AddOutputFlag(cmd *cobra.Command, outputFormat *string) {
	cmd.Flags().StringVarP(outputFormat, "output", "o", "table",
		L(`Output format: table|json|yaml|custom-columns=SPEC|custom-columns-file=PATH
SPEC is a comma-separated list of HEADER:JSONPATH pairs (e.g. ID:.id,NAME:.name)`))
}

func PrintOutput(format string, items []map[string]any, cols []ColumnDef, out io.Writer) error {
	switch {
	case format == "json":
		return printJSON(items, out)
	case format == "yaml":
		return printYAML(items, out)
	case strings.HasPrefix(format, "custom-columns="):
		return printColumns(items, strings.TrimPrefix(format, "custom-columns="), out)
	case strings.HasPrefix(format, "custom-columns-file="):
		return printColumnsFromFile(items, strings.TrimPrefix(format, "custom-columns-file="), out)
	default:
		return printTable(items, cols, out)
	}
}

func printColumns(items []map[string]any, spec string, out io.Writer) error {
	parsed := parseCustomColumns(spec)
	if len(parsed) == 0 {
		return fmt.Errorf("custom-columns format specified but no valid columns given")
	}
	return printTable(items, parsed, out)
}

func printColumnsFromFile(items []map[string]any, path string, out io.Writer) error {
	fileCols, err := parseCustomColumnsFile(path)
	if err != nil {
		return err
	}
	return printTable(items, fileCols, out)
}

func printJSON(items []map[string]any, out io.Writer) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, string(data))
	return err
}

func printYAML(items []map[string]any, out io.Writer) error {
	data, err := yaml.Marshal(items)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(out, string(data))
	return err
}

func printTable(items []map[string]any, cols []ColumnDef, out io.Writer) error {
	w := tabwriter.NewWriter(out, 0, 0, 4, ' ', 0)

	for i, col := range cols {
		if i > 0 {
			fmt.Fprint(w, "\t")
		}
		fmt.Fprint(w, col.Header)
	}
	fmt.Fprintln(w)

	for _, item := range items {
		for i, col := range cols {
			if i > 0 {
				fmt.Fprint(w, "\t")
			}
			if val, ok := fieldValue(item, col.Field); ok {
				fmt.Fprint(w, formatValue(val))
			}
		}
		fmt.Fprintln(w)
	}

	return w.Flush()
}

func formatValue(v any) string {
	if f, ok := v.(float64); ok && f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}
	return fmt.Sprintf("%v", v)
}

func parseCustomColumns(spec string) []ColumnDef {
	var cols []ColumnDef

	for _, part := range strings.Split(spec, ",") {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			cols = append(cols, ColumnDef{
				Header: strings.TrimSpace(kv[0]),
				Field:  normalizeFieldPath(kv[1]),
			})
		}
	}
	return cols
}

func parseCustomColumnsFile(path string) ([]ColumnDef, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	spec := strings.TrimSpace(string(content))
	spec = strings.ReplaceAll(spec, "\n", ",")
	spec = strings.ReplaceAll(spec, "\r", "")

	cols := parseCustomColumns(spec)
	if len(cols) == 0 {
		return nil, fmt.Errorf("invalid custom-columns file format")
	}
	return cols, nil
}

func normalizeFieldPath(path string) string {
	return strings.TrimPrefix(strings.TrimSpace(path), ".")
}

func fieldValue(item map[string]any, path string) (any, bool) {
	if path == "" {
		return nil, false
	}

	var current any = item
	for _, part := range strings.Split(normalizeFieldPath(path), ".") {
		if part == "" {
			return nil, false
		}

		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}

		current, ok = m[part]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

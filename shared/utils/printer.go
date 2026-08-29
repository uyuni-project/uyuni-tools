// SPDX-FileCopyrightText: 2026 Jayprakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/util/jsonpath"
)

type ColumnDef struct {
	Header string
	Field  string
}

// jsonpathList mirrors a kubectl List so templates can use .items.
type jsonpathList[T any] struct {
	Items []T `json:"items"`
}

func AddOutputFlag(cmd *cobra.Command, outputFormat *string) {
	cmd.Flags().StringVarP(outputFormat, "output", "o", "table",
		L("Output format. One of: table, json, yaml, custom-columns=SPEC, "+
			"custom-columns-file=PATH, jsonpath=TEMPLATE, go-template=TEMPLATE, go-template-file=PATH. "+
			"SPEC is a comma-separated list of HEADER:FIELD pairs using the Go field names, "+
			"e.g. ID:ID,NAME:Name. TEMPLATE runs against a list of items using the JSON field names, "+
			"e.g. {.items[*].id}"))
}

func PrintOutput[T any](format string, items []T, cols []ColumnDef, out io.Writer) error {
	switch {
	case format == "" || format == "table":
		return printTable(items, cols, out)
	case format == "json":
		return printJSON(items, out)
	case format == "yaml":
		return printYAML(items, out)
	case strings.HasPrefix(format, "custom-columns="):
		return printColumns(items, strings.TrimPrefix(format, "custom-columns="), out)
	case strings.HasPrefix(format, "custom-columns-file="):
		return printColumnsFromFile(items, strings.TrimPrefix(format, "custom-columns-file="), out)
	case strings.HasPrefix(format, "jsonpath="):
		return printJSONPath(items, strings.TrimPrefix(format, "jsonpath="), out)
	case strings.HasPrefix(format, "go-template="):
		return printGoTemplate(items, strings.TrimPrefix(format, "go-template="), out)
	case strings.HasPrefix(format, "go-template-file="):
		return printGoTemplateFromFile(items, strings.TrimPrefix(format, "go-template-file="), out)
	default:
		return fmt.Errorf(L("unsupported output format %q"), format)
	}
}

func printColumns[T any](items []T, spec string, out io.Writer) error {
	parsed := parseCustomColumns(spec)
	if len(parsed) == 0 {
		return errors.New(L("custom-columns format specified but no valid columns given"))
	}
	return printTable(items, parsed, out)
}

func printColumnsFromFile[T any](items []T, path string, out io.Writer) error {
	spec, err := parseCustomColumnsFile(path)
	if err != nil {
		return err
	}
	return printColumns(items, spec, out)
}

// printJSONPath runs a kubectl-style JSONPath template against the item list.
func printJSONPath[T any](items []T, template string, out io.Writer) error {
	if strings.TrimSpace(template) == "" {
		return errors.New(L("jsonpath format specified but no template given"))
	}

	root, err := outputListData(items)
	if err != nil {
		return err
	}

	j := jsonpath.New("output")
	if err := j.Parse(template); err != nil {
		return err
	}
	if err := j.Execute(out, root); err != nil {
		return Errorf(err, L("failed to execute jsonpath %q"), template)
	}
	return nil
}

// printGoTemplate runs a Go text/template against the same items list as jsonpath.
func printGoTemplate[T any](items []T, tmpl string, out io.Writer) error {
	if strings.TrimSpace(tmpl) == "" {
		return errors.New(L("go-template format specified but no template given"))
	}

	root, err := outputListData(items)
	if err != nil {
		return err
	}

	t, err := template.New("output").Parse(tmpl)
	if err != nil {
		return err
	}
	if err := t.Execute(out, root); err != nil {
		return Errorf(err, L("failed to execute go-template %q"), tmpl)
	}
	return nil
}

func printGoTemplateFromFile[T any](items []T, path string, out io.Writer) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return printGoTemplate(items, string(content), out)
}

// outputListData builds {"items":[list]} for jsonpath and go-template.
// UseNumber keeps large IDs exact instead of float scientific notation.
func outputListData[T any](items []T) (any, error) {
	data, err := json.Marshal(jsonpathList[T]{Items: items})
	if err != nil {
		return nil, err
	}

	var root any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&root); err != nil {
		return nil, err
	}
	return root, nil
}

func printJSON[T any](items []T, out io.Writer) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, string(data))
	return err
}

func printYAML[T any](items []T, out io.Writer) error {
	data, err := yaml.Marshal(items)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(out, string(data))
	return err
}

func printTable[T any](items []T, cols []ColumnDef, out io.Writer) error {
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
				fmt.Fprint(w, val)
			}
		}
		fmt.Fprintln(w)
	}

	return w.Flush()
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

func parseCustomColumnsFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	spec := strings.TrimSpace(string(content))
	spec = strings.ReplaceAll(spec, "\n", ",")
	spec = strings.ReplaceAll(spec, "\r", "")
	if spec == "" {
		return "", errors.New(L("invalid custom-columns file format"))
	}
	return spec, nil
}

func normalizeFieldPath(path string) string {
	return strings.TrimPrefix(strings.TrimSpace(path), ".")
}

func fieldValue[T any](item T, path string) (any, bool) {
	parts := strings.Split(normalizeFieldPath(path), ".")
	return resolveFieldPath(reflect.ValueOf(item), parts)
}

func resolveFieldPath(val reflect.Value, parts []string) (any, bool) {
	// Unwrap pointers and interfaces automatically
	if val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return nil, false
		}
		return resolveFieldPath(val.Elem(), parts)
	}

	// Base case: All path segments consumed
	if len(parts) == 0 {
		if !val.IsValid() {
			return nil, false
		}
		return val.Interface(), true
	}

	if val.Kind() != reflect.Struct {
		return nil, false
	}

	// Direct case-sensitive Go field lookup
	fieldVal := val.FieldByName(parts[0])
	if !fieldVal.IsValid() {
		return nil, false
	}

	// Recurse down the remaining path
	return resolveFieldPath(fieldVal, parts[1:])
}

// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"

	"gopkg.in/yaml.v2"
)

func printJSON[T any](items []T, out io.Writer) error {
	bytes, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, string(bytes))
	return err
}

func printYAML[T any](items []T, out io.Writer) error {
	bytes, err := yaml.Marshal(items)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(out, string(bytes))
	return err
}

func printTable[T any](items []T, cols []ColumnDef, out io.Writer) error {
	w := tabwriter.NewWriter(out, 0, 0, 4, ' ', 0)

	// Print headers
	for i, col := range cols {
		if i > 0 {
			fmt.Fprint(w, "\t")
		}
		fmt.Fprint(w, col.Header)
	}
	fmt.Fprintln(w)

	// Print rows
	for _, item := range items {
		v := reflect.ValueOf(item)
		for i, col := range cols {
			if i > 0 {
				fmt.Fprint(w, "\t")
			}
			field := v.FieldByName(col.Field)
			if field.IsValid() {
				fmt.Fprintf(w, "%v", field.Interface())
			}
		}
		fmt.Fprintln(w)
	}

	return w.Flush()
}

func printOutput[T any](format string, items []T, cols []ColumnDef, out io.Writer) error {
	switch format {
	case "json":
		return printJSON(items, out)
	case "yaml":
		return printYAML(items, out)
	default:
		return printTable(items, cols, out)
	}
}

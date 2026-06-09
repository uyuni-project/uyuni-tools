// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package get

import "os"

func runGet[T any](flags *getOptions, res Resource[T], name string) error {
	client, err := newClient(flags.ConnectionDetails)
	if err != nil {
		return err
	}

	var items []T
	if name != "" {
		item, err := res.Get(client, name)
		if err != nil {
			return err
		}
		items = []T{item}
	} else {
		items, err = res.List(client)
		if err != nil {
			return err
		}
	}

	cols := res.Columns()

	return printOutput(flags.OutputFormat, items, cols, os.Stdout)
}

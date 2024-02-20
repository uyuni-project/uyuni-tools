// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/uyuni-project/uyuni-tools/mgradm/cmd"
)

// Run runs the `mgradm` root command.
func Run() error {
	run, err := cmd.NewUyuniadmCommand()
	if err != nil {
		return err
	}
	return run.Execute()
}

func main() {
	if err := Run(); err != nil {
		os.Exit(1)
	}
}

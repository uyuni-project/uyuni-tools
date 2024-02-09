// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd"
)

// Run runs the `mgrpxy` root command.
func Run() error {
	run, err := cmd.NewUyuniproxyCommand()
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

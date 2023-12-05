// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd"
)

// Run runs the `mgrctl` root command
func Run() error {
	return cmd.NewUyunictlCommand().Execute()
}

func main() {
	if err := Run(); err != nil {
		os.Exit(1)
	}
}

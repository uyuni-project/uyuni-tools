// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/chai2010/gettext-go"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Run runs the `mgrctl` root command.
func Run() error {
	gettext.BindLocale(gettext.New("mgrctl", utils.LocaleRoot))
	run, err := cmd.NewUyunictlCommand()
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

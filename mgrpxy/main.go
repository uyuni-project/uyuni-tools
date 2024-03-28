// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/chai2010/gettext-go"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/cmd"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Run runs the `mgrpxy` root command.
func Run() error {
	gettext.BindLocale(gettext.New("mgrpxy", utils.LocaleRoot))
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

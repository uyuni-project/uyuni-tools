// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/chai2010/gettext-go"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd"
	l10n_utils "github.com/uyuni-project/uyuni-tools/shared/l10n/utils"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Run runs the `mgrctl` root command.
func Run() error {
	gettext.BindLocale(gettext.New("mgrctl", utils.LocaleRoot, l10n_utils.New(utils.LocaleRoot)))
	cobra.EnableCaseInsensitive = true
	run := cmd.NewUyunictlCommand()

	return run.Execute()
}

func main() {
	if err := Run(); err != nil {
		os.Exit(1)
	}
}

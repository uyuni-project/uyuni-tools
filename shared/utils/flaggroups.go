// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// Group Structure to manage groups for commands.
type Group struct {
	ID    string
	Title string
}

// FlagHelpGroupAnnotation is an annotation to store the flag group to.
const FlagHelpGroupAnnotation = "cobra_annotation_flag_help_group"

var commandGroups = make(map[*cobra.Command][]Group)

func usageByFlagHelpGroupID(cmd *cobra.Command, groupID string) string {
	fs := &flag.FlagSet{}

	cmd.LocalFlags().VisitAll(func(f *flag.Flag) {
		if _, ok := f.Annotations[FlagHelpGroupAnnotation]; !ok {
			if groupID == "" {
				fs.AddFlag(f)
			}
			return
		}

		if id := f.Annotations[FlagHelpGroupAnnotation][0]; id == groupID {
			fs.AddFlag(f)
		}
	})

	return fs.FlagUsages()
}

func usageFunc(cmd *cobra.Command) error {
	flagsUsage := ""
	for _, group := range commandGroups[cmd] {
		flagsUsage += group.Title + ":\n"
		flagsUsage += usageByFlagHelpGroupID(cmd, group.ID)
		flagsUsage += "\n"
	}

	genericFlagsUsage := usageByFlagHelpGroupID(cmd, "")
	if len(genericFlagsUsage) > 0 {
		flagsUsage = L("Flags:\n") + genericFlagsUsage + "\n" + flagsUsage
	}

	template := cmd.UsageTemplate()
	re := regexp.MustCompile(`(?s)\{\{if \.HasAvailableLocalFlags\}\}.*?\{\{end\}\}`)
	template = re.ReplaceAllString(template, "\n\n"+flagsUsage)
	cmd.SetUsageTemplate(template)

	// call the original UsageFunc with the modified template
	blankCmd := cobra.Command{}
	cmd.SetUsageFunc(blankCmd.UsageFunc())
	origUsageFunc := cmd.UsageFunc()
	cmd.SetUsageFunc(usageFunc)

	return origUsageFunc(cmd)
}

// AddFlagHelpGroup adds a new flags group.
func AddFlagHelpGroup(cmd *cobra.Command, groups ...*Group) error {
	for _, group := range groups {
		commandGroups[cmd] = append(commandGroups[cmd], *group)
	}

	cmd.SetUsageFunc(usageFunc)
	return nil
}

// AddFlagToHelpGroupID adds a flag to a group.
func AddFlagToHelpGroupID(cmd *cobra.Command, flag, groupID string) error {
	lf := cmd.Flags()

	found := false
	for _, existing := range commandGroups[cmd] {
		if existing.ID == groupID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf(L("no such flag help group: %v"), groupID)
	}

	err := lf.SetAnnotation(flag, FlagHelpGroupAnnotation, []string{groupID})
	if err != nil {
		return err
	}

	return nil
}

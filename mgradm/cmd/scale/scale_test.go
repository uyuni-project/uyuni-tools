// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package scale

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils/flagstests"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestParamsParsing(t *testing.T) {
	args := []string{
		"--replicas", "2",
		"--pullPolicy", "never",
		"--registry", "myOldRegistry",
		"--coco-image", "cocoimg",
		"--coco-tag", "cocotag",
		"--hubxmlrpc-image", "hubimg",
		"--hubxmlrpc-tag", "hubtag",
		"--saline-image", "salineimg",
		"--saline-tag", "salinetag",
		"some-service",
	}
	args = append(args, flagstests.RegistryImageFlagsTestArgs...)
	args = append(args, flagstests.SCCFlagTestArgs...)

	// Test function asserting that the args are properly parsed
	tester := func(_ *types.GlobalFlags, flags *scaleFlags, _ *cobra.Command, args []string) error {
		testutils.AssertEquals(t, "Error parsing --replicas", 2, flags.Replicas)
		testutils.AssertEquals(t, "Error parsing --pullPolicy", "never", flags.Image.PullPolicy)
		testutils.AssertEquals(t, "Error parsing --coco-image", "cocoimg", flags.Coco.Image.Name)
		testutils.AssertEquals(t, "Error parsing --coco-tag", "cocotag", flags.Coco.Image.Tag)
		testutils.AssertEquals(t, "Error parsing --hubxmlrpc-image", "hubimg", flags.HubXmlrpc.Image.Name)
		testutils.AssertEquals(t, "Error parsing --hubxmlrpc-tag", "hubtag", flags.HubXmlrpc.Image.Tag)
		testutils.AssertEquals(t, "Error parsing --saline-image", "salineimg", flags.Saline.Image.Name)
		testutils.AssertEquals(t, "Error parsing --saline-tag", "salinetag", flags.Saline.Image.Tag)
		testutils.AssertEquals(t, "Error parsing the service name", "some-service", args[0])
		flagstests.AssertRegistryFlag(t, &flags.Image.Registry)
		flagstests.AssertSCCFlag(t, &flags.SCC)
		return nil
	}

	globalFlags := types.GlobalFlags{}
	cmd := newCmd(&globalFlags, tester)

	testutils.AssertHasAllFlags(t, cmd, args)

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("command failed with error: %s", err)
	}
}

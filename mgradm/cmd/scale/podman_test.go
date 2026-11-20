// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package scale

import (
	"path"
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestServiceImage(t *testing.T) {
	originalGetServiceImage := getServiceImage
	defer func() {
		getServiceImage = originalGetServiceImage
	}()

	const currentImage = "registry.suse.com/suse/multi-linux-manager/5.2/x86_64/server-saline:5.2.0"
	getServiceImage = func(_ string) string {
		return currentImage
	}

	testCases := []struct {
		args     []string
		expected types.ImageFlags
	}{
		{[]string{}, types.ImageFlags{Name: currentImage}},
		{[]string{"--saline-image", "salineimg"}, types.ImageFlags{Name: "salineimg"}},
		{[]string{"--saline-tag", "salinetag"},
			types.ImageFlags{Name: path.Join(utils.DefaultImagePrefix, "server-saline"), Tag: "salinetag"}},
	}

	for _, testCase := range testCases {
		tester := func(_ *types.GlobalFlags, flags *scaleFlags, cmd *cobra.Command, _ []string) error {
			actual := serviceImage(cmd, "saline", flags.Saline.Image, podman.SalineService)
			testutils.AssertEquals(t, "wrong image", testCase.expected, actual)
			return nil
		}
		cmd := newCmd(&types.GlobalFlags{}, tester)
		cmd.SetArgs(append(testCase.args, podman.SalineService))
		if err := cmd.Execute(); err != nil {
			t.Errorf("command failed with error: %s", err)
		}
	}
}

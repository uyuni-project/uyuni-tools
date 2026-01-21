// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func TestHasDebugPorts(t *testing.T) {
	data := map[string]bool{
		`[Service]
ExecStart=/bin/sh -c '/usr/bin/podman run \
        --name uyuni-server \
        --hostname uyuni-server.mgr.internal \
        --rm --cap-add NET_RAW \
        -p 80:80 \
        -p 8003:8003 \
        -p 4505:4505`: true,
		`[Service]
ExecStart=/bin/sh -c '/usr/bin/podman run \
        --name uyuni-server \
        --hostname uyuni-server.mgr.internal \
        --rm --cap-add NET_RAW \
        -p 80:80 \
        -p 4505:4505`: false,
	}

	for definition, expected := range data {
		actual := hasDebugPorts([]byte(definition))
		testutils.AssertEquals(t, "Unexpected result for "+definition, expected, actual)
	}
}

func TestGetMirrorPath(t *testing.T) {
	data := map[string]string{
		`[Service]
ExecStart=/bin/sh -c '/usr/bin/podman run \
        --name uyuni-server \
        --hostname uyuni-server.mgr.internal \
        --rm --cap-add NET_RAW \
        -p 80:80 \
        -p 4505:4505`: "",
		`[Service]
ExecStart=/bin/sh -c '/usr/bin/podman run \
        --name uyuni-server \
        --hostname uyuni-server.mgr.internal \
        --rm --cap-add NET_RAW \
		-v   /path/to/mirror:/mirror \
        -p 80:80 \
        -p 4505:4505`: "/path/to/mirror",
		`[Service]
ExecStart=/bin/sh -c '/usr/bin/podman run \
        --name uyuni-server \
        --hostname uyuni-server.mgr.internal \
		--rm --cap-add NET_RAW -v /path/to/mirror:/mirror \
        -p 80:80 \
        -p 4505:4505`: "/path/to/mirror",
	}

	for definition, expected := range data {
		actual := getMirrorPath([]byte(definition))
		testutils.AssertEquals(t, "Unexpected result for "+definition, expected, actual)
	}
}

func TestRunPgsqlVersionUpgrade(t *testing.T) {
	cases := []struct {
		registry      string
		image         types.ImageFlags
		upgradeImage  types.ImageFlags
		expectedImage string
	}{
		// Default Uyuni case with global tag set
		{
			"registry.opensuse.org/uyuni",
			types.ImageFlags{
				Name: "registry.opensuse.org/uyuni/server",
				Registry: types.Registry{
					Host: "registry.opensuse.org/uyuni",
				},
				Tag:        "2025.08",
				PullPolicy: "ifnotpresent",
			},
			types.ImageFlags{},
			"registry.opensuse.org/uyuni/server-database-migration:2025.08",
		},
		// own registry case with a special image for the main server but not upgrade
		{
			"registry.example.com/product",
			types.ImageFlags{
				Name: "registry.example.com/product/server",
				Registry: types.Registry{
					Host: "registry.opensuse.org/uyuni",
				},
				Tag:        "fix-123",
				PullPolicy: "always",
			},
			types.ImageFlags{
				Name: "registry.example.com/product/server-database-migration",
				Tag:  "4.5.2",
			},
			"registry.example.com/product/server-database-migration:4.5.2",
		},
	}

	expectedAuthfile := "authfile to pass"
	for i, testCase := range cases {
		prepareImage = func(authFile string, image string, pullPolicy string, _ bool) (string, error) {
			// test that the image computation
			testutils.AssertEquals(t, "auth file not passed down", expectedAuthfile, authFile)
			testutils.AssertEquals(t, fmt.Sprintf("case %d: wrong image", i), testCase.expectedImage, image)
			testutils.AssertEquals(t, fmt.Sprintf("case %d: wrong pull policy", i), testCase.image.PullPolicy, pullPolicy)
			return image, nil
		}
		runContainer = func(_ string, image string, _ []types.VolumeMount, _ []string, _ []string) error {
			testutils.AssertEquals(t, fmt.Sprintf("case %d: wrong image used for container", i), testCase.expectedImage, image)
			return nil
		}
		_ = RunPgsqlVersionUpgrade(expectedAuthfile, testCase.image, testCase.upgradeImage, "14", "16")
	}
}

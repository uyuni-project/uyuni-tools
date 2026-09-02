// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	sharedPodman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestPrepareHostPreservesInspectionErrors(t *testing.T) {
	originalInspect := inspect
	t.Cleanup(func() { inspect = originalInspect })
	inspect = func(_, _ string) (*utils.InspectData, error) {
		return nil, errors.New("database inspection failed")
	}

	_, err := prepareHost("server-image", "postgres-image")
	if err == nil {
		t.Fatal("expected prepareHost to return an error")
	}
	if !strings.Contains(err.Error(), "cannot inspect podman values: database inspection failed") {
		t.Fatalf("expected the original inspection error, got: %v", err)
	}
}

func TestEnsureServicesRunning(t *testing.T) {
	cases := []struct {
		name    string
		running []string
		wantErr bool
	}{
		{name: "both services running", running: []string{sharedPodman.ServerService, sharedPodman.DBService}},
		{name: "server stopped", running: []string{sharedPodman.DBService}, wantErr: true},
		{name: "database stopped", running: []string{sharedPodman.ServerService}, wantErr: true},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			systemd := sharedPodman.NewSystemdWithDriver(&testutils.FakeSystemdDriver{Running: testCase.running})
			err := ensureServicesRunning(systemd)
			if testCase.wantErr && err == nil {
				t.Fatal("expected an error")
			}
			if !testCase.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if testCase.wantErr && !strings.Contains(err.Error(), "Uyuni server and database services must be running") {
				t.Fatalf("expected an actionable error message, got: %v", err)
			}
		})
	}
}

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
			"registry.opensuse.org",
			types.ImageFlags{
				Name: "uyuni/server",
				Registry: types.Registry{
					Host: "registry.opensuse.org",
				},
				Tag:        "2025.08",
				PullPolicy: "ifnotpresent",
			},
			types.ImageFlags{
				Name: "uyuni/server-database-migration",
			},
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
		_ = RunPgsqlVersionUpgrade(expectedAuthfile, testCase.image, testCase.upgradeImage, []types.VolumeMount{})
	}
}

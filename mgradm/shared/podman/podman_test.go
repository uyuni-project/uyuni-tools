// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestHasDebugPorts(t *testing.T) {
	data := map[string]bool{
		`[Service]
ExecStart=/bin/sh -c '/usr/bin/podman run \
        --name uyuni-server \
        --hostname uyuni-server.mgr.internal \
        --rm --cap-add NET_RAW --tmpfs /run -v cgroup:/sys/fs/cgroup:rw \
        -p 80:80 \
        -p 8003:8003 \
        -p 4505:4505`: true,
		`[Service]
ExecStart=/bin/sh -c '/usr/bin/podman run \
        --name uyuni-server \
        --hostname uyuni-server.mgr.internal \
        --rm --cap-add NET_RAW --tmpfs /run -v cgroup:/sys/fs/cgroup:rw \
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
        --rm --cap-add NET_RAW --tmpfs /run -v cgroup:/sys/fs/cgroup:rw \
        -p 80:80 \
        -p 4505:4505`: "",
		`[Service]
ExecStart=/bin/sh -c '/usr/bin/podman run \
        --name uyuni-server \
        --hostname uyuni-server.mgr.internal \
        --rm --cap-add NET_RAW --tmpfs /run -v cgroup:/sys/fs/cgroup:rw \
		-v   /path/to/mirror:/mirror \
        -p 80:80 \
        -p 4505:4505`: "/path/to/mirror",
		`[Service]
ExecStart=/bin/sh -c '/usr/bin/podman run \
        --name uyuni-server \
        --hostname uyuni-server.mgr.internal \
		--rm --cap-add NET_RAW -v /path/to/mirror:/mirror --tmpfs /run -v cgroup:/sys/fs/cgroup:rw \
        -p 80:80 \
        -p 4505:4505`: "/path/to/mirror",
	}

	for definition, expected := range data {
		actual := getMirrorPath([]byte(definition))
		testutils.AssertEquals(t, "Unexpected result for "+definition, expected, actual)
	}
}

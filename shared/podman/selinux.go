// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import "github.com/rs/zerolog"

// IsSELinuxEnabled reports whether SELinux is enabled or not.
// It relies on selinuxenabled tool.
func IsSELinuxEnabled() bool {
	_, err := runCmdOutput(zerolog.DebugLevel, "selinuxenabled")
	return err == nil
}

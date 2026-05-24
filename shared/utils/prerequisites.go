// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"net"

	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// CheckPort checks if a given port is available.
func CheckPort(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return Errorf(err, L("port %d is already in use"), port)
	}
	l.Close()
	return nil
}

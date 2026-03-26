//go:build !linux
// +build !linux

// SPDX-FileCopyrightText: 2024 SUSE LLC
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

// CheckStorage checks if the given path has at least requiredMinGB free space.
func CheckStorage(path string, requiredMinGB uint64) error {
	// Storage check is currently only implemented for Linux.
	return nil
}

// CheckMemory checks if the system has at least requiredMinGB memory.
func CheckMemory(requiredMinGB uint64) error {
	// Memory check is currently only implemented for Linux.
	return nil
}

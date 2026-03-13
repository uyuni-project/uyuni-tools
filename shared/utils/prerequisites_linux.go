//go:build linux
// +build linux

// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"net"
	"syscall"

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
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return Errorf(err, L("failed to check storage for %s"), path)
	}

	freeSpace := stat.Bavail * uint64(stat.Bsize)
	requiredSpace := requiredMinGB * 1024 * 1024 * 1024

	if freeSpace < requiredSpace {
		return fmt.Errorf(
			L("insufficient storage in %s: requires at least %d GB, but only %.2f GB available"),
			path,
			requiredMinGB,
			float64(freeSpace)/(1024*1024*1024),
		)
	}
	return nil
}

// CheckMemory checks if the system has at least requiredMinGB memory.
func CheckMemory(requiredMinGB uint64) error {
	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err != nil {
		return Errorf(err, L("failed to read system memory"))
	}

	totalRAM := uint64(info.Totalram) * uint64(info.Unit)
	requiredRAM := requiredMinGB * 1024 * 1024 * 1024

	if totalRAM < requiredRAM {
		return fmt.Errorf(
			L("insufficient memory: requires at least %d GB, but system has %.2f GB"),
			requiredMinGB,
			float64(totalRAM)/(1024*1024*1024),
		)
	}
	return nil
}

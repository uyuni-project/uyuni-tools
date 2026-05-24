//go:build !linux
// +build !linux

// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// CheckStorage checks if the given path has at least requiredMinGB free space.
func CheckStorage(path string, requiredMinGB uint64) error {
	// Storage check is currently only implemented for Linux.
	log.Warn().Msg(L("storage check is not implemented for this operating system"))
	return nil
}

// CheckMemory checks if the system has at least requiredMinGB memory.
func CheckMemory(requiredMinGB uint64) error {
	// Memory check is currently only implemented for Linux.
	log.Warn().Msg(L("memory check is not implemented for this operating system"))
	return nil
}

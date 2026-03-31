// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

// PortMap describes a port.
type PortMap struct {
	Exposed  int
	Port     int
	Protocol string
}

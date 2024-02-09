// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

// PortMap describes a port.
type PortMap struct {
	Name     string
	Exposed  int
	Port     int
	Protocol string
}

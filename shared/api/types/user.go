// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

// User describes an Uyuni user in the API.
type User struct {
	Login     string
	Password  string
	FirstName string
	LastName  string
	Email     string
}

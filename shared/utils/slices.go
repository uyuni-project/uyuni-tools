// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

// Contains returns true if a string is contained in a string slice.
func Contains(slice []string, needle string) bool {
	for _, item := range slice {
		if item == needle {
			return true
		}
	}
	return false
}

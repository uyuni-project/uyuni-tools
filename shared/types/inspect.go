// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

/* InspectData represents CLI command to run in the container
* and the variable where the output is stored.
 */
type InspectData struct {
	Variable string
	CLI      string
}

// NewInspectData creates an InspectData instance.
func NewInspectData(variable string, cli string) InspectData {
	return InspectData{
		Variable: variable,
		CLI:      cli,
	}
}

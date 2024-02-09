// SPDX-FileCopyrightText: 2024 SUSE LLC
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

/* InspectFile represent where the inspect file should be stored
* and the command to run in the container.
 */
type InspectFile struct {
	Directory string
	Basename  string
	Commands  []InspectData
}

// NewInspectData creates an InspectData instance.
func NewInspectData(variable string, cli string) InspectData {
	return InspectData{
		Variable: variable,
		CLI:      cli,
	}
}

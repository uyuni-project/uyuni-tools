// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

type InspectData struct {
	Variable string
	CLI      string
}

type InspectFile struct {
	Directory string
	Basename  string
	Commands  []InspectData
}

func InspectDataConstructor(variable string, cli string) InspectData {
	return InspectData{
		Variable: variable,
		CLI:      cli,
	}
}

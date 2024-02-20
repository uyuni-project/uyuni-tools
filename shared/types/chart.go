// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

// ChartFlags represents the flags required by charts.
type ChartFlags struct {
	Namespace string
	Chart     string
	Version   string
	Values    string
}

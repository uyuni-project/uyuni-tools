// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

type ChartFlags struct {
	Namespace string
	Chart     string
	Version   string
	Values    string
}

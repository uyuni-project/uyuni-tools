// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

type UyuniServiceOption struct {
	Name        string
	Value       interface{}
	Description string
}

type UyuniServiceReplicas struct {
	Max     uint
	Min     uint
	Default uint
}

type UyuniService struct {
	Name        string
	Image       ImageFlags
	Description string
	Replicas    UyuniServiceReplicas
	Options     []UyuniServiceOption
}

var SingleMandatoryReplica = UyuniServiceReplicas{Max: 1, Min: 1, Default: 1}
var SingleOptionalReplica = UyuniServiceReplicas{Max: 1, Min: 0, Default: 0}

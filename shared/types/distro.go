// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

// Distribution contains information about the distribution.
type Distribution struct {
	TreeLabel    string
	BasePath     string
	ChannelLabel string
	InstallType  string
}

type DistributionDetails struct {
	Name    string
	Version string
	Arch    Arch
}

type Arch string

const (
	UnknownArch Arch = "unknown"
	AMD64       Arch = "x86_64"
	AArch64     Arch = "aarch64"
	S390X       Arch = "s390x"
	PPC64LE     Arch = "ppc64le"
)

func GetArch(a string) Arch {
	switch a {
	case "x86_64":
		return AMD64
	case "aarch64":
		return AArch64
	case "s390x":
		return S390X
	case "ppc64le":
		return PPC64LE
	}
	return UnknownArch
}

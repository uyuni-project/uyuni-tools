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

// DistributionDetails contains distro details passed from the command line.
type DistributionDetails struct {
	Name    string
	Version string
	Arch    Arch
}

// Arch type to store architecture.
type Arch string

// Constants for supported archhitectures.
const (
	UnknownArch Arch = "unknown"
	AMD64       Arch = "x86_64"
	AArch64     Arch = "aarch64"
	S390X       Arch = "s390x"
	PPC64LE     Arch = "ppc64le"
)

// GetArch translates string representation of architecture to Arch type.
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

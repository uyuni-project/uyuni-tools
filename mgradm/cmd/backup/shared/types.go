// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

type Flagpole struct {
	SkipVolumes  []string `mapstructure:"skipvolumes"`
	ExtraVolumes []string `mapstructure:"extravolumes"`
	SkipDatabase bool     `mapstructure:"skipdatabase"`
	SkipImages   bool     `mapstructure:"skipimages"`
	SkipConfig   bool     `mapstructure:"skipconfig"`
	NoRestart    bool     `mapstructure:"norestart"`
	DryRun       bool     `mapstructure:"dryrun"`
	ForceRestore bool     `mapstructure:"force"`
	SkipExisting bool     `mapstructure:"continue"`
	SkipVerify   bool     `mapstructure:"skipverify"`
}

// Backup error indicating if something was already backed up (resp. restored) or not.
type BackupError struct {
	Err         error
	DataRemains bool
	Abort       bool
}

func (e *BackupError) Error() string {
	return e.Err.Error()
}

func (e *BackupError) Unwrap() error {
	return e.Err
}

// Wrap error with metadata indicating this error was fatal and job was aborted.
func AbortError(err error, dataRemains bool) error {
	if err == nil {
		return nil
	}
	return &BackupError{
		Err:         err,
		DataRemains: dataRemains,
		Abort:       true,
	}
}

// Wrap error with metadata indicating this error was not fatal.
func ReportError(err error) error {
	if err == nil {
		return nil
	}
	return &BackupError{
		Err:         err,
		DataRemains: true,
		Abort:       false,
	}
}

// Map of podman secret name and value.
type BackupSecretMap struct {
	Name   string
	Secret string
}

type NetworkSubnet struct {
	Subnet  string
	Gateway string
}

type PodanNetworkConfigData struct {
	Subnets           []NetworkSubnet `mapstructure:"subnets"`
	NetworkInsterface string          `mapstructure:"network_interface"`
	NetworkDNSServers []string        `mapstructure:"network_dns_servers"`
}

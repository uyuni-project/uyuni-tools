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
	DryRun       bool     `manstructure:"dryrun"`
}

// Backup error indicating if something was already backed up or not.
type BackupError struct {
	Err         error
	DataRemains bool
}

func (e *BackupError) Error() string {
	return e.Err.Error()
}

func (e *BackupError) Unwrap() error {
	return e.Err
}

func ReportError(err error, dataRemain bool) *BackupError {
	return &BackupError{
		Err:         err,
		DataRemains: dataRemain,
	}
}

// Map of podman secret name and value.
type BackupSecretMap struct {
	Name   string
	Secret string
}

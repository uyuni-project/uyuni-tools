// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"

	"github.com/spf13/viper"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// ReadInspectData returns an unmarshalled object of type T from the data as a string.
//
// This function is most likely to be used for the implementation of the inspectors, but can also be used directly.
func ReadInspectData[T any](data []byte) (*T, error) {
	viper.SetConfigType("env")
	if err := viper.MergeConfig(bytes.NewBuffer(data)); err != nil {
		return nil, Error(err, L("cannot read config"))
	}

	var inspectResult T
	if err := viper.Unmarshal(&inspectResult); err != nil {
		return nil, Error(err, L("failed to unmarshal the inspected data"))
	}
	return &inspectResult, nil
}

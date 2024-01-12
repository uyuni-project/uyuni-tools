// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const envPrefix = "UYUNI"
const appName = "uyuni-tools"

func ReadConfig(configPath string, configFilename string, cmd *cobra.Command) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigType("yaml")
	v.SetConfigName(configFilename)

	if configPath != "" {
		log.Info().Msgf("Using config file %s", configPath)
		v.SetConfigFile(configPath)
	} else {
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				log.Err(err).Msg("Failed to find home directory")
			} else {
				xdgConfigHome = path.Join(home, ".config")
			}
		}
		if xdgConfigHome != "" {
			v.AddConfigPath(path.Join(xdgConfigHome, appName))
		}
		v.AddConfigPath(".")
	}

	if err := bindFlags(cmd, v); err != nil {
		return nil, err
	}

	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// TODO Provide help on the config file format
			return nil, fmt.Errorf("failed to parse configuration file %s: %s", v.ConfigFileUsed(), err)
		}
	}

	v.SetEnvPrefix(envPrefix)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.AutomaticEnv()

	return v, nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
	var errors []error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", ".")
		if err := v.BindPFlag(configName, f); err != nil {
			errors = append(errors, fmt.Errorf("failed to bind %s config to parameter %s: %s", configName, f.Name, err))
		}
	})

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

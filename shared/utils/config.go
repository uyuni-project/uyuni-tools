// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
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

func ReadConfig(configPath string, configFilename string, cmd *cobra.Command) *viper.Viper {
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
			}
			xdgConfigHome = path.Join(home, ".config")
		}
		v.AddConfigPath(path.Join(xdgConfigHome, appName))
		v.AddConfigPath(".")
	}

	bindFlags(cmd, v)

	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// TODO Provide help on the config file format
			log.Fatal().Err(err).Msgf("Failed to parse configuration file %s", v.ConfigFileUsed())
		}
	}

	v.SetEnvPrefix(envPrefix)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.AutomaticEnv()

	return v
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", ".")
		if err := v.BindPFlag(configName, f); err != nil {
			log.Fatal().Err(err).Msgf("Failed to bind %s config to parameter %s", configName, f.Name)
		}
	})
}

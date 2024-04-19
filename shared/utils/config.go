// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

const envPrefix = "UYUNI"
const appName = "uyuni-tools"
const configFilename = "config.yaml"

// ReadConfig parse configuration file and env variables a return parameters.
func ReadConfig(configPath string, cmd *cobra.Command) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigType("yaml")
	v.SetConfigName(configFilename)

	if configPath != "" {
		log.Info().Msgf(L("Using config file %s"), configPath)
		v.SetConfigFile(configPath)
	} else {
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				log.Err(err).Msg(L("Failed to find home directory"))
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
			return nil, fmt.Errorf(L("failed to parse configuration file %s: %s"), v.ConfigFileUsed(), err)
		}
	}

	v.SetEnvPrefix(envPrefix)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.AutomaticEnv()

	return v, nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable).
func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
	var errors []error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", ".")
		if err := v.BindPFlag(configName, f); err != nil {
			errors = append(errors, fmt.Errorf(L("failed to bind %s config to parameter %s: %s"), configName, f.Name, err))
		}
	})

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// GetLocalizedUsageTemplate provides the help template, but localized.
func GetLocalizedUsageTemplate() string {
	return L(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
}

// GetConfigHelpCommand provides a help command describing the config file and environment variables.
func GetConfigHelpCommand() *cobra.Command {
	var configTemplate = L(`
Configuration:

  All the non-global flags can alternatively be passed as configuration.
  
  The configuration file is a YAML file with entries matching the flag name.
  The name of a flag is the part after the '--' of the command line parameter.
  Every '_' character in the flag name means a nested property.
  
  For instance the '--tz CEST' and '--ssl-password secret' will be mapped to
  this YAML configuration:
  
    tz: CEST
    ssl:
      password: secret
  
  The configuration file will be searched in the following places and order:
  · $XDG_CONFIG_HOME/{{ .Name }}/{{ .ConfigFile }}
  · $HOME/.config/{{ .Name }}/{{ .ConfigFile }}
  · $PWD/{{ .ConfigFile }}
  · the value of the --config flag


Environment variables:

  All the non-global flags can also be passed as environment variables.
  
  The environment variable name is the flag name with '-' replaced by with '_'
  and the {{ .EnvPrefix }} prefix.
  
  For example the '--tz CEST' flag will be mapped to '{{ .EnvPrefix }}_TZ'
  and '--ssl-password' flags to '{{ .EnvPrefix }}_SSL_PASSWORD' 
`)

	cmd := &cobra.Command{
		Use:   "config",
		Short: L("help on configuration file and environment variables"),
	}
	t := template.Must(template.New("help").Parse(configTemplate))
	var helpBuilder strings.Builder
	if err := t.Execute(&helpBuilder, configTemplateData{
		EnvPrefix:  envPrefix,
		Name:       appName,
		ConfigFile: configFilename,
	}); err != nil {
		log.Fatal().Err(err).Msg(L("failed to compute config help command"))
	}
	cmd.SetHelpTemplate(helpBuilder.String())
	return cmd
}

type configTemplateData struct {
	EnvPrefix  string
	ConfigFile string
	Name       string
}

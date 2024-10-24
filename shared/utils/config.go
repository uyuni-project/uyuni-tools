// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"errors"
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

// GlobalConfigFilename is the path for the global configuration.
const GlobalConfigFilename = "/etc/uyuni/uyuni-tools.yaml"

func addConfigurationFile(v *viper.Viper, cmd *cobra.Command, configFilename string) error {
	if FileExists(configFilename) {
		v.SetConfigFile(configFilename)
	}
	if err := bindFlags(cmd, v); err != nil {
		return err
	}
	if err := v.MergeInConfig(); err != nil {
		// It's okay if there isn't a config file
		var configFileNotFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundErr) {
			// TODO Provide help on the config file format
			return Errorf(err, L("failed to parse configuration file %s"), v.ConfigFileUsed())
		}
	}
	return nil
}

// GetUserConfigDir returns the user configuration directory.
//
// Can be $XDG_CONFIG_HOME or `homedir/.config`.
func GetUserConfigDir() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Warn().Err(err).Msg(L("Failed to find home directory"))
		} else {
			xdgConfigHome = path.Join(home, ".config")
		}
	}
	return xdgConfigHome
}

// ReadConfig parse configuration file and env variables a return parameters.
func ReadConfig(cmd *cobra.Command, configPaths ...string) (*viper.Viper, error) {
	v := viper.New()

	for _, configPath := range configPaths {
		if err := addConfigurationFile(v, cmd, configPath); err != nil {
			return v, err
		}
	}

	// once global configuration are set, set the local config file as default
	v.SetConfigType("yaml")
	v.SetConfigName(configFilename)

	xdgConfigHome := GetUserConfigDir()
	if xdgConfigHome != "" {
		v.AddConfigPath(path.Join(xdgConfigHome, appName))
	}
	v.AddConfigPath(".")

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
			errors = append(errors, Errorf(err, L("failed to bind %[1]s config to parameter %[2]s"), configName, f.Name))
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
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

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
  · /etc/uyuni/uyuni-tools.yaml
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
		Short: L("Help on configuration file and environment variables"),
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

// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type distribution struct {
	TreeLabel    string
	BasePath     string
	ChannelLabel string
	InstallType  string
}

type distroDetails struct {
	InstallType  string
	ChannelLabel string
}

func getDetailsFromDistro(distro string, version string) (distroDetails, error) {
	installerMapping := map[string]distroDetails{
		"AlmaLinux 9": {
			InstallType:  "rhel_9",
			ChannelLabel: "",
		},
		"AlmaLinux 8": {
			InstallType:  "rhel_8",
			ChannelLabel: "",
		},
		"SUSE Linux Enterprise 15": {
			InstallType:  "sles15generic",
			ChannelLabel: "",
		},
		"SUSE Linux Enterprise 12": {
			InstallType:  "sles12generic",
			ChannelLabel: "",
		},
	}
	lookupname := fmt.Sprintf("%s %s", distro, version)
	val, ok := installerMapping[lookupname]
	if ok {
		return val, nil
	}
	return distroDetails{}, fmt.Errorf("Unkown distribution '%s'", lookupname)
}

func umountAndRemove(mountpoint string) {
	umount_cmd := []string{
		"/usr/bin/umount",
		mountpoint,
	}

	if err := utils.RunCmd("/usr/bin/sudo", umount_cmd...); err != nil {
		log.Fatal().Err(err).Msgf("Unable to unmount iso file, leaving %s intact", mountpoint)
	}

	os.Remove(mountpoint)
}

func detectDistro(path string, distro *distribution) error {
	treeinfopath := filepath.Join(path, ".treeinfo")
	log.Trace().Msgf("Reading .treeinfo %s", treeinfopath)
	custom_viper := viper.New()
	custom_viper.SetConfigType("ini")
	custom_viper.SetConfigName(".treeinfo")
	custom_viper.AddConfigPath(path)
	if err := custom_viper.ReadInConfig(); err != nil {
		return err
	}

	dname := custom_viper.GetString("release.name")
	dversion := custom_viper.GetString("release.version")
	log.Debug().Msgf("Detected distro %s, version %s", dname, dversion)
	details, err := getDetailsFromDistro(dname, dversion)
	if err != nil {
		return err
	}

	*distro = distribution{
		InstallType:  details.InstallType,
		TreeLabel:    dname,
		ChannelLabel: details.ChannelLabel,
	}
	return nil
}

func registerDistro(connection *api.ConnectionDetails, distro *distribution) error {
	client := api.Init(connection)
	if err := client.Login(connection.User, connection.Password); err != nil {
		log.Error().Msg("Unable to login and register the distribution. Manual distro registration is required.")
		return err
	}
	data := map[string]interface{}{
		"treeLabel":    distro.TreeLabel,
		"basePath":     distro.BasePath,
		"channelLabel": distro.ChannelLabel,
		"installType":  distro.InstallType,
	}

	_, err := client.Post("kickstart/tree/create", data)
	if err != nil {
		log.Error().Msg("Unable to register the distribution. Manual distro registration is required.")
		return err
	}
	log.Info().Msgf("Distribution %s successfuly registered", distro.TreeLabel)
	return nil
}

func distCp(globalFlags *types.GlobalFlags, flags *flagpole, apiFlags *api.ConnectionDetails, cmd *cobra.Command, distroName string, source string) {
	cnx := utils.NewConnection(flags.Backend)
	log.Info().Msgf("Copying distribution %s\n", distroName)
	if !utils.FileExists(source) {
		log.Fatal().Msgf("Source %s does not exists", source)
	}

	dstpath := "/srv/www/distributions/" + distroName
	if utils.TestExistenceInPod(cnx, dstpath) {
		log.Fatal().Msgf("Distribution already exists: %s\n", dstpath)
	}

	srcdir := source
	if strings.HasSuffix(source, ".iso") {
		log.Debug().Msg("Source is an iso file")
		tmpdir, err := os.MkdirTemp("", "mgrctl")
		if err != nil {
			log.Fatal().Err(err)
		}
		srcdir = tmpdir
		defer umountAndRemove(srcdir)

		mount_cmd := []string{
			"/usr/bin/mount",
			"-o", "ro,loop",
			source,
			srcdir,
		}
		if out, err := utils.RunCmdOutput(zerolog.DebugLevel, "/usr/bin/sudo", mount_cmd...); err != nil {
			log.Debug().Msgf("Error mounting iso: '%s'", out)
			log.Error().Msg("Unable to mount iso file. Mount manually and try again")
		}
	}

	utils.Copy(cnx, srcdir, "server:"+dstpath, "tomcat", "susemanager")

	log.Info().Msg("Distribution has been copied")

	if apiFlags.User != "" {
		distro := distribution{
			BasePath: dstpath,
		}
		if err := detectDistro(srcdir, &distro); err != nil {
			log.Error().Msg(err.Error())
			return
		}

		if err := registerDistro(apiFlags, &distro); err != nil {
			log.Error().Msg(err.Error())
		}
	}
}

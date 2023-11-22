// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func umountAndRemove(mountpoint string) {
	umountCmd := []string{
		"/usr/bin/umount",
		mountpoint,
	}

	if err := utils.RunCmd("/usr/bin/sudo", umountCmd...); err != nil {
		log.Fatal().Err(err).Msgf("Unable to unmount iso file, leaving %s intact", mountpoint)
	}

	os.Remove(mountpoint)
}

func registerDistro(connection *api.ConnectionDetails, distro *types.Distribution) error {
	client, err := api.Init(connection)
	if err != nil {
		log.Error().Msg("Unable to login and register the distribution. Manual distro registration is required.")
		return err
	}
	data := map[string]interface{}{
		"treeLabel":    distro.TreeLabel,
		"basePath":     distro.BasePath,
		"channelLabel": distro.ChannelLabel,
		"installType":  distro.InstallType,
	}

	_, err = client.Post("kickstart/tree/create", data)
	if err != nil {
		log.Error().Msg("Unable to register the distribution. Manual distro registration is required.")
		return err
	}
	log.Info().Msgf("Distribution %s successfuly registered", distro.TreeLabel)
	return nil
}

func distCp(globalFlags *types.GlobalFlags, flags *flagpole, apiFlags *api.ConnectionDetails, cmd *cobra.Command, distroName string, source string, channelLabel string) {
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

		mountCmd := []string{
			"/usr/bin/mount",
			"-o", "ro,loop",
			source,
			srcdir,
		}
		if out, err := utils.RunCmdOutput(zerolog.DebugLevel, "/usr/bin/sudo", mountCmd...); err != nil {
			log.Debug().Msgf("Error mounting iso: '%s'", out)
			log.Error().Msg("Unable to mount iso file. Mount manually and try again")
		}
	}

	utils.Copy(cnx, srcdir, "server:"+dstpath, "tomcat", "susemanager")

	log.Info().Msg("Distribution has been copied")

	if apiFlags.User != "" {
		distro := types.Distribution{
			BasePath: dstpath,
		}
		if err := detectDistro(srcdir, channelLabel, flags, &distro); err != nil {
			log.Error().Msg(err.Error())
			return
		}

		if err := registerDistro(apiFlags, &distro); err != nil {
			log.Error().Msg(err.Error())
		}
	}
}

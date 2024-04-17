// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package distro

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func umountAndRemove(mountpoint string) {
	umountCmd := []string{
		"/usr/bin/umount",
		mountpoint,
	}

	if err := utils.RunCmd("/usr/bin/sudo", umountCmd...); err != nil {
		log.Fatal().Err(err).Msgf(L("Unable to unmount ISO image, leaving %s intact"), mountpoint)
	}

	os.Remove(mountpoint)
}

func registerDistro(connection *api.ConnectionDetails, distro *types.Distribution) error {
	client, err := api.Init(connection)
	if err != nil {
		return fmt.Errorf(L("unable to login and register the distribution. Manual distro registration is required: %s"), err)
	}
	data := map[string]interface{}{
		"treeLabel":    distro.TreeLabel,
		"basePath":     distro.BasePath,
		"channelLabel": distro.ChannelLabel,
		"installType":  distro.InstallType,
	}

	_, err = client.Post("kickstart/tree/create", data)
	if err != nil {
		return fmt.Errorf(L("unable to register the distribution. Manual distro registration is required: %s"), err)
	}
	log.Info().Msgf(L("Distribution %s successfully registered"), distro.TreeLabel)
	return nil
}

func distroCp(
	globalFlags *types.GlobalFlags,
	flags *flagpole,
	cmd *cobra.Command,
	args []string,
) error {
	distroName := args[1]
	source := args[0]

	var channelLabel string
	if len(args) == 3 {
		channelLabel = args[2]
	} else {
		channelLabel = ""
	}
	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)
	log.Info().Msgf(L("Copying distribution %s\n"), distroName)
	if !utils.FileExists(source) {
		return fmt.Errorf(L("source %s does not exists"), source)
	}

	dstpath := "/srv/www/distributions/" + distroName
	if cnx.TestExistenceInPod(dstpath) {
		return fmt.Errorf(L("distribution already exists: %s"), dstpath)
	}

	srcdir := source
	if strings.HasSuffix(source, ".iso") {
		log.Debug().Msg("Source is an ISO image")
		tmpdir, err := os.MkdirTemp("", "mgrctl")
		if err != nil {
			return err
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
			log.Debug().Msgf("Error mounting ISO image: '%s'", out)
			return errors.New(L("unable to mount ISO image. Mount manually and try again"))
		}
	}

	if err := cnx.Copy(srcdir, "server:"+dstpath, "tomcat", "susemanager"); err != nil {
		return fmt.Errorf(L("cannot copy %s: %s"), dstpath, err)
	}

	log.Info().Msg(L("Distribution has been copied"))

	if flags.ConnectionDetails.User != "" {
		distro := types.Distribution{
			BasePath: dstpath,
		}
		if err := detectDistro(srcdir, channelLabel, flags, &distro); err != nil {
			return err
		}

		if err := registerDistro(&flags.ConnectionDetails, &distro); err != nil {
			return err
		}
	}
	return nil
}

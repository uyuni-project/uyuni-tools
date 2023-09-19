package distro

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func umountAndRemove(mountpoint string) {
	umount_cmd := []string{
		"/usr/bin/umount",
		mountpoint,
	}

	if err := utils.RunRawCmd("/usr/bin/sudo", umount_cmd, true); err != nil {
		log.Fatal().Err(err).Msgf("Unable to unmount iso file, leaving %s intact", mountpoint)
	}

	os.Remove(mountpoint)
}

func distCp(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, distroName string, source string) {
	log.Info().Msgf("Copying distribution %s\n", distroName)
	if !utils.FileExists(source) {
		log.Fatal().Msgf("Source %s does not exists", source)
	}

	dstpath := "/srv/www/distributions/" + distroName
	if utils.TestExistenceInPod(globalFlags, flags.Backend, dstpath) {
		log.Fatal().Msgf("Distribution already exists: %s\n", dstpath)
	}

	srcdir := source
	if strings.HasSuffix(source, ".iso") {
		log.Debug().Msg("Source is an iso file")
		tmpdir, err := os.MkdirTemp("", "uyunictl")
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
		if out, err := utils.RunCmdOutput("/usr/bin/sudo", mount_cmd...); err != nil {
			log.Debug().Msgf("output %s", out)
			log.Error().Err(err).Msg("Unable to mount iso file. Mount manually and try again")
			return
		}
	}

	utils.Copy(globalFlags, flags.Backend, srcdir, "server:"+dstpath, "tomcat", "susemanager")

	log.Info().Msg("Distribution has been copied")
}

package distro

import (
	"errors"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	Backend string
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	flags := &flagpole{}

	distroCmd := &cobra.Command{
		Use:     "distribution",
		Short:   "Distribution management",
		Long:    "Tools and utilities for distribution management",
		Aliases: []string{"distro"},
	}

	cpCmd := &cobra.Command{
		Use:   "copy [path/to/source] [distribution name]",
		Short: "copy distribution files from iso to the container",
		Long: `takes a path to iso file or directory with mounted iso and copies it into the container.
	Distribution name specifies the destination directory under /srv/www/distributions.`,
		Args:    cobra.ExactArgs(2),
		Aliases: []string{"cp"},
		Run: func(cmd *cobra.Command, args []string) {
			viper := utils.ReadConfig(globalFlags.ConfigPath, "ctlconfig", cmd)
			if err := viper.Unmarshal(&flags); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshall configuration")
			}
			distCp(globalFlags, flags, cmd, args[1], args[0])
		},
	}

	distroCmd.AddCommand(cpCmd)
	return distroCmd
}

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
	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		log.Fatal().Err(err).Msgf("Source %s does not exists", source)
	}

	dstpath := "/srv/www/distributions/" + distroName
	if utils.TestExistence(globalFlags, flags.Backend, dstpath) {
		log.Fatal().Msgf("Distribution already exists: %s\n", dstpath)
	}

	srcdir := source
	if strings.HasSuffix(source, ".iso") {
		log.Debug().Msg("Source is an iso file")
		tmpdir, err := os.MkdirTemp("", "uyuni-tools")
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
		if err := utils.RunRawCmd("/usr/bin/sudo", mount_cmd, true); err != nil {
			log.Fatal().Err(err).Msg("Unable to mount iso file. Mount manually and try again")
		}
	}

	utils.Copy(globalFlags, flags.Backend, srcdir, "server:"+dstpath, "tomcat", "susemanager")

	log.Info().Msg("Distribution copied")
}

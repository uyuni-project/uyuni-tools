package distro

import (
	"errors"
	"log"
	"os"
	"strings"

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
				log.Fatalf("Failed to unmarshall configuration: %s\n", err)
			}
			distCp(globalFlags, flags, cmd, args[1], args[0])
		},
	}

	distroCmd.AddCommand(cpCmd)
	return distroCmd
}

func distCp(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, distroName string, source string) {
	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Source %s does not exists", source)
	}

	srcdir := source
	if strings.HasSuffix(source, ".iso") {
		srcdir, err := os.MkdirTemp("", "uyuni-tools")
		if err != nil {
			log.Fatal(err)
		}
		mount_cmd := []string{
			"/usr/bin/mount",
			"-o loop",
			source,
			srcdir,
		}
		utils.RunCmd("/usr/bin/sudo", mount_cmd, "Unable to mount iso file. Mount manually and try again", globalFlags.Verbose)
		defer os.Remove(srcdir)
	}

	dstpath := "/srv/www/distributions/" + distroName
	if utils.TestExistence(globalFlags, flags.Backend, dstpath) {
		log.Fatalf("Distribution already exists: %s\n", dstpath)
	}

	utils.Copy(globalFlags, flags.Backend, srcdir, "server:"+dstpath, "tomcat", "susemanager")
}

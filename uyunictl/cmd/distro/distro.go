package distro

import (
	"log"

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
			distCp(globalFlags, flags, cmd, args)
		},
	}

	distroCmd.AddCommand(cpCmd)
	return distroCmd
}

func distCp(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, args []string) {
	dstpath := "/srv/www/distributions/" + args[1]
	if utils.TestExistence(globalFlags, flags.Backend, dstpath) {
		log.Fatalf("Distribution already exists: %s\n", dstpath)
	}
	utils.Copy(globalFlags, flags.Backend, args[0], "server:"+dstpath, "tomcat", "susemanager")
}

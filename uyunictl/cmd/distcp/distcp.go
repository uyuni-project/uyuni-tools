package distcp

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	User  string
	Group string
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	flags := &flagpole{}

	cpCmd := &cobra.Command{
		Use:   "distcp [path/to/mounted/iso] [distribution name]",
		Short: "copy distribution files from monted iso to container",
		Long: `distcp takes a path to mounted iso and copies it into the container.
	Distribution name specifies the destination diretory under /srv/www/distributions.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			run(globalFlags, flags, cmd, args)
		},
	}

	cpCmd.Flags().StringVar(&flags.User, "user", "", "User or UID to set on the destination files")
	cpCmd.Flags().StringVar(&flags.Group, "group", "", "Group or GID to set on the destination files")
	return cpCmd
}

func run(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, args []string) {
	dstpath := "/srv/www/distributions/" + args[1]
	if utils.TestExistence(dstpath) {
		log.Fatalf("Distribution already exists: %s\n", dstpath)
	}
	utils.Copy(globalFlags, args[0], "server:" + dstpath, flags.User, flags.Group)
}

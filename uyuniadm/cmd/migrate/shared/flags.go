package shared

import (
	"github.com/spf13/cobra"
	cmd_utils "github.com/uyuni-project/uyuni-tools/uyuniadm/shared/utils"
)

type MigrateFlags struct {
	Image cmd_utils.ImageFlags `mapstructure:",squash"`
}

func AddMigrateFlags(cmd *cobra.Command) {
	cmd_utils.AddImageFlag(cmd)
}

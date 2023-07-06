package migrate

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func migrateToKubernetes(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, args []string) {
	// TODO Run the migration job

	// TODO prepare the values.yaml and deploy helm chart
}

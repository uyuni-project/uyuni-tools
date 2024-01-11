package podman

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/shared"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	podman_utils "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type podmanUpgradeFlags struct {
	shared.UpgradeFlags `mapstructure:",squash"`
	Podman              podman_utils.PodmanFlags
	MirrorPath          string
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	upgradeCmd := &cobra.Command{
		Use:   "podman",
		Short: "upgrade a local server on podman",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags podmanUpgradeFlags
			if err := viper.Unmarshal(&flags); err != nil {
				log.Fatal().Err(err).Msg("Failed to Unmarshal configuration")
			}

			upgradePodman(globalFlags, &flags, cmd, args)
		},
	}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list available tag for an image",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			viper, _ := utils.ReadConfig(globalFlags.ConfigPath, cmd)

			var flags podmanUpgradeFlags
			if err := viper.Unmarshal(&flags); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshall configuration")
			}
			tags := podman_utils.ShowAvailableTag(flags.Image.Name)
			log.Info().Msgf("Available Tags for image: %s", flags.Image.Name)
			for _, value := range tags {
				log.Info().Msgf("%s", value)
			}
		},
	}
	shared.AddUpgradeListFlags(listCmd)
	upgradeCmd.AddCommand(listCmd)

	shared.AddUpgradeFlags(upgradeCmd)
	podman.AddPodmanInstallFlag(upgradeCmd)

	return upgradeCmd
}

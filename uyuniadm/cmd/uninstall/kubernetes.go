package uninstall

import (
	"log"

	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func uninstallForKubernetes(globalFlags *types.GlobalFlags, dryRun bool) {
	// Run helm uninstall
	if dryRun {
		log.Println("Would run helm uninstall uyuni")
	} else {
		utils.RunCmd("helm", []string{"uninstall", "uyuni"}, "Failed to run helm uninstall uyuni", globalFlags.Verbose)
	}
}

package migrate

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/templates"
)

func migrateToKubernetes(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, args []string) {
	scriptDir := generateMigrationScript(args[0], true)
	defer os.RemoveAll(scriptDir)

	runMigrationJob(scriptDir, flags, globalFlags.Verbose)

	// TODO Watch the logs and wait for the end of the job

	// TODO prepare the values.yaml and deploy helm chart
}

func runMigrationJob(tmpPath string, flags *flagpole, verbose bool) {
	sshAuthSocket := getSshAuthSocket()

	// Find ssh config to mount it in the container
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to find home directory to look for SSH config")
	}
	sshConfigPath := filepath.Join(homedir, ".ssh", "config")
	sshKnownhostsPath := filepath.Join(homedir, ".ssh", "known_hosts")

	volumes := make(map[string]templates.Volume)
	for name, path := range utils.VOLUMES {
		volumes[name] = templates.Volume{Path: path}
	}
	volumes["ssh-auth-socket"] = templates.Volume{HostPath: sshAuthSocket, Path: "/tmp/ssh_auth_sock"}
	volumes["ssh-config"] = templates.Volume{HostPath: sshConfigPath, Path: "/root/.ssh/config"}
	volumes["ssh-known-hosts"] = templates.Volume{HostPath: sshKnownhostsPath, Path: "/root/.ssh/known_hosts"}

	templateData := templates.KubernetesMigrateJobTemplateData{
		Volumes: volumes,
		Image:   flags.Image,
		Tag:     flags.ImageTag,
	}

	// TODO PVCs and PVs need to be ready before this

	migrationYamlPath := filepath.Join(tmpPath, "migration-job.yaml")
	if err = utils.WriteTemplateToFile(templateData, migrationYamlPath, 0600, true); err != nil {
		log.Fatalf("Failed to generate migration job description: %s\n", err)
	}

	utils.RunCmd("kubectl", []string{"apply", "-f", migrationYamlPath}, "Failed to start migration job", verbose)
}

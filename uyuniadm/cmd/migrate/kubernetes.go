package migrate

import (
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func migrateToKubernetes(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, args []string) {
	scriptDir := generateMigrationScript(args[0], true)
	defer os.RemoveAll(scriptDir)

	runMigrationJob(scriptDir, flags, globalFlags.Verbose)

	// TODO Watch the logs and wait for the end of the job

	// TODO prepare the values.yaml and deploy helm chart
}

type volume struct {
	HostPath string
	Path     string
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

	volumes := make(map[string]volume)
	for name, path := range utils.VOLUMES {
		volumes[name] = volume{Path: path}
	}
	volumes["ssh-auth-socket"] = volume{HostPath: sshAuthSocket, Path: "/tmp/ssh_auth_sock"}
	volumes["ssh-config"] = volume{HostPath: sshConfigPath, Path: "/root/.ssh/config"}
	volumes["ssh-known-hosts"] = volume{HostPath: sshKnownhostsPath, Path: "/root/.ssh/known_hosts"}

	model := struct {
		Volumes map[string]volume
		Image   string
		Tag     string
	}{
		Volumes: volumes,
		Image:   flags.Image,
		Tag:     flags.ImageTag,
	}

	// TODO PVCs and PVs need to be ready before this

	migrationYamlPath := filepath.Join(tmpPath, "migration-job.yaml")
	file, err := os.OpenFile(migrationYamlPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Fail to open migration job definition file: %s\n", err)
	}
	defer file.Close()

	t := template.Must(template.New("job").Parse(migrationJob))
	if err = t.Execute(file, model); err != nil {
		log.Fatalf("Failed to generate migration job description: %s\n", err)
	}

	utils.RunCmd("kubectl", []string{"apply", "-f", migrationYamlPath}, "Failed to start migration job", verbose)
}

const migrationJob = `apiVersion: batch/v1
kind: Job
metadata:
  name: uyuni-migration
spec:
  backoffLimit: 1
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: rsync-var-pgsql
        image: {{ .Image }}:{{ .Tag }}
        command: [ "/var/lib/uyuni-tools/migrate.sh" ]
        env:
          - name: SSH_AUTH_SOCK
            value: /tmp/ssh_auth_sock
        volumeMounts:
        {{- range $name, $volume := .Volumes }}
          - mountPath: {{ $volume.Path }}
            name: {{ $name }}
        {{- end }}
      volumes:
	  {{- range $name, $volume := .Volumes }}
        - name: {{ $name }}
		{{- if eq $volume.HostPath "" }}
          persistentVolumeClaim:
            claimName: {{ $name }}
		{{- else }}
          hostPath:
            path: {{ $volume.HostPath }}
	    {{- end }}
      {{- end }}
`

package install

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	image       string
	imageTag    string
	timeZone    string
	manager     manager
	reportDB    database
	cert        certificate
	enableTftp  bool
	sccUser     string
	sccPassword string
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	flags := &flagpole{}

	installCmd := &cobra.Command{
		Use:   "install [fqdn]",
		Short: "install a new server from scratch",
		Long: `Install a new server from scratch

The install command assumes the following:
  * podman or kubectl is installed locally
  * if kubectl is installed, a working kubeconfig should be set to connect to the cluster to deploy to

When installing on kubernetes, the helm values file will be overridden with the values from the uyuniadm parameters or configuration.

NOTE: for now installing on a remote cluster or podman is not supported!
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viper := utils.ReadConfig(globalFlags.ConfigPath, "admconfig", cmd)
			checkParameters(viper, flags)
			command := utils.GetCommand()
			switch command {
			case "podman":
				installForPodman(viper, globalFlags, cmd, args)
			case "kubectl":
				installForKubernetes(viper, globalFlags, cmd, args)
			}
		},
	}

	installCmd.Flags().String("image", "registry.opensuse.org/uyuni/server", "Image")
	installCmd.Flags().String("tag", "latest", "Tag Image")

	installCmd.Flags().String("tz", "Etc/UTC", "Time zone to set on the server. Defaults to the host timezone")
	installCmd.Flags().String("email", "admin@example.com", "Manager E-Mail")
	installCmd.Flags().String("emailfrom", "admin@example.com", "E-Mail sending the notifications")
	installCmd.Flags().String("mirrorPath", "", "Path to mirrored packages mounted on the host")
	installCmd.Flags().String("issParent", "", "Inter Server Sync v1 parent fully qualified domain name")
	installCmd.Flags().String("db-user", "spacewalk", "Manager User")
	installCmd.Flags().String("db-password", "", "Manager Password")
	installCmd.Flags().String("db-name", "susemanager", "Database Name")
	installCmd.Flags().String("db-host", "localhost", "Database Host")
	installCmd.Flags().Int("db-port", 5432, "Database Port")
	installCmd.Flags().String("db-protocol", "tcp", "Database Protocol")
	installCmd.Flags().String("db-admin-user", "", "External database admin user name")
	installCmd.Flags().String("db-admin-password", "", "External database admin password")
	installCmd.Flags().String("db-provider", "", "External database provider. Possible values 'aws'")

	installCmd.Flags().Bool("tftp", true, "Enable TFTP")
	installCmd.Flags().String("reportdb-name", "reportdb", "Report Database Name")
	installCmd.Flags().String("reportdb-host", "", "Report Database Host. Defaults to the selected FQDN")
	installCmd.Flags().Int("reportdb-port", 5432, "Report Database Port")
	installCmd.Flags().String("reportdb-user", "pythia_susemanager", "Report Database username")
	installCmd.Flags().String("reportdb-password", "", "Report Database Password")

	installCmd.Flags().Bool("cert-useexisting", false, "Use Existing Certificate")
	installCmd.Flags().StringArray("cert-cname", []string{}, "SSL Certificate cnames separated by commas")
	installCmd.Flags().String("cert-country", "DE", "SSL Certificate Country")
	installCmd.Flags().String("cert-state", "Bayern", "SSL Certificate State")
	installCmd.Flags().String("cert-city", "Nuernberg", "SSL Certificate City")
	installCmd.Flags().String("cert-org", "SUSE", "SSL Certificate Organization")
	installCmd.Flags().String("cert-ou", "SUSE", "SSL Certificate Organization Unit")
	installCmd.Flags().String("cert-password", "", "SSL Certificate Password")
	installCmd.Flags().String("cert-email", "ca-admin@example.com", "SSL Certificate E-Mail")

	installCmd.Flags().String("scc-user", "", "SUSE Customer Center username")
	installCmd.Flags().String("scc-password", "", "SUSE Customer Center password")

	installCmd.Flags().StringArray("podman-arg", []string{}, "Extra arguments to pass to podman")

	installCmd.Flags().String("helm-namespace", "default", "Kubernetes namespace to install uyuni to")
	installCmd.Flags().String("helm-chart", "oci://registry.opensuse.org/uyuni/proxy", "URL to the uyuni helm chart")
	installCmd.Flags().String("helm-values", "", "Path to a values YAML file to use for helm install")

	return installCmd
}

func checkParameters(viper *viper.Viper, flags *flagpole) {
	utils.AskPasswordIfMissing(viper, "db.password", "Database user password: ")
	if !flags.cert.useExistingCertificate {
		utils.AskPasswordIfMissing(viper, "cert.password", "Password for the CA certificate to generate: ")
	}

	// Use the host timezone if the user didn't define one
	if viper.GetString("tz") == "" {
		viper.Set("tz", utils.GetLocalTimezone())
	}
}

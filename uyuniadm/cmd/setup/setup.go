package setup

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

	setupCmd := &cobra.Command{
		Use:   "setup [fqdn]",
		Short: "setup a new server from scratch",
		Long: `Setup a new server from scratch

The setup command assumes the following:
  * podman or kubectl is installed locally
  * if kubectl is installed, a working kubeconfig should be set to connect to the cluster to deploy to

NOTE: for now installing on a remote cluster or podman is not supported yet!
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viper := utils.ReadConfig(globalFlags.ConfigPath, "admconfig", cmd)
			checkParameters(viper, flags)
			command := utils.GetCommand()
			switch command {
			case "podman":
				setupForPodman(viper, globalFlags, cmd, args)
			case "kubectl":
				setupForKubernetes(viper, globalFlags, cmd, args)
			}
		},
	}

	setupCmd.Flags().String("image", "registry.opensuse.org/uyuni/server", "Image")
	setupCmd.Flags().String("tag", "latest", "Tag Image")

	setupCmd.Flags().String("tz", "Etc/UTC", "Time zone to set on the server. Defaults to the host timezone")
	setupCmd.Flags().String("email", "admin@example.com", "Manager E-Mail")
	setupCmd.Flags().String("emailfrom", "admin@example.com", "E-Mail sending the notifications")
	setupCmd.Flags().String("mirrorPath", "", "Path to mirrored packages mounted on the host")
	setupCmd.Flags().String("issParent", "", "Inter Server Sync v1 parent fully qualified domain name")
	setupCmd.Flags().String("db-user", "spacewalk", "Manager User")
	setupCmd.Flags().String("db-password", "", "Manager Password")
	setupCmd.Flags().String("db-name", "susemanager", "Database Name")
	setupCmd.Flags().String("db-host", "localhost", "Database Host")
	setupCmd.Flags().Int("db-port", 5432, "Database Port")
	setupCmd.Flags().String("db-protocol", "tcp", "Database Protocol")
	setupCmd.Flags().String("db-admin-user", "", "External database admin user name")
	setupCmd.Flags().String("db-admin-password", "", "External database admin password")
	setupCmd.Flags().String("db-provider", "", "External database provider. Possible values 'aws'")

	setupCmd.Flags().Bool("tftp", true, "Enable TFTP")
	setupCmd.Flags().String("reportdb-name", "reportdb", "Report Database Name")
	setupCmd.Flags().String("reportdb-host", "", "Report Database Host. Defaults to the selected FQDN")
	setupCmd.Flags().Int("reportdb-port", 5432, "Report Database Port")
	setupCmd.Flags().String("reportdb-user", "pythia_susemanager", "Report Database username")
	setupCmd.Flags().String("reportdb-password", "", "Report Database Password")

	setupCmd.Flags().Bool("cert-useexisting", false, "Use Existing Certificate")
	setupCmd.Flags().StringArray("cert-cname", []string{}, "SSL Certificate cnames separated by commas")
	setupCmd.Flags().String("cert-country", "DE", "SSL Certificate Country")
	setupCmd.Flags().String("cert-state", "Bayern", "SSL Certificate State")
	setupCmd.Flags().String("cert-city", "Nuernberg", "SSL Certificate City")
	setupCmd.Flags().String("cert-org", "SUSE", "SSL Certificate Organization")
	setupCmd.Flags().String("cert-ou", "SUSE", "SSL Certificate Organization Unit")
	setupCmd.Flags().String("cert-password", "", "SSL Certificate Password")
	setupCmd.Flags().String("cert-email", "ca-admin@example.com", "SSL Certificate E-Mail")

	setupCmd.Flags().String("scc-user", "", "SUSE Customer Center username")
	setupCmd.Flags().String("scc-password", "", "SUSE Customer Center password")

	return setupCmd
}

func checkParameters(viper *viper.Viper, flags *flagpole) {
	utils.AskPasswordIfMissing(viper, "db.password", "Database user password: ")
	if !flags.cert.useExistingCertificate {
		utils.AskPasswordIfMissing(viper, "cert.password", "Password for the CA certificate to generate: ")
	}
}

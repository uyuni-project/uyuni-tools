package install

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	cmd_utils "github.com/uyuni-project/uyuni-tools/uyuniadm/shared/utils"
)

type DbFlags struct {
	Host     string
	Name     string
	Port     int
	User     string
	Password string
	Protocol string
	Provider string
	Admin    struct {
		User     string
		Password string
	}
}

type SslCertFlags struct {
	UseExisting bool
	Cnames      []string `mapstructure:"cname"`
	Country     string
	State       string
	City        string
	Org         string
	OU          string
	Password    string
	Email       string
}

type SccFlags struct {
	User     string
	Password string
}

type InstallFlags struct {
	TZ         string
	Email      string
	EmailFrom  string
	IssParent  string
	MirrorPath string
	Tftp       bool
	Db         DbFlags
	ReportDb   DbFlags
	Cert       SslCertFlags
	Scc        SccFlags
	Image      cmd_utils.ImageFlags `mapstructure:",squash"`
	Podman     cmd_utils.PodmanFlags
	Helm       cmd_utils.HelmFlags
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

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
			var flags InstallFlags
			if err := viper.Unmarshal(&flags); err != nil {
				log.Fatalf("Failed to unmarshall configuration: %s\n", err)
			}
			command := utils.GetCommand("")
			checkParameters(cmd, &flags, command)
			switch command {
			case "podman":
				installForPodman(globalFlags, &flags, cmd, args)
			case "kubectl":
				installForKubernetes(globalFlags, &flags, cmd, args)
			}
		},
	}

	installCmd.Flags().String("tz", "Etc/UTC", "Time zone to set on the server. Defaults to the host timezone")
	installCmd.Flags().String("email", "admin@example.com", "Administrator e-mail")
	installCmd.Flags().String("emailfrom", "admin@example.com", "E-Mail sending the notifications")
	installCmd.Flags().String("mirrorPath", "", "Path to mirrored packages mounted on the host")
	installCmd.Flags().String("issParent", "", "Inter Server Sync v1 parent fully qualified domain name")
	installCmd.Flags().String("db-user", "spacewalk", "Database user")
	installCmd.Flags().String("db-password", "", "Database password")
	installCmd.Flags().String("db-name", "susemanager", "Database name")
	installCmd.Flags().String("db-host", "localhost", "Database host")
	installCmd.Flags().Int("db-port", 5432, "Database port")
	installCmd.Flags().String("db-protocol", "tcp", "Database protocol")
	installCmd.Flags().String("db-admin-user", "", "External database admin user name")
	installCmd.Flags().String("db-admin-password", "", "External database admin password")
	installCmd.Flags().String("db-provider", "", "External database provider. Possible values 'aws'")

	installCmd.Flags().Bool("tftp", true, "Enable TFTP")
	installCmd.Flags().String("reportdb-name", "reportdb", "Report database name")
	installCmd.Flags().String("reportdb-host", "", "Report database host. Defaults to the selected FQDN")
	installCmd.Flags().Int("reportdb-port", 5432, "Report database port")
	installCmd.Flags().String("reportdb-user", "pythia_susemanager", "Report Database username")
	installCmd.Flags().String("reportdb-password", "", "Report database password. Randomly generated by default")

	installCmd.Flags().Bool("cert-useexisting", false, "Use existing SSL certificate")
	installCmd.Flags().StringSlice("cert-cname", []string{}, "SSL certificate cnames separated by commas")
	installCmd.Flags().String("cert-country", "DE", "SSL certificate country")
	installCmd.Flags().String("cert-state", "Bayern", "SSL certificate state")
	installCmd.Flags().String("cert-city", "Nuernberg", "SSL certificate city")
	installCmd.Flags().String("cert-org", "SUSE", "SSL certificate organization")
	installCmd.Flags().String("cert-ou", "SUSE", "SSL certificate organization unit")
	installCmd.Flags().String("cert-password", "", "Password for the CA certificate to generate")
	installCmd.Flags().String("cert-email", "ca-admin@example.com", "SSL certificate E-Mail")

	installCmd.Flags().String("scc-user", "", "SUSE Customer Center username")
	installCmd.Flags().String("scc-password", "", "SUSE Customer Center password")

	cmd_utils.AddImageFlag(installCmd)
	cmd_utils.AddPodmanInstallFlag(installCmd)
	cmd_utils.AddHelmInstallFlag(installCmd)

	return installCmd
}

func checkParameters(cmd *cobra.Command, flags *InstallFlags, command string) {
	utils.AskPasswordIfMissing(&flags.Db.Password, cmd.Flag("db-password").Usage)

	// Since we use cert-manager for self-signed certificates on kubernetes we don't need password for it
	if !flags.Cert.UseExisting && command == "podman" {
		utils.AskPasswordIfMissing(&flags.Cert.Password, cmd.Flag("cert-password").Usage)
	}

	// Use the host timezone if the user didn't define one
	if flags.TZ == "" {
		flags.TZ = utils.GetLocalTimezone()
	}

	utils.AskIfMissing(&flags.Email, cmd.Flag("email").Usage)
	utils.AskIfMissing(&flags.EmailFrom, cmd.Flag("emailfrom").Usage)
}

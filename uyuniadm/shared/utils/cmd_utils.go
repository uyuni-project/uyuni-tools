package utils

import (
	"github.com/spf13/cobra"
)

var DefaultImage = "registry.opensuse.org/uyuni/server"
var DefaultTag = "latest"
var DefaultChart = "oci://registry.opensuse.org/uyuni/server"

type PodmanFlags struct {
	Args []string `mapstructure:"arg"`
}

type ChartFlags struct {
	Namespace string
	Chart     string
	Version   string
	Values    string
}

type HelmFlags struct {
	Uyuni       ChartFlags
	CertManager ChartFlags
}

type ImageFlags struct {
	Name string `mapstructure:"image"`
	Tag  string
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

func AddPodmanInstallFlag(cmd *cobra.Command) {
	cmd.Flags().StringSlice("podman-arg", []string{}, "Extra arguments to pass to podman")
}

func AddHelmInstallFlag(cmd *cobra.Command) {
	cmd.Flags().String("helm-uyuni-namespace", "default", "Kubernetes namespace where to install uyuni")
	cmd.Flags().String("helm-uyuni-chart", DefaultChart, "URL to the uyuni helm chart")
	cmd.Flags().String("helm-uyuni-version", "", "Version of the uyuni helm chart")
	cmd.Flags().String("helm-uyuni-values", "", "Path to a values YAML file to use for Uyuni helm install")
	cmd.Flags().String("helm-certmanager-namespace", "cert-manager", "Kubernetes namespace where to install cert-manager")
	cmd.Flags().String("helm-certmanager-chart", "", "URL to the cert-manager helm chart. To be used for offline installations")
	cmd.Flags().String("helm-certmanager-version", "", "Version of the cert-manager helm chart")
	cmd.Flags().String("helm-certmanager-values", "", "Path to a values YAML file to use for cert-manager helm install")
}

func AddImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("image", DefaultImage, "Image")
	cmd.Flags().String("tag", DefaultTag, "Tag Image")
}

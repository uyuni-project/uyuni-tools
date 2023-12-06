// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var DefaultImage = path.Join(utils.DefaultNamespace, "server")

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

type SslCertFlags struct {
	Cnames   []string `mapstructure:"cname"`
	Country  string
	State    string
	City     string
	Org      string
	OU       string
	Password string
	Email    string
	Ca       ssl.CaChain
	Server   ssl.SslPair
}

func (f *SslCertFlags) UseExisting() bool {
	return f.Server.Cert != "" && f.Server.Key != "" && f.Ca.Root != ""
}

// Checks that all the required flags are passed if using 3rd party certificates.
func (f *SslCertFlags) CheckParameters() {
	if !f.UseExisting() && (f.Server.Cert != "" || f.Server.Key != "" || f.Ca.Root != "") {
		log.Fatal().Msg("Server certificate, key and root CA need to be all provided")
	}
}

func AddHelmInstallFlag(cmd *cobra.Command) {
	defaultChart := fmt.Sprintf("oci://%s/server-helm", utils.DefaultNamespace)

	cmd.Flags().String("helm-uyuni-namespace", "default", "Kubernetes namespace where to install uyuni")
	cmd.Flags().String("helm-uyuni-chart", defaultChart, "URL to the uyuni helm chart")
	cmd.Flags().String("helm-uyuni-version", "", "Version of the uyuni helm chart")
	cmd.Flags().String("helm-uyuni-values", "", "Path to a values YAML file to use for Uyuni helm install")
	cmd.Flags().String("helm-certmanager-namespace", "cert-manager", "Kubernetes namespace where to install cert-manager")
	cmd.Flags().String("helm-certmanager-chart", "", "URL to the cert-manager helm chart. To be used for offline installations")
	cmd.Flags().String("helm-certmanager-version", "", "Version of the cert-manager helm chart")
	cmd.Flags().String("helm-certmanager-values", "", "Path to a values YAML file to use for cert-manager helm install")
}

func AddImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("image", DefaultImage, "Image")
	cmd.Flags().String("tag", utils.DefaultTag, "Tag Image")

	// Podman:
	//   Never, just check and fail if needed
	//   IfNotPresent, check and pull
	//   Always, pull without checking
	// Kubernetes -> set helm values
	cmd.Flags().String("pullPolicy", "IfNotPresent",
		"set whether to pull the images or not. The value can be one of 'Never', 'IfNotPresent' or 'Always'")
}

func AddMigrationImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("migration-image", "", "Migration image")
	cmd.Flags().String("migration-tag", utils.DefaultTag, "Migration image tag")
}

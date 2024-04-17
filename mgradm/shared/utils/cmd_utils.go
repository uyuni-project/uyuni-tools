// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var defaultImage = path.Join(utils.DefaultNamespace, "server")

// HelmFrags stores Uyuni and Cert Manager Helm information.
type HelmFlags struct {
	Uyuni       types.ChartFlags
	CertManager types.ChartFlags
}

// SslCertFlags can store SSL Certs information.
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

// UseExisting return true if existing SSL Cert can be used.
func (f *SslCertFlags) UseExisting() bool {
	return f.Server.Cert != "" && f.Server.Key != "" && f.Ca.Root != ""
}

// Checks that all the required flags are passed if using 3rd party certificates.
func (f *SslCertFlags) CheckParameters() {
	if !f.UseExisting() && (f.Server.Cert != "" || f.Server.Key != "" || f.Ca.Root != "") {
		log.Fatal().Msg(L("Server certificate, key and root CA need to be all provided"))
	}
}

// AddHelmInstallFlag add Helm install flags to a command.
func AddHelmInstallFlag(cmd *cobra.Command) {
	defaultChart := fmt.Sprintf("oci://%s/server-helm", utils.DefaultNamespace)

	cmd.Flags().String("helm-uyuni-namespace", "default", L("Kubernetes namespace where to install uyuni"))
	cmd.Flags().String("helm-uyuni-chart", defaultChart, L("URL to the uyuni helm chart"))
	cmd.Flags().String("helm-uyuni-version", "", L("Version of the uyuni helm chart"))
	cmd.Flags().String("helm-uyuni-values", "", L("Path to a values YAML file to use for Uyuni helm install"))
	cmd.Flags().String("helm-certmanager-namespace", "cert-manager", L("Kubernetes namespace where to install cert-manager"))
	cmd.Flags().String("helm-certmanager-chart", "", L("URL to the cert-manager helm chart. To be used for offline installations"))
	cmd.Flags().String("helm-certmanager-version", "", L("Version of the cert-manager helm chart"))
	cmd.Flags().String("helm-certmanager-values", "", L("Path to a values YAML file to use for cert-manager helm install"))
}

// AddimageFlag add Image flags to a command.
func AddImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("image", defaultImage, L("Image"))
	cmd.Flags().String("tag", utils.DefaultTag, L("Tag Image"))

	utils.AddPullPolicyFlag(cmd)
}

// AddMigrationImageFlag add Migration Image flags to a command.
func AddMigrationImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("migration-image", "", L("Migration image"))
	cmd.Flags().String("migration-tag", utils.DefaultTag, L("Migration image tag"))
	cmd.Flags().String("migration-pullPolicy", "IfNotPresent",
		L("set whether to pull the migrattion images or not. The value can be one of 'Never', 'IfNotPresent' or 'Always'"))
}

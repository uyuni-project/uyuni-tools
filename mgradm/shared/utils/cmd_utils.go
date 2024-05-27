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

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "helm", Title: L("Helm Chart Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "helm-uyuni-namespace", "helm")
	_ = utils.AddFlagToHelpGroupID(cmd, "helm-uyuni-chart", "helm")
	_ = utils.AddFlagToHelpGroupID(cmd, "helm-uyuni-version", "helm")
	_ = utils.AddFlagToHelpGroupID(cmd, "helm-uyuni-values", "helm")
	_ = utils.AddFlagToHelpGroupID(cmd, "helm-certmanager-namespace", "helm")
	_ = utils.AddFlagToHelpGroupID(cmd, "helm-certmanager-chart", "helm")
	_ = utils.AddFlagToHelpGroupID(cmd, "helm-certmanager-version", "helm")
	_ = utils.AddFlagToHelpGroupID(cmd, "helm-certmanager-values", "helm")
}

// AddContainerImageFlags add container image flags to command.
func AddContainerImageFlags(cmd *cobra.Command, container string, displayName string) {
	cmd.Flags().String(container+"-image", "",
		fmt.Sprintf(L("Image for %s container, overrides the namespace if set"), displayName))
	cmd.Flags().String(container+"-tag", "",
		fmt.Sprintf(L("Tag for %s container, overrides the global value if set"), displayName))
}

// AddImageFlag add Image flags to a command.
func AddImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("image", defaultImage, L("Image"))
	cmd.Flags().String("tag", utils.DefaultTag, L("Tag Image"))

	utils.AddPullPolicyFlag(cmd)

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "image", Title: L("Image Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "image", "image")
	_ = utils.AddFlagToHelpGroupID(cmd, "tag", "image")
	_ = utils.AddFlagToHelpGroupID(cmd, "pullPolicy", "image")
}

// AddDbUpgradeImageFlag add Database upgrade image flags to a command.
func AddDbUpgradeImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("dbupgrade-image", "", L("Database upgrade image"))
	cmd.Flags().String("dbupgrade-tag", "latest", L("Database upgrade image tag"))
	cmd.Flags().String("dbupgrade-pullPolicy", "IfNotPresent",
		L("set whether to pull the database upgrade images or not. The value can be one of 'Never', 'IfNotPresent' or 'Always'"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "dbupgrade-image", Title: L("Database Upgrade Image Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "dbupgrade-image", "dbupgrade-image")
	_ = utils.AddFlagToHelpGroupID(cmd, "dbupgrade-tag", "dbupgrade-image")
	_ = utils.AddFlagToHelpGroupID(cmd, "dbupgrade-pullPolicy", "dbupgrade-image")
}

// AddMirrorFlag adds the flag for the mirror.
func AddMirrorFlag(cmd *cobra.Command) {
	cmd.Flags().String("mirror", "", L("Path to mirrored packages mounted on the host"))
}

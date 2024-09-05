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

var defaultImage = path.Join(utils.DefaultRegistry, "server")

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
	defaultChart := fmt.Sprintf("oci://%s/server-helm", utils.DefaultHelmRegistry)

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
func AddContainerImageFlags(
	cmd *cobra.Command,
	container string,
	displayName string,
	groupName string,
	imageName string,
) {
	defaultImage := path.Join(utils.DefaultRegistry, imageName)
	cmd.Flags().String(container+"-image", defaultImage,
		fmt.Sprintf(L("Image for %s container"), displayName))
	cmd.Flags().String(container+"-tag", "",
		fmt.Sprintf(L("Tag for %s container, overrides the global value if set"), displayName))

	if groupName != "" {
		_ = utils.AddFlagToHelpGroupID(cmd, container+"-image", groupName)
		_ = utils.AddFlagToHelpGroupID(cmd, container+"-tag", groupName)
	}
}

// AddSCCFlag add SCC flags to a command.
func AddSCCFlag(cmd *cobra.Command) {
	cmd.Flags().String("scc-user", "", L("SUSE Customer Center username. It will be used as SCC credentials for products synchronization and to pull images from registry.suse.com"))
	cmd.Flags().String("scc-password", "", L("SUSE Customer Center password. It will be used as SCC credentials for products synchronization and to pull images from registry.suse.com"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "scc", Title: L("SUSE Customer Center Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "scc-user", "scc")
	_ = utils.AddFlagToHelpGroupID(cmd, "scc-password", "scc")
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

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "dbupgrade-image", Title: L("Database Upgrade Image Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "dbupgrade-image", "dbupgrade-image")
	_ = utils.AddFlagToHelpGroupID(cmd, "dbupgrade-tag", "dbupgrade-image")
	_ = utils.AddFlagToHelpGroupID(cmd, "dbupgrade-pullPolicy", "dbupgrade-image")
}

// AddMirrorFlag adds the flag for the mirror.
func AddMirrorFlag(cmd *cobra.Command) {
	cmd.Flags().String("mirror", "", L("Path to mirrored packages mounted on the host"))
}

// AddCocoFlag adds the confidential computing related parameters to cmd.
func AddCocoFlag(cmd *cobra.Command) {
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "coco-container", Title: L("Confidential Computing Flags")})
	AddContainerImageFlags(cmd, "coco", L("confidential computing attestation"), "coco-container", "server-attestation")
	cmd.Flags().Int("coco-replicas", 0, L("How many replicas of the confidential computing container should be started"))
	_ = utils.AddFlagToHelpGroupID(cmd, "coco-replicas", "coco-container")
}

// AddUpgradeCocoFlag adds the confidential computing related parameters to cmd upgrade.
func AddUpgradeCocoFlag(cmd *cobra.Command) {
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "coco-container", Title: L("Confidential Computing Flags")})
	AddContainerImageFlags(cmd, "coco", L("confidential computing attestation"), "coco-container", "server-attestation")
	cmd.Flags().Int("coco-replicas", 0, L("How many replicas of the confidential computing container should be started. Leave it unset if you want to keep the previous number of replicas."))
	_ = utils.AddFlagToHelpGroupID(cmd, "coco-replicas", "coco-container")
}

// AddHubXmlrpcFlags adds hub XML-RPC related parameters to cmd.
func AddHubXmlrpcFlags(cmd *cobra.Command) {
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "hubxmlrpc-container", Title: L("Hub XML-RPC API")})
	AddContainerImageFlags(cmd, "hubxmlrpc", L("Hub XML-RPC API"), "hubxmlrpc-container", "server-hub-xmlrpc-api")
	cmd.Flags().Int("hubxmlrpc-replicas", 0, L("How many replicas of the Hub XML-RPC API service container should be started."))
	_ = utils.AddFlagToHelpGroupID(cmd, "hubxmlrpc-replicas", "hubxmlrpc-container")
}

// AddUpgradeHubXmlrpcFlags adds hub XML-RPC related parameters to cmd upgrade.
func AddUpgradeHubXmlrpcFlags(cmd *cobra.Command) {
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "hubxmlrpc-container", Title: L("Hub XML-RPC API")})
	AddContainerImageFlags(cmd, "hubxmlrpc", L("Hub XML-RPC API"), "hubxmlrpc-container", "server-hub-xmlrpc-api")
	cmd.Flags().Int("hubxmlrpc-replicas", 0, L("How many replicas of the Hub XML-RPC API service container should be started. Leave it unset if you want to keep the previous number of replicas."))
	_ = utils.AddFlagToHelpGroupID(cmd, "hubxmlrpc-replicas", "hubxmlrpc-container")
}

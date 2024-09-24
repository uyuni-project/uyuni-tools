// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var defaultImage = path.Join(utils.DefaultRegistry, "server")

// UseExisting return true if existing SSL Cert can be used.
func (f *InstallSSLFlags) UseExisting() bool {
	return f.Server.Cert != "" && f.Server.Key != "" && f.Ca.Root != "" && f.Ca.Key == ""
}

// UseMigratedCa returns true if a migrated CA and key can be used.
func (f *InstallSSLFlags) UseMigratedCa() bool {
	return f.Ca.Root != "" && f.Ca.Key != ""
}

// CheckParameters checks that all the required flags are passed if using 3rd party certificates.
func (f *InstallSSLFlags) CheckParameters() {
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
	cmd.Flags().String("helm-certmanager-namespace", "cert-manager",
		L("Kubernetes namespace where to install cert-manager"),
	)
	cmd.Flags().String("helm-certmanager-chart", "",
		L("URL to the cert-manager helm chart. To be used for offline installations"),
	)
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

const volumesFlagsGroupID = "volumes"

// AddVolumesFlags adds the Kubernetes volumes configuration parameters to the command.
func AddVolumesFlags(cmd *cobra.Command) {
	cmd.Flags().String("volumes-class", "", L("Default storage class for all the volumes"))
	cmd.Flags().String("volumes-mirror", "",
		L("PersistentVolume name to use as a mirror. Empty means no mirror is used"),
	)

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: volumesFlagsGroupID, Title: L("Volumes Configuration Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "volumes-class", volumesFlagsGroupID)
	_ = utils.AddFlagToHelpGroupID(cmd, "volumes-mirror", volumesFlagsGroupID)

	addVolumeFlags(cmd, "database", "var-pgsql", "50Gi")
	addVolumeFlags(cmd, "packages", "var-spacewalk", "100Gi")
	addVolumeFlags(cmd, "www", "srv-www", "100Gi")
	addVolumeFlags(cmd, "cache", "var-cache", "10Gi")
}

func addVolumeFlags(cmd *cobra.Command, name string, volumeName string, size string) {
	sizeName := fmt.Sprintf("volumes-%s-size", name)
	cmd.Flags().String(
		sizeName, size, fmt.Sprintf(L("Requested size for the %s volume"), volumeName),
	)
	_ = utils.AddFlagToHelpGroupID(cmd, sizeName, volumesFlagsGroupID)

	className := fmt.Sprintf("volumes-%s-class", name)
	cmd.Flags().String(
		className, "", fmt.Sprintf(L("Requested storage class for the %s volume"), volumeName),
	)
	_ = utils.AddFlagToHelpGroupID(cmd, className, volumesFlagsGroupID)
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
	cmd.Flags().String("scc-user", "", L(`SUSE Customer Center username.
It will be used as SCC credentials for products synchronization and to pull images from registry.suse.com`))
	cmd.Flags().String("scc-password", "", L(`SUSE Customer Center password.
It will be used as SCC credentials for products synchronization and to pull images from registry.suse.com`))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "scc", Title: L("SUSE Customer Center Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "scc-user", "scc")
	_ = utils.AddFlagToHelpGroupID(cmd, "scc-password", "scc")
}

// AddImageFlag add Image flags to a command.
func AddImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("image", defaultImage, L("Image"))
	cmd.Flags().String("registry", utils.DefaultRegistry, L("Specify a private registry where pull the images"))
	cmd.Flags().String("tag", utils.DefaultTag, L("Tag Image"))

	utils.AddPullPolicyFlag(cmd)

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "image", Title: L("Image Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "image", "image")
	_ = utils.AddFlagToHelpGroupID(cmd, "registry", "") // without group, since this flag is applied to all the images
	_ = utils.AddFlagToHelpGroupID(cmd, "tag", "image")
	_ = utils.AddFlagToHelpGroupID(cmd, "pullPolicy", "") // without group, since this flag is applied to all the images
}

// AddDBUpgradeImageFlag add Database upgrade image flags to a command.
func AddDBUpgradeImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("dbupgrade-image", "", L("Database upgrade image"))
	cmd.Flags().String("dbupgrade-tag", "latest", L("Database upgrade image tag"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "dbupgrade-image", Title: L("Database Upgrade Image Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "dbupgrade-image", "dbupgrade-image")
	_ = utils.AddFlagToHelpGroupID(cmd, "dbupgrade-tag", "dbupgrade-image")
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
	cmd.Flags().Int("coco-replicas", 0, L(`How many replicas of the confidential computing container should be started.
Leave it unset if you want to keep the previous number of replicas.`))
	_ = utils.AddFlagToHelpGroupID(cmd, "coco-replicas", "coco-container")
}

// AddHubXmlrpcFlags adds hub XML-RPC related parameters to cmd.
func AddHubXmlrpcFlags(cmd *cobra.Command) {
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "hubxmlrpc-container", Title: L("Hub XML-RPC API")})
	AddContainerImageFlags(cmd, "hubxmlrpc", L("Hub XML-RPC API"), "hubxmlrpc-container", "server-hub-xmlrpc-api")
	cmd.Flags().Int("hubxmlrpc-replicas", 0,
		L("How many replicas of the Hub XML-RPC API service container should be started."),
	)
	_ = utils.AddFlagToHelpGroupID(cmd, "hubxmlrpc-replicas", "hubxmlrpc-container")
}

// AddUpgradeHubXmlrpcFlags adds hub XML-RPC related parameters to cmd upgrade.
func AddUpgradeHubXmlrpcFlags(cmd *cobra.Command) {
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "hubxmlrpc-container", Title: L("Hub XML-RPC API")})
	AddContainerImageFlags(cmd, "hubxmlrpc", L("Hub XML-RPC API"), "hubxmlrpc-container", "server-hub-xmlrpc-api")
	cmd.Flags().Int("hubxmlrpc-replicas", 0,
		L(`How many replicas of the Hub XML-RPC API service container should be started.
Leave it unset if you want to keep the previous number of replicas.`))
	_ = utils.AddFlagToHelpGroupID(cmd, "hubxmlrpc-replicas", "hubxmlrpc-container")
}

// AddSalineFlag adds the Saline related parameters to cmd.
func AddSalineFlag(cmd *cobra.Command) {
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "saline-container", Title: L("Saline Flags")})
	AddContainerImageFlags(cmd, "saline", L("Saline"), "saline-container", "server-saline")
	cmd.Flags().Int("saline-replicas", 0, L(`How many replicas of the Saline container should be started
(only 0 or 1 supported for now)`))
	cmd.Flags().Int("saline-port", 8216, L("Saline port (default: 8216)"))
	_ = utils.AddFlagToHelpGroupID(cmd, "saline-replicas", "saline-container")
	_ = utils.AddFlagToHelpGroupID(cmd, "saline-port", "saline-container")
}

// AddUpgradeSalineFlag adds the Saline related parameters to cmd upgrade.
func AddUpgradeSalineFlag(cmd *cobra.Command) {
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "saline-container", Title: L("Saline Flags")})
	AddContainerImageFlags(cmd, "saline", L("Saline"), "saline-container", "server-saline")
	cmd.Flags().Int("saline-replicas", 0, L(`How many replicas of the Saline container should be started.
Leave it unset if you want to keep the previous number of replicas.
(only 0 or 1 supported for now)`))
	cmd.Flags().Int("saline-port", 8216, L("Saline port (default: 8216)"))
	_ = utils.AddFlagToHelpGroupID(cmd, "saline-replicas", "saline-container")
	_ = utils.AddFlagToHelpGroupID(cmd, "saline-port", "saline-container")
}

// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var defaultImage = path.Join(utils.DefaultImagePrefix, "server")

// UseProvided return true if server can use an SSL Cert provided by flags.
func (f *InstallSSLFlags) UseProvided() bool {
	return f.Server.IsDefined() && f.Ca.IsThirdParty()
}

// UseProvidedDB return true if DB can use an SSL Cert provided by flags.
func (f *InstallSSLFlags) UseProvidedDB() bool {
	return f.DB.IsDefined() && f.DB.CA.IsThirdParty()
}

// UseMigratedCa returns true if a migrated CA and key can be used.
func (f *InstallSSLFlags) UseMigratedCa() bool {
	return f.Ca.Root != "" && f.Ca.Key != ""
}

// CheckParameters checks that all the required flags are passed if using 3rd party certificates.
//
// localDB indicates whether the SSL certificates for the database need to be checked.
// Those are not needed for external databases.
func (f *InstallSSLFlags) CheckParameters(localDB bool) {
	if !f.UseProvided() && (f.Server.Cert != "" || f.Server.Key != "" || f.Ca.IsDefined()) {
		log.Fatal().Msg(L("Server certificate, key and root CA need to be all provided"))
	}

	if f.UseProvided() && localDB && !f.DB.IsDefined() {
		log.Fatal().Msg(L("Database certificate and key need to be provided"))
	}
}

// CheckUpgradeParameters checks that all the required flags are passed if using 3rd party certificates.
//
// localDB indicates whether the SSL certificates for the database need to be checked.
// Those are not needed for external databases.
func (f *InstallSSLFlags) CheckUpgradeParameters(localDB bool) {
	if !f.UseProvidedDB() && (f.DB.Cert != "" || f.DB.Key != "" || f.DB.IsDefined()) {
		log.Fatal().Msg(L("DB certificate, key and root CA need to be all provided"))
	}

	if f.UseProvided() && localDB && !f.DB.IsDefined() {
		log.Fatal().Msg(L("Database certificate and key need to be provided"))
	}
}

// AddHelmInstallFlag add Helm install flags to a command.
func AddHelmInstallFlag(cmd *cobra.Command) {
	cmd.Flags().String("kubernetes-uyuni-namespace", "default", L("Kubernetes namespace where to install uyuni"))
	cmd.Flags().String("kubernetes-certmanager-namespace", "cert-manager",
		L("Kubernetes namespace where to install cert-manager"),
	)
	cmd.Flags().String("kubernetes-certmanager-chart", "",
		L("URL to the cert-manager helm chart. To be used for offline installations"),
	)
	cmd.Flags().String("kubernetes-certmanager-version", "", L("Version of the cert-manager helm chart"))
	cmd.Flags().String("kubernetes-certmanager-values", "",
		L("Path to a values YAML file to use for cert-manager helm install"),
	)

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "helm", Title: L("Helm Chart Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "kubernetes-uyuni-namespace", "helm")
	_ = utils.AddFlagToHelpGroupID(cmd, "kubernetes-certmanager-namespace", "helm")
	_ = utils.AddFlagToHelpGroupID(cmd, "kubernetes-certmanager-chart", "helm")
	_ = utils.AddFlagToHelpGroupID(cmd, "kubernetes-certmanager-version", "helm")
	_ = utils.AddFlagToHelpGroupID(cmd, "kubernetes-certmanager-values", "helm")
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
	defaultImage := path.Join(utils.DefaultImagePrefix, imageName)
	cmd.Flags().String(container+"-image", defaultImage,
		fmt.Sprintf(L("Image for %s container"), displayName))
	cmd.Flags().String(container+"-tag", "",
		fmt.Sprintf(L("Tag for %s container, overrides the global value if set"), displayName))

	if groupName != "" {
		_ = utils.AddFlagToHelpGroupID(cmd, container+"-image", groupName)
		_ = utils.AddFlagToHelpGroupID(cmd, container+"-tag", groupName)
	}
}

// addReportDBFlags add ReportDB flags to a command.
func AddReportDBFlags(cmd *cobra.Command) {
	cmd.Flags().String("reportdb-name", "reportdb", L("Report database name"))
	cmd.Flags().String("reportdb-host", "reportdb", L("Report database host"))
	cmd.Flags().Int("reportdb-port", 5432, L("Report database port"))
	cmd.Flags().String("reportdb-user", "pythia_susemanager", L("Report Database username"))
	cmd.Flags().String("reportdb-password", "", L("Report database password. Randomly generated by default"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "reportdb", Title: L("Report DB Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "reportdb-name", "reportdb")
	_ = utils.AddFlagToHelpGroupID(cmd, "reportdb-host", "reportdb")
	_ = utils.AddFlagToHelpGroupID(cmd, "reportdb-port", "reportdb")
	_ = utils.AddFlagToHelpGroupID(cmd, "reportdb-user", "reportdb")
	_ = utils.AddFlagToHelpGroupID(cmd, "reportdb-password", "reportdb")
}

// addDBFlags add DB flags to a command.
func AddDBFlags(cmd *cobra.Command) {
	cmd.Flags().String("db-user", "spacewalk", L("Database user"))
	cmd.Flags().String("db-password", "", L("Database password. Randomly generated by default"))
	cmd.Flags().String("db-name", "susemanager", L("Database name"))
	cmd.Flags().String("db-host", "db", L("Database host"))
	cmd.Flags().Int("db-port", 5432, L("Database port"))
	cmd.Flags().String("db-admin-user", "postgres", L("Database admin user name"))
	cmd.Flags().String("db-admin-password", "", L("Database admin password"))
	cmd.Flags().String("db-provider", "", L("External database provider. Possible values 'aws'"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "db", Title: L("Database Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "db-user", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-password", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-name", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-host", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-port", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-admin-user", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-admin-password", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-provider", "db")
}

// AddSCCFlag add SCC flags to a command.
func AddSCCFlag(cmd *cobra.Command) {
	cmd.Flags().String("scc-user", "", L(`SUSE Customer Center username.
It will be used as SCC credentials for products synchronization and to pull images from SCC registry`))
	cmd.Flags().String("scc-password", "", L(`SUSE Customer Center password.
It will be used as SCC credentials for products synchronization and to pull images from SCC registry`))
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "scc", Title: L("SUSE Customer Center Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "scc-user", "scc")
	_ = utils.AddFlagToHelpGroupID(cmd, "scc-password", "scc")
}

// AddImageFlag add Image flags to a command.
func AddImageFlag(cmd *cobra.Command) {
	cmd.Flags().String("image", defaultImage, L("Image"))
	cmd.Flags().String("tag", utils.DefaultTag, L("Tag Image"))

	utils.AddPullPolicyFlag(cmd)
	utils.AddRegistryFlag(cmd)

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "image", Title: L("Image Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "image", "image")
	_ = utils.AddFlagToHelpGroupID(cmd, "tag", "image")
	// without group, since this flag is applied to all the images
	_ = utils.AddFlagToHelpGroupID(cmd, "pullPolicy", "")
	_ = utils.AddFlagToHelpGroupID(cmd, "registry", "")
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
	AddContainerImageFlags(cmd, "coco", L("Confidential computing attestation"), "coco-container", "server-attestation")
	cmd.Flags().Int("coco-replicas", 0, L("How many replicas of the confidential computing container should be started"))
	_ = utils.AddFlagToHelpGroupID(cmd, "coco-replicas", "coco-container")
}

// AddUpgradeCocoFlag adds the confidential computing related parameters to cmd upgrade.
func AddUpgradeCocoFlag(cmd *cobra.Command) {
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "coco-container", Title: L("Confidential Computing Flags")})
	AddContainerImageFlags(cmd, "coco", L("Confidential computing attestation"), "coco-container", "server-attestation")
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
	cmd.Flags().Int("saline-port", 8216, L("Saline port"))
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
	cmd.Flags().Int("saline-port", 8216, L("Saline port"))
	_ = utils.AddFlagToHelpGroupID(cmd, "saline-replicas", "saline-container")
	_ = utils.AddFlagToHelpGroupID(cmd, "saline-port", "saline-container")
}

// AddPgsqlFlags adds PostgreSQL related parameters to cmd.
func AddPgsqlFlags(cmd *cobra.Command) {
	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "pgsql-container", Title: L("PostgreSQL Database Container Flags")})
	AddContainerImageFlags(cmd, "pgsql", L("PostgreSQL Database"), "pgsql-container", "server-postgresql")
}

// AddServerFlags add flags common to install, upgrade and migration.
func AddServerFlags(cmd *cobra.Command) {
	AddImageFlag(cmd)
	AddSCCFlag(cmd)
	AddPgsqlFlags(cmd)
	AddDBFlags(cmd)
	AddReportDBFlags(cmd)
	ssl.AddSSLGenerationFlags(cmd)
	ssl.AddSSLThirdPartyFlags(cmd)
	ssl.AddSSLDBThirdPartyFlags(cmd)

	cmd.Flags().String("ssl-password", "", L("Password for the CA key to generate"))
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-password", ssl.GeneratedFlagsGroup)
}

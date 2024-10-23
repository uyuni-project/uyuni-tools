// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const migrationDataPvcName = "migration-data"

func migrateToKubernetes(
	globalFlags *types.GlobalFlags,
	flags *kubernetes.KubernetesServerFlags,
	cmd *cobra.Command,
	args []string,
) error {
	namespace := flags.Helm.Uyuni.Namespace

	// Check the for the required SSH key and configuration
	if err := checkSsh(namespace, &flags.Ssh); err != nil {
		return err
	}

	serverImage, err := utils.ComputeImage(flags.Image.Registry, utils.DefaultTag, flags.Image)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	fqdn := args[0]
	if err := utils.IsValidFQDN(fqdn); err != nil {
		return err
	}

	mounts := kubernetes.GetServerMounts()
	mounts = kubernetes.TuneMounts(mounts, &flags.Volumes)

	// Add a mount and volume for the extracted data
	migrationDataVolume := types.VolumeMount{Name: migrationDataPvcName, MountPath: "/var/lib/uyuni-tools"}
	migrationMounts := append(mounts, migrationDataVolume)

	if err := shared_kubernetes.CreatePersistentVolumeClaims(namespace, migrationMounts); err != nil {
		return err
	}

	if err = startMigrationJob(
		namespace,
		serverImage,
		flags.Image.PullPolicy,
		fqdn,
		flags.Migration.User,
		flags.Migration.Prepare,
		migrationMounts,
	); err != nil {
		return err
	}

	// Wait for ever for the job to finish: the duration of this job depends on the amount of data to copy
	if err := shared_kubernetes.WaitForJob(namespace, migrationJobName, -1); err != nil {
		return err
	}

	// Read the extracted data from the migration volume
	extractedData, err := extractMigrationData(namespace, serverImage, migrationDataVolume)
	if err != nil {
		return err
	}

	flags.Installation.TZ = extractedData.Data.Timezone
	flags.Installation.Debug.Java = extractedData.Data.Debug
	if extractedData.Data.HasHubXmlrpcApi {
		flags.HubXmlrpc.Replicas = 1
		flags.HubXmlrpc.IsChanged = true
	}
	flags.Installation.Db.User = extractedData.Data.DbUser
	flags.Installation.Db.Password = extractedData.Data.DbPassword
	// TODO Are those two really needed in migration?
	flags.Installation.Db.Name = extractedData.Data.DbName
	flags.Installation.Db.Port = extractedData.Data.DbPort

	sslDir, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(sslDir)

	// Extract the SSL data as files and pass them as arguments to share code with installation.
	if err := writeToFile(
		extractedData.CaCert, path.Join(sslDir, "ca.crt"), &flags.Installation.Ssl.Ca.Root,
	); err != nil {
		return err
	}

	if err := writeToFile(
		extractedData.CaKey, path.Join(sslDir, "ca.key"), &flags.Installation.Ssl.Ca.Key,
	); err != nil {
		return err
	}

	if err := writeToFile(
		extractedData.ServerCert, path.Join(sslDir, "srv.crt"), &flags.Installation.Ssl.Server.Cert,
	); err != nil {
		return err
	}

	if err := writeToFile(
		extractedData.ServerKey, path.Join(sslDir, "srv.key"), &flags.Installation.Ssl.Server.Key,
	); err != nil {
		return err
	}

	// TODO All the other extracted data should be moved to the inspect step and shared with Upgrade.
	return kubernetes.Reconcile(flags, fqdn)
}

func writeToFile(content string, file string, flag *string) error {
	if content != "" {
		if err := os.WriteFile(file, []byte(content), 0600); err != nil {
			return utils.Errorf(err, L("failed to write certificate to %s"), file)
		}
		*flag = file
	}
	return nil
}

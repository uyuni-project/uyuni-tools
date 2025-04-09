// SPDX-FileCopyrightText: 2025 SUSE LLC
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
	_ *types.GlobalFlags,
	flags *kubernetes.KubernetesServerFlags,
	_ *cobra.Command,
	args []string,
) error {
	namespace := flags.Kubernetes.Uyuni.Namespace

	// Create the namespace if not present
	if err := kubernetes.CreateNamespace(namespace); err != nil {
		return err
	}

	// Check the for the required SSH key and configuration
	if err := checkSSH(namespace, &flags.SSH); err != nil {
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

	// Create a secret using SCC credentials if any are provided
	pullSecret, err := shared_kubernetes.GetRegistrySecret(
		flags.Kubernetes.Uyuni.Namespace, &flags.Installation.SCC, shared_kubernetes.ServerApp,
	)
	if err != nil {
		return err
	}

	jobName, err := startMigrationJob(
		namespace,
		serverImage,
		flags.Image.PullPolicy,
		pullSecret,
		fqdn,
		flags.Migration.User,
		flags.Migration.Prepare,
		migrationMounts,
	)
	if err != nil {
		return err
	}

	// Wait for ever for the job to finish: the duration of this job depends on the amount of data to copy
	if err := shared_kubernetes.WaitForJob(namespace, jobName, -1); err != nil {
		return err
	}

	// Read the extracted data from the migration volume
	extractedData, err := extractMigrationData(
		namespace, serverImage, flags.Image.PullPolicy, pullSecret, migrationDataVolume,
	)
	if err != nil {
		return err
	}

	flags.TZ = extractedData.Data.Timezone
	flags.Installation.Debug.Java = extractedData.Data.Debug
	if extractedData.Data.HasHubXmlrpcAPI {
		flags.HubXmlrpc.Replicas = 1
		flags.HubXmlrpc.IsChanged = true
	}
	flags.Installation.DB.User = extractedData.Data.DBUser
	flags.Installation.DB.Password = extractedData.Data.DBPassword
	// TODO Are those two really needed in migration?
	flags.Installation.DB.Name = extractedData.Data.DBName
	flags.Installation.DB.Port = extractedData.Data.DBPort

	sslDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	// Extract the SSL data as files and pass them as arguments to share code with installation.
	if err := writeToFile(
		extractedData.CaCert, path.Join(sslDir, "ca.crt"), &flags.Installation.SSL.Server.CA.Root,
	); err != nil {
		return err
	}

	// The CA key shouldn't be stored as a temporary file.
	if extractedData.CaKey != "" {
		flags.Installation.SSL.Server.Pair.Key = extractedData.CaKey
	}

	if err := writeToFile(
		extractedData.ServerCert, path.Join(sslDir, "srv.crt"), &flags.Installation.SSL.Server.Pair.Cert,
	); err != nil {
		return err
	}

	if err := writeToFile(
		extractedData.ServerKey, path.Join(sslDir, "srv.key"), &flags.Installation.SSL.Server.Pair.Key,
	); err != nil {
		return err
	}

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

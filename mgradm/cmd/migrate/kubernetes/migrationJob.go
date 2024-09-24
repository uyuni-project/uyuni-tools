// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const migrationJobName = "uyuni-data-sync"

// Prepares and starts the synchronization job.
//
// This assumes the SSH key is stored in an uyuni-migration-key secret
// and the SSH config in an uyuni-migration-ssh ConfigMap with config and known_hosts keys.
func startMigrationJob(
	namespace string,
	serverImage string,
	pullPolicy string,
	pullSecret string,
	fqdn string,
	user string,
	prepare bool,
	mounts []types.VolumeMount,
) (string, error) {
	job, err := getMigrationJob(
		namespace,
		serverImage,
		pullPolicy,
		pullSecret,
		mounts,
		fqdn,
		user,
		prepare,
	)
	if err != nil {
		return "", err
	}

	// Run the job
	return job.ObjectMeta.Name, kubernetes.Apply([]runtime.Object{job}, L("failed to run the migration job"))
}

func getMigrationJob(
	namespace string,
	image string,
	pullPolicy string,
	pullSecret string,
	mounts []types.VolumeMount,
	sourceFqdn string,
	user string,
	prepare bool,
) (*batch.Job, error) {
	// Add mount and volume for the uyuni-migration-key secret with key and key.pub items
	keyMount := core.VolumeMount{Name: "ssh-key", MountPath: "/root/.ssh/id_rsa", SubPath: "id_rsa"}
	pubKeyMount := core.VolumeMount{Name: "ssh-key", MountPath: "/root/.ssh/id_rsa.pub", SubPath: "id_rsa.pub"}

	keyVolume := kubernetes.CreateSecretVolume("ssh-key", "uyuni-migration-key")
	var keyMode int32 = 0600
	keyVolume.VolumeSource.Secret.Items = []core.KeyToPath{
		{Key: "key", Path: "id_rsa", Mode: &keyMode},
		{Key: "key.pub", Path: "id_rsa.pub"},
	}

	// Add mounts and volume for the uyuni-migration-ssh config map
	// We need one mount for each file using subPath to not have 2 mounts on the same folder
	knownHostsMount := core.VolumeMount{Name: "ssh-conf", MountPath: "/root/.ssh/known_hosts", SubPath: "known_hosts"}
	sshConfMount := core.VolumeMount{Name: "ssh-conf", MountPath: "/root/.ssh/config", SubPath: "config"}
	sshVolume := kubernetes.CreateConfigVolume("ssh-conf", "uyuni-migration-ssh")

	// Prepare the script
	scriptData := templates.MigrateScriptTemplateData{
		Volumes:    utils.ServerVolumeMounts,
		SourceFqdn: sourceFqdn,
		User:       user,
		Kubernetes: true,
		Prepare:    prepare,
	}

	job, err := kubernetes.GetScriptJob(namespace, migrationJobName, image, pullPolicy, pullSecret, mounts, scriptData)
	if err != nil {
		return nil, err
	}

	// Append the extra volumes and mounts
	volumeMounts := job.Spec.Template.Spec.Containers[0].VolumeMounts
	volumes := job.Spec.Template.Spec.Volumes

	volumeMounts = append(volumeMounts, keyMount, pubKeyMount, knownHostsMount, sshConfMount)
	volumes = append(volumes, keyVolume, sshVolume)

	job.Spec.Template.Spec.Containers[0].VolumeMounts = volumeMounts
	job.Spec.Template.Spec.Volumes = volumes

	initScript := `cp -a /etc/systemd/system/multi-user.target.wants/. /mnt/etc-systemd-multi`

	job.Spec.Template.Spec.InitContainers = []core.Container{
		{
			Name:            "init-volumes",
			Image:           image,
			ImagePullPolicy: kubernetes.GetPullPolicy(pullPolicy),
			Command:         []string{"sh", "-c", initScript},
			VolumeMounts: []core.VolumeMount{
				{Name: "etc-systemd-multi", MountPath: "/mnt/etc-systemd-multi"},
			},
		},
	}

	return job, nil
}

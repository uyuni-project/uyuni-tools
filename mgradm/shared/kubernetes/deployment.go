// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"os"

	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ServerDeployName is the name of the server deployment.
const ServerDeployName = "uyuni"

// CreateServerDeployment creates a new deployment of the server.
func CreateServerDeployment(
	namespace string,
	serverImage string,
	pullPolicy string,
	timezone string,
	debug bool,
	mirrorPvName string,
) error {
	if mirrorPvName != "" {
		// Create a PVC using the required mirror PV
		if err := kubernetes.CreatePersistentVolumeClaimForVolume(namespace, mirrorPvName); err != nil {
			return err
		}
	}

	serverDeploy := getServerDeployment(
		namespace, serverImage, kubernetes.GetPullPolicy(pullPolicy), timezone, debug, mirrorPvName,
	)

	tempDir, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	return kubernetes.Apply([]runtime.Object{serverDeploy}, L("failed to create the server deployment"))
}

func getServerDeployment(
	namespace string,
	image string,
	pullPolicy core.PullPolicy,
	timezone string,
	debug bool,
	mirrorPvName string,
) *apps.Deployment {
	var replicas int32 = 1

	envs := []core.EnvVar{
		{Name: "TZ", Value: timezone},
	}

	mounts := GetServerMounts()

	if mirrorPvName != "" {
		// Add a mount for the mirror
		mounts = append(mounts, types.VolumeMount{MountPath: "/mirror", Name: mirrorPvName})

		// Add the environment variable for the deployment to use the mirror
		// This doesn't makes sense for migration as the setup script is not executed
		envs = append(envs, core.EnvVar{Name: "MIRROR_PATH", Value: "/mirror"})
	}

	// Convert our mounts to Kubernetes objects
	volumeMounts := kubernetes.ConvertVolumeMounts(mounts)

	// The init mounts are the same mounts but in /mnt just for the init container populating the volumes
	var initMounts []core.VolumeMount
	for _, mount := range volumeMounts {
		initMount := mount.DeepCopy()
		initMount.MountPath = "/mnt" + initMount.MountPath
		initMounts = append(initMounts, *initMount)
	}

	volumes := kubernetes.CreateVolumes(mounts)

	runMount, runVolume := kubernetes.CreateTmpfsMount("/run", "256Mi")
	cgroupMount, cgroupVolume := kubernetes.CreateHostPathMount(
		"/sys/fs/cgroup", "/sys/fs/cgroup", core.HostPathDirectory,
	)

	caMount := core.VolumeMount{
		Name:      "ca-cert",
		MountPath: "/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT",
		ReadOnly:  true,
		SubPath:   "ca.crt",
	}
	tlsKeyMount := core.VolumeMount{Name: "tls-key", MountPath: "/etc/pki/spacewalk-tls"}

	caVolume := kubernetes.CreateConfigVolume("ca-cert", "uyuni-ca")
	tlsKeyVolume := kubernetes.CreateSecretVolume("tls-key", "uyuni-cert")
	var keyMode int32 = 0600
	tlsKeyVolume.VolumeSource.Secret.Items = []core.KeyToPath{
		{Key: "tls.crt", Path: "spacewalk.crt"},
		{Key: "tls.key", Path: "spacewalk.key", Mode: &keyMode},
	}

	initMounts = append(initMounts, tlsKeyMount)
	volumeMounts = append(volumeMounts, runMount, cgroupMount, caMount, tlsKeyMount)
	volumes = append(volumes, runVolume, cgroupVolume, caVolume, tlsKeyVolume)

	// Compute the needed ports
	ports := []types.PortMap{
		utils.NewPortMap("http", 80, 80),
		utils.NewPortMap("https", 443, 443),
	}
	ports = append(ports, utils.TCP_PORTS...)
	ports = append(ports, utils.UDP_PORTS...)
	if debug {
		ports = append(ports, utils.DEBUG_PORTS...)
	}

	deployment := apps.Deployment{
		TypeMeta: meta.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
		ObjectMeta: meta.ObjectMeta{
			Name:      ServerDeployName,
			Namespace: namespace,
			Labels:    map[string]string{"app": kubernetes.ServerApp},
		},
		Spec: apps.DeploymentSpec{
			Replicas: &replicas,
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{"app": kubernetes.ServerApp},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{"app": kubernetes.ServerApp},
				},
				Spec: core.PodSpec{
					InitContainers: []core.Container{
						{
							Name:            "init-volumes",
							Image:           image,
							ImagePullPolicy: pullPolicy,
							Command:         []string{"sh", "-x", "-c", initScript},
							VolumeMounts:    initMounts,
						},
					},
					Containers: []core.Container{
						{
							Name:            "uyuni",
							Image:           image,
							ImagePullPolicy: pullPolicy,
							Lifecycle: &core.Lifecycle{
								PreStop: &core.LifecycleHandler{
									Exec: &core.ExecAction{
										Command: []string{"/bin/sh", "-c", "spacewalk-service stop && systemctl stop postgresql"},
									},
								},
							},
							Ports: kubernetes.ConvertPortMaps(ports),
							Env:   envs,
							ReadinessProbe: &core.Probe{
								ProbeHandler: core.ProbeHandler{
									HTTPGet: &core.HTTPGetAction{
										Port: intstr.FromInt(80),
										Path: "/rhn/manager/login",
									},
								},
								PeriodSeconds:    30,
								TimeoutSeconds:   20,
								FailureThreshold: 5,
							},
							LivenessProbe: &core.Probe{
								ProbeHandler: core.ProbeHandler{
									HTTPGet: &core.HTTPGetAction{
										Port: intstr.FromInt(80),
										Path: "/rhn/manager/login",
									},
								},
								InitialDelaySeconds: 60,
								PeriodSeconds:       60,
								TimeoutSeconds:      20,
								FailureThreshold:    5,
							},
							VolumeMounts: volumeMounts,
						},
					},
					Volumes: volumes,
				},
			},
		},
	}

	return &deployment
}

const initScript = `
# Fill he empty volumes
for vol in /var/lib/cobbler \
		   /var/lib/salt \
		   /var/lib/pgsql \
		   /var/cache \
		   /var/log \
		   /srv/salt \
		   /srv/www \
		   /srv/tftpboot \
		   /srv/formula_metadata \
		   /srv/pillar \
		   /srv/susemanager \
		   /srv/spacewalk \
		   /root \
		   /etc/apache2 \
		   /etc/rhn \
		   /etc/systemd/system/multi-user.target.wants \
		   /etc/systemd/system/sockets.target.wants \
		   /etc/salt \
		   /etc/tomcat \
		   /etc/cobbler \
		   /etc/sysconfig \
		   /etc/postfix \
		   /etc/sssd \
		   /etc/pki/tls
do
	chown --reference=$vol /mnt$vol;
	chmod --reference=$vol /mnt$vol;
	if [ -z "$(ls -A /mnt$vol)" ]; then
    	cp -a $vol/. /mnt$vol;
		if [ "$vol" = "/srv/www" ]; then
            ln -s /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT /mnt$vol/RHN-ORG-TRUSTED-SSL-CERT;
		fi

		if [ "$vol" = "/etc/pki/tls" ]; then
              ln -s /etc/pki/spacewalk-tls/spacewalk.crt /mnt/etc/pki/tls/certs/spacewalk.crt;
              ln -s /etc/pki/spacewalk-tls/spacewalk.key /mnt/etc/pki/tls/private/spacewalk.key;
              cp /etc/pki/spacewalk-tls/spacewalk.key /mnt/etc/pki/tls/private/pg-spacewalk.key;
              chown postgres:postgres /mnt/etc/pki/tls/private/pg-spacewalk.key;
		fi
	fi
done
`

// GetServerMounts returns the volume mounts required for the server pod.
func GetServerMounts() []types.VolumeMount {
	// Filter out the duplicate mounts to avoid issues applying the jobs
	serverMounts := utils.ServerVolumeMounts
	mounts := []types.VolumeMount{}
	mountsSet := map[string]types.VolumeMount{}
	for _, mount := range serverMounts {
		switch mount.Name {
		// Skip mounts that are not PVCs
		case "ca-cert", "tls-key":
			continue
		}
		if _, exists := mountsSet[mount.Name]; !exists {
			mounts = append(mounts, mount)
			mountsSet[mount.Name] = mount
		}
	}

	return mounts
}

// TuneMounts adjusts the server mounts with the size and storage class passed by as parameters.
func TuneMounts(mounts []types.VolumeMount, flags *cmd_utils.VolumesFlags) []types.VolumeMount {
	tunedMounts := []types.VolumeMount{}
	for _, mount := range mounts {
		class := flags.Class
		var volumeFlags *cmd_utils.VolumeFlags
		switch mount.Name {
		case "var-pgsql":
			volumeFlags = &flags.Database
		case "var-spacewalk":
			volumeFlags = &flags.Packages
		case "var-cache":
			volumeFlags = &flags.Cache
		case "srv-www":
			volumeFlags = &flags.Www
		}
		if volumeFlags != nil {
			if volumeFlags.Class != "" {
				class = volumeFlags.Class
			}
			mount.Size = volumeFlags.Size
		}
		mount.Class = class
		tunedMounts = append(tunedMounts, mount)
	}
	return tunedMounts
}

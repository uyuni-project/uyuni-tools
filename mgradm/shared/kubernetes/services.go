// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const webServiceName = "web"
const saltServiceName = "salt"
const cobblerServiceName = "cobbler"
const reportdbServiceName = "report-db"
const taskoServiceName = "taskomatic"
const tftpServiceName = "tftp"

// CreateServices creates the kubernetes services for the server.
//
// If debug is true, the Java debug ports will be exposed.
func CreateServices(namespace string, debug bool) error {
	reportDbPorts := []types.PortMap{
		utils.NewPortMap("pgsql", 5432, 5432),
		utils.NewPortMap("exporter", 9187, 9187),
	}
	saltPorts := []types.PortMap{
		utils.NewPortMap("publish", 4505, 4505),
		utils.NewPortMap("request", 4506, 4506),
		// TODO Add the salt API if configured
	}

	taskoPorts := []types.PortMap{
		utils.NewPortMap("jmx", 5556, 5556),
		utils.NewPortMap("metrics", 9800, 9800),
	}
	tomcatPorts := []types.PortMap{
		utils.NewPortMap("jmx", 5557, 5557),
	}

	if debug {
		taskoPorts = append(taskoPorts, utils.NewPortMap("debug", 8001, 8001))
		tomcatPorts = append(tomcatPorts, utils.NewPortMap("debug", 8003, 8003))
	}

	services := []runtime.Object{
		getService(namespace, webServiceName, core.ProtocolTCP, utils.NewPortMap("web", 80, 80)),
		getService(namespace, saltServiceName, core.ProtocolTCP, saltPorts...),
		getService(namespace, cobblerServiceName, core.ProtocolTCP, utils.NewPortMap("cobbler", 25151, 25151)),
		getService(namespace, reportdbServiceName, core.ProtocolTCP, reportDbPorts...),
		getService(namespace, tftpServiceName, core.ProtocolUDP, utils.NewPortMap("tftp", 69, 69)),
		getService(namespace, "tomcat", core.ProtocolTCP, tomcatPorts...),
		getService(namespace, taskoServiceName, core.ProtocolTCP, taskoPorts...),
	}

	if debug {
		services = append(services,
			getService(namespace, "search", core.ProtocolTCP, utils.NewPortMap("debug", 8002, 8002)),
		)
	}

	return kubernetes.Apply(services, L("failed to create the service"))
}

func getService(namespace string, name string, protocol core.Protocol, ports ...types.PortMap) *core.Service {
	serviceType := core.ServiceTypeClusterIP // TODO make configurable to allow NodePort and maybe LoadBalancer

	portObjs := []core.ServicePort{}
	for _, port := range ports {
		portObjs = append(portObjs, core.ServicePort{
			Name:       port.Name,
			Port:       int32(port.Exposed),
			TargetPort: intstr.FromInt(port.Port),
			Protocol:   protocol,
		})
	}

	return &core.Service{
		TypeMeta: meta.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: meta.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels:    map[string]string{"app": kubernetes.ServerApp},
		},
		Spec: core.ServiceSpec{
			Ports:    portObjs,
			Selector: map[string]string{"app": kubernetes.ServerApp},
			Type:     serviceType,
		},
	}
}

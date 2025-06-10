// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Map the service names to the component label that are not the server.
var serviceMap = map[string]string{
	utils.DBServiceName:       kubernetes.DBComponent,
	utils.ReportdbServiceName: kubernetes.DBComponent,
}

// CreateServices creates the kubernetes services for the server.
//
// If debug is true, the Java debug ports will be exposed.
func CreateServices(namespace string, debug bool) error {
	services := GetServices(namespace, debug)
	for _, svc := range services {
		if !hasCustomService(namespace, svc.Name) {
			if err := kubernetes.Apply([]*core.Service{svc}, L("failed to create the service")); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetServices creates the definitions of all the services of the server.
//
// If debug is true, the Java debug ports will be exposed.
func GetServices(namespace string, debug bool) []*core.Service {
	ports := utils.GetServerPorts(debug)
	ports = append(ports, utils.DBPorts...)
	ports = append(ports, utils.ReportDBPorts...)

	servicesPorts := map[string][]types.PortMap{}
	for _, port := range ports {
		svcPorts := servicesPorts[port.Service]
		if svcPorts == nil {
			svcPorts = []types.PortMap{}
		}
		svcPorts = append(svcPorts, port)
		servicesPorts[port.Service] = svcPorts
	}

	services := []*core.Service{}
	for _, svcPorts := range servicesPorts {
		protocol := core.ProtocolTCP
		if svcPorts[0].Protocol == "udp" {
			protocol = core.ProtocolUDP
		}
		// Do we have a split component for that service already?
		component := kubernetes.ServerComponent
		if comp, exists := serviceMap[svcPorts[0].Service]; exists {
			component = comp
		}
		services = append(services,
			getService(namespace, kubernetes.ServerApp, component, svcPorts[0].Service, protocol, svcPorts...),
		)
	}
	return services
}

func getService(
	namespace string,
	app string,
	component string,
	name string,
	protocol core.Protocol,
	ports ...types.PortMap,
) *core.Service {
	// TODO make configurable to allow NodePort and maybe LoadBalancer for exposed services.
	serviceType := core.ServiceTypeClusterIP

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
			Labels:    kubernetes.GetLabels(app, component),
		},
		Spec: core.ServiceSpec{
			Ports:    portObjs,
			Selector: map[string]string{kubernetes.ComponentLabel: component},
			Type:     serviceType,
		},
	}
}

func hasCustomService(namespace string, name string) bool {
	out, err := utils.RunCmdOutput(
		zerolog.DebugLevel, "kubectl", "get", "svc", "-n", namespace, name,
		"-l", fmt.Sprintf("%s!=%s", kubernetes.AppLabel, kubernetes.ServerApp),
		"-o", "jsonpath={.items[?(@.metadata.name=='db')]}",
	)
	// Custom services don't have our app label!
	return err == nil && strings.TrimSpace(string(out)) != ""
}

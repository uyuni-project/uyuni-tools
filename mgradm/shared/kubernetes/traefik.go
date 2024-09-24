// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"

	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// CreateTraefikRoutes creates the routes and middleware wiring the traefik endpoints to their service.
func CreateTraefikRoutes(namespace string, hub bool, debug bool) error {
	routeTemplate := template.Must(template.New("ingressRoute").Parse(ingressRouteTemplate))

	tempDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	filePath := path.Join(tempDir, "routes.yaml")
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return utils.Errorf(err, L("failed to open %s for writing"), filePath)
	}
	defer file.Close()

	// Write the SSL Redirect middleware
	_, err = file.WriteString(fmt.Sprintf(`
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: uyuni-https-redirect
  namespace: "%s"
  labels:
    %s: %s
spec:
  redirectScheme:
    scheme: https
    permanent: true
`, namespace, kubernetes.AppLabel, kubernetes.ServerApp))
	if err != nil {
		return utils.Errorf(err, L("failed to write traefik middleware and routes to file"))
	}

	// Write the routes from the endpoint to the services
	for _, endpoint := range getPortList(hub, debug) {
		_, err := file.WriteString("---\n")
		if err != nil {
			return utils.Errorf(err, L("failed to write traefik middleware and routes to file"))
		}
		if err := getTraefixRoute(routeTemplate, file, namespace, endpoint); err != nil {
			return err
		}
	}
	if err := file.Close(); err != nil {
		return utils.Errorf(err, L("failed to close traefik middleware and routes file"))
	}

	if _, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "apply", "-f", filePath); err != nil {
		return utils.Errorf(err, L("failed to create traefik middleware and routes"))
	}
	return nil
}

func getTraefixRoute(t *template.Template, writer io.Writer, namespace string, endpoint types.PortMap) error {
	endpointName := kubernetes.GetTraefikEndpointName(endpoint)
	protocol := "TCP"
	if endpoint.Protocol == "udp" {
		protocol = "UDP"
	}

	data := routeData{
		Name:      endpointName + "-route",
		Namespace: namespace,
		EndPoint:  endpointName,
		Service:   endpoint.Service,
		Port:      endpoint.Exposed,
		Protocol:  protocol,
	}
	if err := t.Execute(writer, data); err != nil {
		return utils.Errorf(err, L("failed to write traefik routes to file"))
	}
	return nil
}

type routeData struct {
	Name      string
	Namespace string
	EndPoint  string
	Service   string
	Port      int
	Protocol  string
}

const ingressRouteTemplate = `
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute{{ .Protocol }}
metadata:
  name: {{ .Name }}
  namespace: "{{ .Namespace }}"
  labels:
    ` + kubernetes.AppLabel + ": " + kubernetes.ServerApp + `
spec:
  entryPoints:
    - {{ .EndPoint }}
  routes:
    - services:
      - name: {{ .Service }}
        port: {{ .Port }}
{{- if eq .Protocol "TCP" }}
      match: ` + "HostSNI(`*`)" + `
{{- end }}
`

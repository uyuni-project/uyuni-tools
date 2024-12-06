// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestGetTraefikRouteTCP(t *testing.T) {
	routeTemplate := template.Must(template.New("ingressRoute").Parse(ingressRouteTemplate))

	var buf bytes.Buffer
	err := getTraefixRoute(routeTemplate, &buf, "foo", utils.NewPortMap("svcname", "port1", 123, 456))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	actual := buf.String()
	expected := `
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRouteTCP
metadata:
  name: svcname-port1-route
  namespace: "foo"
  labels:
    app.kubernetes.io/part-of: uyuni
spec:
  entryPoints:
    - svcname-port1
  routes:
    - services:
      - name: svcname
        port: 123
      match: ` + "HostSNI(`*`)\n"
	testutils.AssertEquals(t, "Wrong traefik route generated", expected, actual)
}

func TestGetTraefikRouteUDP(t *testing.T) {
	routeTemplate := template.Must(template.New("ingressRoute").Parse(ingressRouteTemplate))

	var buf bytes.Buffer
	err := getTraefixRoute(routeTemplate, &buf, "foo",
		types.PortMap{
			Service:  "svcname",
			Name:     "port1",
			Exposed:  123,
			Port:     456,
			Protocol: "udp",
		})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	actual := buf.String()
	expected := `
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRouteUDP
metadata:
  name: svcname-port1-route
  namespace: "foo"
  labels:
    app.kubernetes.io/part-of: uyuni
spec:
  entryPoints:
    - svcname-port1
  routes:
    - services:
      - name: svcname
        port: 123
`
	testutils.AssertEquals(t, "Wrong traefik route generated", expected, actual)
}

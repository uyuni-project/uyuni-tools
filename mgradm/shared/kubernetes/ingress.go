// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	net "k8s.io/api/networking/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CertSecretName is the name of the server SSL certificate secret to use.
const CertSecretName = "uyuni-cert"

// CreateIngress creates the ingress definitions for Uyuni server.
//
// fqdn is the fully qualified domain name associated with the Uyuni server.
//
// caIssuer is the name of the cert-manager to associate for the SSL routes.
// It can be empty if cert-manager is not used.
//
// ingressName is one of traefik or nginx.
func CreateIngress(namespace string, fqdn string, caIssuer string, ingressName string) error {
	ingresses := []runtime.Object{
		getSslIngress(namespace, fqdn, caIssuer, ingressName),
		getNoSslIngress(namespace, fqdn, ingressName),
	}
	sslRedirectIngress := getSslRedirectIngress(namespace, fqdn, ingressName)
	if sslRedirectIngress != nil {
		ingresses = append(ingresses, sslRedirectIngress)
	}

	return kubernetes.Apply(ingresses, L("failed to create the hub API service"))
}

func getSslIngress(namespace string, fqdn string, caIssuer string, ingressName string) *net.Ingress {
	annotations := map[string]string{}
	if caIssuer != "" {
		annotations["cert-manager.io/issuer"] = caIssuer
	}
	if ingressName == "traefik" {
		annotations["traefik.ingress.kubernetes.io/router.tls"] = "true"
		annotations["traefik.ingress.kubernetes.io/router.tls.domains.n.main"] = fqdn
		annotations["traefik.ingress.kubernetes.io/router.entrypoints"] = "websecure,web"
	}

	ingress := net.Ingress{
		TypeMeta: meta.TypeMeta{APIVersion: "networking.k8s.io/v1", Kind: "Ingress"},
		ObjectMeta: meta.ObjectMeta{
			Namespace:   namespace,
			Name:        "uyuni-ingress-ssl",
			Annotations: annotations,
			Labels:      map[string]string{"app": kubernetes.ServerApp},
		},
		Spec: net.IngressSpec{
			TLS: []net.IngressTLS{
				{Hosts: []string{fqdn}, SecretName: CertSecretName},
			},
			Rules: []net.IngressRule{
				getIngressWebRule(fqdn),
			},
		},
	}

	return &ingress
}

func getSslRedirectIngress(namespace string, fqdn string, ingressName string) *net.Ingress {
	var ingress *net.Ingress

	// Nginx doesn't require a special ingress for the SSL redirection.
	if ingressName == "traefik" {
		ingress = &net.Ingress{
			TypeMeta: meta.TypeMeta{APIVersion: "networking.k8s.io/v1", Kind: "Ingress"},
			ObjectMeta: meta.ObjectMeta{
				Namespace: namespace,
				Name:      "uyuni-ingress-ssl-redirect",
				Annotations: map[string]string{
					"traefik.ingress.kubernetes.io/router.middlewares": "default-uyuni-https-redirect@kubernetescrd",
					"traefik.ingress.kubernetes.io/router.entrypoints": "web",
				},
				Labels: map[string]string{"app": kubernetes.ServerApp},
			},
			Spec: net.IngressSpec{
				Rules: []net.IngressRule{
					getIngressWebRule(fqdn),
				},
			},
		}
	}

	return ingress
}

var noSslPaths = []string{
	"/pub",
	"/rhn/([^/])+/DownloadFile",
	"/(rhn/)?rpc/api",
	"/rhn/errors",
	"/rhn/ty/TinyUrl",
	"/rhn/websocket",
	"/rhn/metrics",
	"/cobbler_api",
	"/cblr",
	"/httpboot",
	"/images",
	"/cobbler",
	"/os-images",
	"/tftp",
	"/docs",
}

func getNoSslIngress(namespace string, fqdn string, ingressName string) *net.Ingress {
	annotations := map[string]string{}
	if ingressName == "nginx" {
		annotations["nginx.ingress.kubernetes.io/ssl-redirect"] = "false"
	}
	if ingressName == "traefik" {
		annotations["traefik.ingress.kubernetes.io/router.tls"] = "false"
		annotations["traefik.ingress.kubernetes.io/router.entrypoints"] = "web"
	}

	pathType := net.PathTypePrefix
	paths := []net.HTTPIngressPath{}
	for _, noSslPath := range noSslPaths {
		paths = append(paths, net.HTTPIngressPath{
			Backend:  webServiceBackend,
			Path:     noSslPath,
			PathType: &pathType,
		})
	}

	ingress := net.Ingress{
		TypeMeta: meta.TypeMeta{APIVersion: "networking.k8s.io/v1", Kind: "Ingress"},
		ObjectMeta: meta.ObjectMeta{
			Namespace:   namespace,
			Name:        "uyuni-ingress-nossl",
			Annotations: annotations,
			Labels:      map[string]string{"app": kubernetes.ServerApp},
		},
		Spec: net.IngressSpec{
			TLS: []net.IngressTLS{
				{Hosts: []string{fqdn}, SecretName: CertSecretName},
			},
			Rules: []net.IngressRule{
				{
					Host: fqdn,
					IngressRuleValue: net.IngressRuleValue{
						HTTP: &net.HTTPIngressRuleValue{Paths: paths},
					},
				},
			},
		},
	}

	return &ingress
}

// build the ingress rule object catching all HTTP traffic.
func getIngressWebRule(fqdn string) net.IngressRule {
	pathType := net.PathTypePrefix

	return net.IngressRule{
		Host: fqdn,
		IngressRuleValue: net.IngressRuleValue{
			HTTP: &net.HTTPIngressRuleValue{
				Paths: []net.HTTPIngressPath{
					{
						Backend:  webServiceBackend,
						Path:     "/",
						PathType: &pathType,
					},
				},
			},
		},
	}
}

var webServiceBackend net.IngressBackend = net.IngressBackend{
	Service: &net.IngressServiceBackend{
		Name: webServiceName,
		Port: net.ServiceBackendPort{Number: 80},
	},
}

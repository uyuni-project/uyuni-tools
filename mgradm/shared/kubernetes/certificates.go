// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"encoding/base64"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func installTLSSecret(namespace string, serverCrt []byte, serverKey []byte, rootCaCrt []byte) error {
	crdsDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	secretPath := filepath.Join(crdsDir, "secret.yaml")
	log.Info().Msg(L("Creating SSL server certificate secret"))
	tlsSecretData := templates.TLSSecretTemplateData{
		Namespace:   namespace,
		Name:        "uyuni-cert",
		Certificate: base64.StdEncoding.EncodeToString(serverCrt),
		Key:         base64.StdEncoding.EncodeToString(serverKey),
		RootCa:      base64.StdEncoding.EncodeToString(rootCaCrt),
	}

	if err = utils.WriteTemplateToFile(tlsSecretData, secretPath, 0500, true); err != nil {
		return utils.Errorf(err, L("Failed to generate uyuni-crt secret definition"))
	}
	err = utils.RunCmd("kubectl", "apply", "-f", secretPath)
	if err != nil {
		return utils.Errorf(err, L("Failed to create uyuni-crt TLS secret"))
	}

	createCaConfig(namespace, rootCaCrt)
	return nil
}

// Install cert-manager and its CRDs using helm in the cert-manager namespace if needed
// and then create a self-signed CA and issuers.
// Returns helm arguments to be added to use the issuer.
func installSSLIssuers(helmFlags *cmd_utils.HelmFlags, sslFlags *cmd_utils.InstallSSLFlags, rootCa string,
	tlsCert *types.SSLPair, kubeconfig, fqdn string, imagePullPolicy string) ([]string, error) {
	// Install cert-manager if needed
	if err := installCertManager(helmFlags, kubeconfig, imagePullPolicy); err != nil {
		return []string{}, utils.Errorf(err, L("cannot install cert manager"))
	}

	log.Info().Msg(L("Creating SSL certificate issuer"))
	crdsDir, cleaner, err := utils.TempDir()
	if err != nil {
		return []string{}, err
	}
	defer cleaner()

	issuerPath := filepath.Join(crdsDir, "issuer.yaml")

	issuerData := templates.IssuerTemplateData{
		Namespace:   helmFlags.Uyuni.Namespace,
		Country:     sslFlags.Country,
		State:       sslFlags.State,
		City:        sslFlags.City,
		Org:         sslFlags.Org,
		OrgUnit:     sslFlags.OU,
		Email:       sslFlags.Email,
		Fqdn:        fqdn,
		RootCa:      rootCa,
		Key:         tlsCert.Key,
		Certificate: tlsCert.Cert,
	}

	if err = utils.WriteTemplateToFile(issuerData, issuerPath, 0500, true); err != nil {
		return []string{}, utils.Errorf(err, L("failed to generate issuer definition"))
	}

	err = utils.RunCmd("kubectl", "apply", "-f", issuerPath)
	if err != nil {
		log.Fatal().Err(err).Msg(L("Failed to create issuer"))
	}

	// Wait for issuer to be ready
	for i := 0; i < 60; i++ {
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "-o=jsonpath={.status.conditions[*].type}",
			"issuer", "uyuni-ca-issuer", "-n", issuerData.Namespace)
		if err == nil && string(out) == "Ready" {
			return []string{"--set-json", "ingressSSLAnnotations={\"cert-manager.io/issuer\": \"uyuni-ca-issuer\"}"}, nil
		}
		time.Sleep(1 * time.Second)
	}
	log.Fatal().Msg(L("Issuer didn't turn ready after 60s"))
	return []string{}, nil
}

func installCertManager(helmFlags *cmd_utils.HelmFlags, kubeconfig string, imagePullPolicy string) error {
	if !kubernetes.IsDeploymentReady("", "cert-manager") {
		log.Info().Msg(L("Installing cert-manager"))
		repo := ""
		chart := helmFlags.CertManager.Chart
		version := helmFlags.CertManager.Version
		namespace := helmFlags.CertManager.Namespace

		args := []string{
			"--set", "crds.enabled=true",
			"--set-json", "global.commonLabels={\"installedby\": \"mgradm\"}",
			"--set", "image.pullPolicy=" + kubernetes.GetPullPolicy(imagePullPolicy),
		}
		extraValues := helmFlags.CertManager.Values
		if extraValues != "" {
			args = append(args, "-f", extraValues)
		}

		// Use upstream chart if nothing defined
		if chart == "" {
			repo = "https://charts.jetstack.io"
			chart = "cert-manager"
		}
		// The installedby label will be used to only uninstall what we installed
		if err := kubernetes.HelmUpgrade(
			kubeconfig, namespace, true, repo, "cert-manager", chart, version, args...,
		); err != nil {
			return utils.Errorf(err, L("cannot run helm upgrade"))
		}
	}

	// Wait for cert-manager to be ready
	err := kubernetes.WaitForDeployment("", "cert-manager-webhook", "webhook")
	if err != nil {
		return utils.Errorf(err, L("cannot deploy"))
	}

	return nil
}

func extractCaCertToConfig(namespace string) {
	// TODO Replace with [trust-manager](https://cert-manager.io/docs/projects/trust-manager/) to automate this
	const jsonPath = "-o=jsonpath={.data.ca\\.crt}"

	log.Info().Msg(L("Extracting CA certificate to a configmap"))
	// Skip extracting if the configmap is already present
	out, err := utils.RunCmdOutput(
		zerolog.DebugLevel, "kubectl", "get", "configmap", "uyuni-ca", jsonPath, "-n", namespace,
	)
	log.Info().Msgf(L("CA cert: %s"), string(out))
	if err == nil && len(out) > 0 {
		log.Info().Msg(L("uyuni-ca configmap already existing, skipping extraction"))
		return
	}

	out, err = utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "secret", "uyuni-ca", jsonPath, "-n", namespace)
	if err != nil {
		log.Fatal().Err(err).Msgf(L("Failed to get uyuni-ca certificate"))
	}

	decoded, err := base64.StdEncoding.DecodeString(string(out))
	if err != nil {
		log.Fatal().Err(err).Msgf(L("Failed to base64 decode CA certificate"))
	}

	createCaConfig(namespace, decoded)
}

func createCaConfig(namespace string, ca []byte) {
	valueArg := "--from-literal=ca.crt=" + string(ca)
	if err := utils.RunCmd("kubectl", "create", "configmap", "uyuni-ca", valueArg, "-n", namespace); err != nil {
		log.Fatal().Err(err).Msg(L("Failed to create uyuni-ca config map from certificate"))
	}
}

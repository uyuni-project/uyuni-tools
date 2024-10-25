// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

// Mappings for the models/Schemas in scope of the proxy API

// ProxyConfigRequestToMap maps the ProxyConfigRequest to a map.
func ProxyConfigRequestToMap(request ProxyConfigRequest) map[string]interface{} {
	return map[string]interface{}{
		"proxyName":       request.ProxyName,
		"proxyPort":       request.ProxyPort,
		"server":          request.Server,
		"maxCache":        request.MaxCache,
		"email":           request.Email,
		"rootCA":          request.RootCA,
		"proxyCrt":        request.ProxyCrt,
		"proxyKey":        request.ProxyKey,
		"intermediateCAs": request.IntermediateCAs,
	}
}

// ProxyConfigGenerateRequestToMap maps the ProxyConfigGenerateRequest to a map.
func ProxyConfigGenerateRequestToMap(request ProxyConfigGenerateRequest) map[string]interface{} {
	return map[string]interface{}{
		"proxyName":  request.ProxyName,
		"proxyPort":  request.ProxyPort,
		"server":     request.Server,
		"maxCache":   request.MaxCache,
		"email":      request.Email,
		"caCrt":      request.CaCrt,
		"caKey":      request.CaKey,
		"caPassword": request.CaPassword,
		"cnames":     request.Cnames,
		"country":    request.Country,
		"state":      request.State,
		"city":       request.City,
		"org":        request.Org,
		"orgUnit":    request.OrgUnit,
		"sslEmail":   request.SslEmail,
	}
}

// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package utils

// KubernetesBuilt is a flag for compiling kubernetes code. True when go:build !nok8s, False when go:build nok8s.
const KubernetesBuilt = true

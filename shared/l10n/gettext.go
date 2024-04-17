// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package l10n

import "github.com/chai2010/gettext-go"

// L localizes a string using the set up gettext domain and locale.
// This is an alias for gettext.Gettext().
func L(message string) string {
	return gettext.Gettext(message)
}

// NL returns a localized message depending on the value of count.
// This is an alias for gettext.NGettext().
func NL(message string, plural string, count int) string {
	return gettext.NGettext(message, plural, count)
}

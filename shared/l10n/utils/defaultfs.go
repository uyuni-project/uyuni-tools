// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package l10n

import "github.com/chai2010/gettext-go"

// DefaultFS providing a empty data if no data is found.
type DefaultFS struct {
	osFs gettext.FileSystem
	gettext.FileSystem
}

// New creates a new DefaultFS delegating to an OS FileSystem.
func New(path string) *DefaultFS {
	return &DefaultFS{
		osFs: gettext.OS(path),
	}
}

// LocaleList gets the list of locales from the underlying os FileSystem.
func (f *DefaultFS) LocaleList() []string {
	return f.osFs.LocaleList()
}

// LoadMessagesFile loads a messages or returns the content of an empty json file.
func (f *DefaultFS) LoadMessagesFile(domain, lang, ext string) ([]byte, error) {
	osFile, err := f.osFs.LoadMessagesFile(domain, lang, ext)
	// Return an empty file by default
	if err != nil {
		return []byte("[]"), nil
	}
	return osFile, nil
}

// LoadResourceFile loads the resource file or returns empty data.
func (f *DefaultFS) LoadResourceFile(domain, lang, ext string) ([]byte, error) {
	osFile, err := f.osFs.LoadResourceFile(domain, lang, ext)
	// Return an empty file by default
	if err != nil {
		return []byte{}, nil
	}
	return osFile, nil
}

// String returns a name of the FileSystem.
func (f *DefaultFS) String() string {
	return "DefaultFS"
}

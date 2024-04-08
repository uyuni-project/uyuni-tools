// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// Extracts a tar.gz file.
func ExtractTarGz(tarballPath string, dstPath string) error {
	reader, err := os.Open(tarballPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	tarReader := tar.NewReader(archive)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path, err := filepath.Abs(filepath.Join(dstPath, header.Name))
		if err != nil {
			return err
		}
		if !strings.HasPrefix(path, dstPath) {
			log.Warn().Msgf(L("Skipping extraction of %s in %s file as it resolves outside the target path"),
				header.Name, tarballPath)
			continue
		}

		info := header.FileInfo()
		if info.IsDir() {
			log.Debug().Msgf("Creating folder %s", path)
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		log.Debug().Msgf("Extracting file %s", path)
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err = io.Copy(file, tarReader); err != nil {
			return err
		}
	}

	return nil
}

// Object holding a .tar.gz to write it to a file.
type TarGz struct {
	fileWriter *os.File
	tarWriter  *tar.Writer
	gzipWriter *gzip.Writer
}

// NewTarGz create a targz object with writers opened.
// A successful call should be followed with a close.
func NewTarGz(path string) (*TarGz, error) {
	var targz TarGz
	var err error
	targz.fileWriter, err = os.Create(path)
	if err != nil {
		return nil, fmt.Errorf(L("failed to write tar.gz to %s: %s"), path, err)
	}

	targz.gzipWriter = gzip.NewWriter(targz.fileWriter)
	targz.tarWriter = tar.NewWriter(targz.gzipWriter)
	return &targz, nil
}

// Close stops all the writers.
func (t *TarGz) Close() {
	t.tarWriter.Close()
	t.gzipWriter.Close()
	t.fileWriter.Close()
}

// AddFile adds the file at filepath to the archive as entrypath.
func (t *TarGz) AddFile(filepath string, entrypath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	header.Name = entrypath
	if err = t.tarWriter.WriteHeader(header); err != nil {
		return err
	}

	if _, err = io.Copy(t.tarWriter, file); err != nil {
		return err
	}
	return nil
}

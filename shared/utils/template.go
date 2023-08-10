package utils

import (
	"io"
	"log"
	"os"
)

type Template interface {
	Render(wr io.Writer) error
}

func WriteTemplateToFile(template Template, path string, perm os.FileMode, overwrite bool) error {
	// Check if the file is existing
	if !overwrite {
		if FileExists(path) {
			log.Fatalf("%s file already present, not overwriting\n", path)
		}
	}

	// Write the configuration
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		log.Fatalf("Failed to open %s for writing: %s\n", path, err)
	}
	defer file.Close()

	return template.Render(file)
}

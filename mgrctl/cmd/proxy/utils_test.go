package proxy_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgrctl/cmd/proxy"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
)

// Test promptForPassword function.
func TestPromptForPassword(t *testing.T) {
	input := "test_password\n"
	expectedOutput := "test_password"

	// Redirect os.Stdin to read from the simulated input
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Write the simulated input to the pipe
	_, err = w.Write([]byte(input))
	if err != nil {
		t.Fatalf("Failed to write to pipe: %v", err)
	}
	w.Close()

	// Save the original os.Stdin and defer restoring it
	oStdin := os.Stdin
	defer func() {
		os.Stdin = oStdin
	}()

	// Set os.Stdin to read from the pipe
	os.Stdin = r

	// Capture the output of the promptForPassword function
	output := proxy.PromptForPassword()

	// Trim the newline character from the output
	output = strings.TrimSpace(output)

	// Compare the output with the expected output
	if output != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, output)
	}
}

// Test getFilename function.
func TestGetFilename(t *testing.T) {
	// Test when output is empty
	filename := proxy.GetFilename("", "testProxy.domain.com")
	test_utils.AssertEquals(t, "", "testProxy-config.tar.gz", filename)

	// Test when output is provided
	filename = proxy.GetFilename("customOutput", "testProxy.domain.com")
	test_utils.AssertEquals(t, "", "customOutput.tar.gz", filename)

	// Test when output is provided
	filename = proxy.GetFilename("/var/customOutputWitPath", "testProxy.domain.com")
	test_utils.AssertEquals(t, "", "/var/customOutputWitPath.tar.gz", filename)
}

func createTestFile(dir string, filename string, content string, t *testing.T) string {
	filepath := path.Join(dir, filename)
	test_utils.WriteFile(t, filepath, content)
	return filepath
}

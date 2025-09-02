// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/onsi/ginkgo/v2"
)

var MgradmPath string
var MgrctlPath string
var MgrpxyPath string
var MgradmConfigPath string

// interaction defines a prompt to look for and the corresponding response to send.
type interaction struct {
	Prompt   string
	Response string
}

// RunCommand executes a command, handling an ordered sequence of interactive prompts.
func RunCommand(command string, args []string, interactions ...interaction) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	var combinedOutputBuffer bytes.Buffer
	var wg sync.WaitGroup

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("error creating stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("error creating stderr pipe: %w", err)
	}

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("error creating stdin pipe: %w", err)
	}

	interactionReader, interactionWriter := io.Pipe()
	mainWriter := io.MultiWriter(&combinedOutputBuffer, ginkgo.GinkgoWriter, interactionWriter)

	ginkgo.By(fmt.Sprintf("Running command: %s %v", command, args))
	if err := cmd.Start(); err != nil {
		err := interactionWriter.Close()
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("error starting command: %w", err)
	}

	// Start a goroutine to read from the interactionReader and handle interactions.
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer stdinPipe.Close()

		interactionIndex := 0
		scanner := bufio.NewScanner(interactionReader)
		for scanner.Scan() {
			handleNextInteraction(scanner.Text(), stdinPipe, interactions, &interactionIndex)
		}
	}()

	// Read from stdout and stderr concurrently, writing to the combined output buffer.
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := io.Copy(mainWriter, stdoutPipe)
		if err != nil {
			return
		}
	}()
	go func() {
		defer wg.Done()
		_, err := io.Copy(mainWriter, stderrPipe)
		if err != nil {
			return
		}
	}()

	err = cmd.Wait()
	interactionWriter.Close()
	wg.Wait()

	return combinedOutputBuffer.String(), err
}

// handleNextInteraction checks a line of output against the next expected interaction.
// If it matches, it sends the corresponding response and advances the interaction index.
func handleNextInteraction(line string, stdinPipe io.WriteCloser, interactions []interaction, interactionIndex *int) {
	// Do nothing if all interactions have been handled.
	if *interactionIndex >= len(interactions) {
		return
	}

	currentInteraction := interactions[*interactionIndex]
	if strings.Contains(line, currentInteraction.Prompt) {
		ginkgo.By(fmt.Sprintf("Answering prompt '%s' with '%s'", currentInteraction.Prompt, currentInteraction.Response))
		_, err := io.WriteString(stdinPipe, currentInteraction.Response+"\n")
		if err != nil {
			// Log the error without stopping the entire process.
			ginkgo.GinkgoWriter.Printf("Error writing to stdin: %v", err)
		}

		*interactionIndex++
	}
}

func RunMgradmCommand(args []string, interactions ...interaction) (string, error) {
	return RunCommand(MgradmPath, args, interactions...)
}

func RunMgrctlCommand(args []string, interactions ...interaction) (string, error) {
	return RunCommand(MgrctlPath, args, interactions...)
}

func RunMgrpxyCommand(args []string, interactions ...interaction) (string, error) {
	return RunCommand(MgrpxyPath, args, interactions...)
}

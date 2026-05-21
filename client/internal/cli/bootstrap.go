package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cc-status/client/internal/claude"
	"cc-status/client/internal/lock"
	"cc-status/client/internal/runtime"
)

type BootstrapOptions struct {
	AppDir            string
	ClaudeProjectsDir string
	EnvLookup         func(string) string
	Stdout            io.Writer
}

func NewBootstrapRunner(options BootstrapOptions) Runner {
	return func(command string) error {
		appDir, err := resolveAppDir(options.AppDir)
		if err != nil {
			return err
		}

		heldLock, err := lock.Acquire(filepath.Join(appDir, "client.lock"))
		if err != nil {
			return err
		}
		defer heldLock.Release()

		state, err := runtime.Bootstrap(runtime.Options{
			AppDir:    appDir,
			EnvLookup: options.EnvLookup,
		})
		if err != nil {
			return err
		}
		defer state.Close()

		if command != "dry-run" {
			return nil
		}

		projectsDir, err := resolveClaudeProjectsDir(options.ClaudeProjectsDir)
		if err != nil {
			return err
		}

		result, err := claude.ScanProjectsDir(projectsDir)
		if err != nil {
			return err
		}

		stdout := options.Stdout
		if stdout == nil {
			stdout = os.Stdout
		}

		_, err = fmt.Fprintf(
			stdout,
			"dry-run summary: files_scanned=%d records=%d errors=%d\n",
			result.FilesScanned,
			len(result.Records),
			len(result.Errors),
		)
		if err != nil {
			return err
		}

		for _, fileError := range result.Errors {
			if _, err := fmt.Fprintf(stdout, "dry-run error: %s\n", fileError); err != nil {
				return err
			}
		}

		return nil
	}
}

func resolveAppDir(appDir string) (string, error) {
	if appDir != "" {
		return appDir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}

	return filepath.Join(homeDir, ".cc-usage-client"), nil
}

func resolveClaudeProjectsDir(projectsDir string) (string, error) {
	if projectsDir != "" {
		return projectsDir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}

	return filepath.Join(homeDir, ".claude", "projects"), nil
}

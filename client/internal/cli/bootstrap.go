package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cc-status/client/internal/claude"
	"cc-status/client/internal/httpclient"
	"cc-status/client/internal/lock"
	"cc-status/client/internal/runtime"
	"cc-status/client/internal/syncer"
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

		stdout := stdoutOrDefault(options.Stdout)
		_, _ = fmt.Fprintf(stdout, "[初始化] 客户端: %s\n", state.ClientID)

		if command != "dry-run" {
			return runSync(state, options, stdout)
		}

		projectsDir, err := resolveClaudeProjectsDir(options.ClaudeProjectsDir)
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintln(stdout, "[扫描] 正在扫描 Claude 使用记录...")
		result, err := claude.ScanProjectsDir(projectsDir)
		if err != nil {
			return err
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

func runSync(state *runtime.State, options BootstrapOptions, stdout io.Writer) error {
	projectsDir, err := resolveClaudeProjectsDir(options.ClaudeProjectsDir)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintln(stdout, "[扫描] 正在扫描 Claude 使用记录...")
	fileResults, err := claude.ScanProjectFiles(projectsDir)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(stdout, "[扫描] 完成, 共扫描 %d 个文件\n", len(fileResults))

	syncClient := httpclient.NewSyncClient(
		state.Config.ServerURL,
		state.Config.AuthToken,
		state.Config.TimeoutSeconds,
	)
	result, err := syncer.RunHappyPath(
		state.Store,
		state.ClientID,
		syncClient,
		fileResults,
		state.Config.BatchSize,
		stdout,
	)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(
		stdout,
		"sync summary: files_scanned=%d records=%d accepted=%d skipped=%d errors=%d\n",
		result.FilesScanned,
		result.Records,
		result.Accepted,
		result.Skipped,
		len(result.Errors),
	)
	return err
}

func stdoutOrDefault(stdout io.Writer) io.Writer {
	if stdout != nil {
		return stdout
	}
	return os.Stdout
}
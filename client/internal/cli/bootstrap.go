package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"cc-status/client/internal/lock"
	"cc-status/client/internal/runtime"
)

type BootstrapOptions struct {
	AppDir    string
	EnvLookup func(string) string
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

		// issue 01 只要求打通启动链路，业务同步逻辑后续 issue 再补齐。
		_ = command
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

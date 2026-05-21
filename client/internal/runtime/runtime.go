package runtime

import (
	"fmt"
	"os"
	"path/filepath"

	"cc-status/client/internal/config"
	"cc-status/client/internal/storage"
)

type Options struct {
	AppDir    string
	EnvLookup func(string) string
}

type State struct {
	AppDir   string
	ClientID string
	Config   config.Config
	Store    *storage.Store
}

func Bootstrap(options Options) (*State, error) {
	appDir := options.AppDir
	if appDir == "" {
		return nil, fmt.Errorf("app dir is required")
	}

	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return nil, fmt.Errorf("create app dir: %w", err)
	}

	cfg, err := config.Load(appDir, options.EnvLookup)
	if err != nil {
		return nil, err
	}

	store, err := storage.Open(filepath.Join(appDir, "client.db"))
	if err != nil {
		return nil, err
	}

	clientID, err := store.LoadOrCreateClientID()
	if err != nil {
		_ = store.Close()
		return nil, err
	}

	return &State{
		AppDir:   appDir,
		ClientID: clientID,
		Config:   cfg,
		Store:    store,
	}, nil
}

func (state *State) Close() error {
	if state == nil {
		return nil
	}
	return state.Store.Close()
}

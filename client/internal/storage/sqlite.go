package storage

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	store := &Store{db: db}
	if err := store.ensureSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (store *Store) Close() error {
	if store == nil || store.db == nil {
		return nil
	}
	return store.db.Close()
}

func (store *Store) LoadOrCreateClientID() (string, error) {
	const key = "client_id"

	var clientID string
	err := store.db.QueryRow(`SELECT value FROM metadata WHERE key = ?`, key).Scan(&clientID)
	if err == nil {
		return clientID, nil
	}

	if err != sql.ErrNoRows {
		return "", fmt.Errorf("query client id: %w", err)
	}

	clientID, err = newUUID()
	if err != nil {
		return "", err
	}

	if _, err := store.db.Exec(
		`INSERT INTO metadata(key, value) VALUES(?, ?)`,
		key,
		clientID,
	); err != nil {
		return "", fmt.Errorf("insert client id: %w", err)
	}

	return clientID, nil
}

func (store *Store) MarkReported(requestID string) error {
	_, err := store.db.Exec(
		`INSERT OR REPLACE INTO reported_ids(request_id, reported_at) VALUES(?, ?)`,
		requestID,
		time.Now().Unix(),
	)
	if err != nil {
		return fmt.Errorf("mark reported: %w", err)
	}
	return nil
}

func (store *Store) HasReported(requestID string) (bool, error) {
	var exists int
	err := store.db.QueryRow(
		`SELECT 1 FROM reported_ids WHERE request_id = ? LIMIT 1`,
		requestID,
	).Scan(&exists)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("query reported state: %w", err)
}

func (store *Store) UpdateSyncState(filePath string, lastModified int64, lastLineOffset int) error {
	_, err := store.db.Exec(
		`INSERT OR REPLACE INTO sync_state(file_path, last_modified, last_line_offset) VALUES(?, ?, ?)`,
		filePath,
		lastModified,
		lastLineOffset,
	)
	if err != nil {
		return fmt.Errorf("update sync state: %w", err)
	}
	return nil
}

func (store *Store) ensureSchema() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS sync_state (
			file_path TEXT PRIMARY KEY,
			last_modified INTEGER NOT NULL,
			last_line_offset INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS reported_ids (
			request_id TEXT PRIMARY KEY,
			reported_at INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS metadata (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)`,
	}

	for _, statement := range statements {
		if _, err := store.db.Exec(statement); err != nil {
			return fmt.Errorf("ensure schema: %w", err)
		}
	}

	return nil
}

func newUUID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate client id: %w", err)
	}

	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		raw[0:4],
		raw[4:6],
		raw[6:8],
		raw[8:10],
		raw[10:16],
	), nil
}

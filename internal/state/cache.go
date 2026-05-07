package state

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type Cache struct {
	LastFetch time.Time      `json:"last_fetch"`
	PRStates  map[int]string `json:"pr_states"`
}

func ReadCache(path string) (*Cache, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &Cache{PRStates: map[int]string{}}, nil
		}
		return nil, err
	}
	var c Cache
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if c.PRStates == nil {
		c.PRStates = map[int]string{}
	}
	return &c, nil
}

func WriteCache(path string, c *Cache) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

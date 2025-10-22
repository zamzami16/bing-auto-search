// The following struct declaration is removed as it is unused and causing syntax errors.
package app

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type DelayRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type ScrollRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type GlobalSetting struct {
	Delay       DelayRange  `json:"delay"`
	TotalSearch int         `json:"total_search"`
	Scroll      ScrollRange `json:"scroll"`
	TotalScroll ScrollRange `json:"total_scroll"`
}

type ViewConfig struct {
	PosX int `json:"pos_x"`
	PosY int `json:"pos_y"`
	DX   int `json:"d_x"`
	DY   int `json:"d_y"`
}

type DesktopConfig struct {
	Name    string       `json:"name"`
	Configs []ViewConfig `json:"configs"`
}

// New config structure
type Config struct {
	GlobalSetting GlobalSetting   `json:"global_setting"`
	Data          []DesktopConfig `json:"data"`
}

func DefaultConfig() Config {
	return Config{
		GlobalSetting: GlobalSetting{
			Delay:       DelayRange{Min: 3, Max: 7},
			TotalSearch: 32,
			Scroll:      ScrollRange{Min: 500, Max: 1000},
			TotalScroll: ScrollRange{Min: 2, Max: 5},
		},
		Data: []DesktopConfig{
			{
				Name: "desktop-1",
				Configs: []ViewConfig{
					{PosX: 400, PosY: 50, DX: 150, DY: 5},
				},
			},
		},
	}
}

func LoadConfig(path string) (Config, error) {
	if path == "" {
		return DefaultConfig(), nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig(), err
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return DefaultConfig(), err
	}
	return c, nil
}

func EnsureConfigFile(path string) (string, error) {
	if path == "" {
		path = filepath.Join("config", "config.json")
	}
	_, err := os.Stat(path)
	if err == nil {
		return path, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return path, err
	}

	// tulis default
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return path, err
	}
	def := DefaultConfig()
	b, _ := json.MarshalIndent(def, "", "  ")
	if writeErr := os.WriteFile(path, b, 0o644); writeErr != nil {
		return path, writeErr
	}
	return path, nil
}

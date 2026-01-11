package app

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type MobileDelayRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type MobileScrollConfig struct {
	Chance int `json:"chance"` // 1/N chance to scroll (e.g., 5 means 1/5 chance)
}

type MobileSearchArea struct {
	X      int `json:"x"`       // Base X coordinate
	Y      int `json:"y"`       // Base Y coordinate
	DeltaX int `json:"delta_x"` // Random range for X
	DeltaY int `json:"delta_y"` // Random range for Y
}

type MobileConfig struct {
	Picks           int                `json:"picks"`             // Number of keywords to pick
	WorkUser        int                `json:"work_user"`         // Android work profile user ID
	WaitAppOpen     int                `json:"wait_app_open"`     // Seconds to wait for manual app open
	TapDelay        MobileDelayRange   `json:"tap_delay"`         // Delay after tap (seconds)
	EnterDelay      MobileDelayRange   `json:"enter_delay"`       // Delay before pressing Enter (seconds)
	AfterEnter      int                `json:"after_enter"`       // Delay after Enter (seconds)
	ClearDelay      int                `json:"clear_delay"`       // Delay after clearing text (milliseconds)
	Scroll          MobileScrollConfig `json:"scroll"`            // Scroll configuration
	ScrollStartFrom int                `json:"scroll_start_from"` // Start occasional scroll from cycle N
	SearchArea      MobileSearchArea   `json:"search_area"`       // Search bar tap area
}

func DefaultMobileConfig() MobileConfig {
	return MobileConfig{
		Picks:           22,
		WorkUser:        10,
		WaitAppOpen:     10,
		TapDelay:        MobileDelayRange{Min: 1, Max: 2},
		EnterDelay:      MobileDelayRange{Min: 1, Max: 2},
		AfterEnter:      3,
		ClearDelay:      150,
		Scroll:          MobileScrollConfig{Chance: 5},
		ScrollStartFrom: 2,
		SearchArea: MobileSearchArea{
			X:      150,
			Y:      130,
			DeltaX: 570,
			DeltaY: 95,
		},
	}
}

func LoadMobileConfig(path string) (MobileConfig, error) {
	if path == "" {
		return DefaultMobileConfig(), nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return DefaultMobileConfig(), err
	}
	var c MobileConfig
	if err := json.Unmarshal(b, &c); err != nil {
		return DefaultMobileConfig(), err
	}
	return c, nil
}

func EnsureMobileConfigFile(path string) (string, error) {
	if path == "" {
		path = filepath.Join("config", "mobile-config.json")
	}
	_, err := os.Stat(path)
	if err == nil {
		return path, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return path, err
	}

	// Write default config
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return path, err
	}
	def := DefaultMobileConfig()
	b, _ := json.MarshalIndent(def, "", "  ")
	if writeErr := os.WriteFile(path, b, 0o644); writeErr != nil {
		return path, writeErr
	}
	return path, nil
}

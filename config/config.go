package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigDir())

	viper.SetEnvPrefix("JELLYFIN")
	viper.AutomaticEnv()

	viper.SetDefault("url", "http://192.168.1.69:8096")

	_ = viper.ReadInConfig()
}

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "jellyfin-cli")
}

func URL() string {
	return viper.GetString("url")
}

func Save() error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	return viper.WriteConfigAs(filepath.Join(dir, "config.yaml"))
}

// Session stores auth info.
type Session struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

func sessionPath() string {
	return filepath.Join(ConfigDir(), "session.json")
}

func LoadSession() (*Session, error) {
	data, err := os.ReadFile(sessionPath())
	if err != nil {
		return nil, fmt.Errorf("not authenticated, run 'jellyfin login'")
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("corrupt session file: %w", err)
	}
	return &s, nil
}

func SaveSession(s *Session) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sessionPath(), data, 0600)
}

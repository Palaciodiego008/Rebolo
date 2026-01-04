package adapters

import (
	"os"
	"gopkg.in/yaml.v3"
	"github.com/Palaciodiego008/rebololang/pkg/rebolo/ports"
)

// YAMLConfig implements ConfigPort
type YAMLConfig struct{}

func NewYAMLConfig() *YAMLConfig {
	return &YAMLConfig{}
}

func (c *YAMLConfig) Load() (ports.ConfigData, error) {
	config := ports.ConfigData{}
	
	// Set defaults
	config.Server.Port = c.GetEnv("PORT", "3000")
	config.Server.Host = c.GetEnv("HOST", "localhost")
	config.App.Env = c.GetEnv("REBOLO_ENV", "development")
	config.Assets.HotReload = config.App.Env == "development"
	
	// Try to load config.yml
	if data, err := os.ReadFile("config.yml"); err == nil {
		yaml.Unmarshal(data, &config)
	}
	
	return config, nil
}

func (c *YAMLConfig) GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

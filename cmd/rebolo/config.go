package main

import "time"

// DevConfig holds configuration for development server
type DevConfig struct {
	// Go server settings
	GoRestartDebounce time.Duration
	GoWatchExtensions []string
	GoSkipDirs        []string

	// Frontend settings
	FrontendWatchExtensions []string
	FrontendSrcDir          string
	FrontendOutDir          string

	// Bun settings
	BunInstallCommand []string
	BunWatchCommand   []string
	BunBuildCommand   []string
}

// DefaultDevConfig returns the default development configuration
func DefaultDevConfig() *DevConfig {
	return &DevConfig{
		GoRestartDebounce:       100 * time.Millisecond,
		GoWatchExtensions:       []string{".go"},
		GoSkipDirs:              []string{"node_modules", ".git", "vendor", "public", "dist"},
		FrontendWatchExtensions: []string{".js", ".css", ".ts", ".jsx", ".tsx"},
		FrontendSrcDir:          "src",
		FrontendOutDir:          "public",
		BunInstallCommand:       []string{"bun", "install"},
		BunWatchCommand:         []string{"bun", "run", "watch"},
		BunBuildCommand:         []string{"bun", "run", "build"},
	}
}

// FieldTypeMapping defines mappings between different type systems
type FieldTypeMapping struct {
	GoTypes   map[string]string
	SQLTypes  map[string]string
	HTMLTypes map[string]string
}

// DefaultFieldTypeMapping returns default type mappings
func DefaultFieldTypeMapping() *FieldTypeMapping {
	return &FieldTypeMapping{
		GoTypes: map[string]string{
			"string":   "string",
			"text":     "string",
			"int":      "int64",
			"integer":  "int64",
			"bool":     "bool",
			"boolean":  "bool",
			"float":    "float64",
			"time":     "time.Time",
			"datetime": "time.Time",
		},
		SQLTypes: map[string]string{
			"string":   "VARCHAR(255)",
			"text":     "TEXT",
			"int":      "BIGINT",
			"integer":  "BIGINT",
			"bool":     "BOOLEAN",
			"boolean":  "BOOLEAN",
			"float":    "DECIMAL",
			"time":     "TIMESTAMP",
			"datetime": "TIMESTAMP",
		},
		HTMLTypes: map[string]string{
			"string":   "text",
			"text":     "textarea",
			"int":      "number",
			"integer":  "number",
			"bool":     "checkbox",
			"boolean":  "checkbox",
			"float":    "number",
			"time":     "datetime-local",
			"datetime": "datetime-local",
		},
	}
}

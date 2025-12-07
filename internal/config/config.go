package config

import (
	"os"
	"strconv"
)

type Config struct {
	APIURL       string
	ItemsPerPage int
	DataFilePath string
}

// Load fetches configuration from env vars or sets defaults
func Load() Config {
	return Config{
		APIURL:       getEnv("ASI67_API_URL", "https://www.asi67.com/webapi/getJson/Templates/ProductsList"),
		ItemsPerPage: getEnvAsInt("ASI67_ITEMS_PER_PAGE", 12),
		DataFilePath: getEnv("DATA_FILE_PATH", "data/seen.json"),
	}
}

// Helper to get env string with default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper to get env int with default
func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

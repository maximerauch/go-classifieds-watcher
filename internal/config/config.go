package config

import (
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	Asi67      Asi67Config
	RememberMe RememberMeConfig
	Email      EmailConfig
}

type Asi67Config struct {
	APIURL       string
	ItemsPerPage int
	DataFilePath string
}

type RememberMeConfig struct {
	SearchURL    string
	DataFilePath string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	From         string
	To           []string
}

// Load fetches configuration from env vars or sets defaults
func Load() AppConfig {
	return AppConfig{
		Asi67: Asi67Config{
			APIURL:       getEnv("ASI67_API_URL", "https://www.asi67.com/webapi/getJson/Templates/ProductsList"),
			ItemsPerPage: getEnvAsInt("ASI67_ITEMS_PER_PAGE", 12),
			DataFilePath: getEnv("ASI67_DATA_FILE_PATH", "data/asi67-seen.json"),
		},

		RememberMe: RememberMeConfig{
			SearchURL:    getEnv("REMEMBERME_SEARCH_URL", "https://remembermefrance.org/pets/?breed=0&pets_search%5Bsexe%5D=all&pets_search%5Bou_est_le_chien%5D=all&pets_search%5Burgence%5D=all"),
			DataFilePath: getEnv("REMEMBERME_DATA_FILE_PATH", "data/rememberme-seen.json"),
		},

		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
			SMTPUser:     getEnv("SMTP_USER", "your.mail@gmail.com"),
			SMTPPassword: getEnv("SMTP_PASSWORD", "your-password"),
			From:         getEnv("EMAIL_FROM", "Watcher Bot <your.mail@gmail.com>"),
			To:           strings.Split(getEnv("EMAIL_TO", "your.mail@gmail.com"), ","),
		},
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

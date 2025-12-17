package config

import (
	"testing"
)

// TestLoad validates the configuration loading logic.
// It ensures environment variables take precedence over defaults
// and that type conversion (string -> int, string -> slice) works as expected.
func TestLoad(t *testing.T) {

	// Case 1: Happy Path
	// Verify that environment variables correctly override default values.
	t.Run("Overrides defaults with environment variables", func(t *testing.T) {
		// Arrange: Inject environment variables.
		t.Setenv("ASI67_API_URL", "https://api.test.com")
		t.Setenv("ASI67_ITEMS_PER_PAGE", "42")                // Tests string-to-int conversion
		t.Setenv("EMAIL_TO", "user1@test.com,user2@test.com") // Tests string splitting
		t.Setenv("SMTP_PORT", "2525")
		t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")

		// Act: Execute the Load function
		cfg := Load()

		// Assert: Verify the resulting struct

		// Verify simple string mapping
		if cfg.Asi67.APIURL != "https://api.test.com" {
			t.Errorf("Asi67.APIURL = %s; want https://api.test.com", cfg.Asi67.APIURL)
		}

		// Verify integer conversion
		if cfg.Asi67.ItemsPerPage != 42 {
			t.Errorf("Asi67.ItemsPerPage = %d; want 42", cfg.Asi67.ItemsPerPage)
		}

		// Verify slice logic (CSV splitting)
		if len(cfg.Email.To) != 2 {
			t.Errorf("Email.To len = %d; want 2", len(cfg.Email.To))
		}
		if cfg.Email.To[0] != "user1@test.com" || cfg.Email.To[1] != "user2@test.com" {
			t.Errorf("Email.To parsing failed, got %v", cfg.Email.To)
		}

		// Verify nested struct mapping
		if cfg.Database.DSN != "postgres://user:pass@localhost:5432/db" {
			t.Errorf("Database.DSN mismatch")
		}
	})

	// Case 2: Edge Case / Resilience
	// Verify that the application falls back to default values instead of crashing
	// when an invalid type is provided (e.g., text for an integer field).
	t.Run("Uses default fallback on invalid integer", func(t *testing.T) {
		// Arrange: Inject a non-numeric string for a numeric field
		t.Setenv("SMTP_PORT", "im-not-a-number")

		// Act
		cfg := Load()

		// Assert: The loader should swallow the error and use the default value (587).
		expectedDefault := 587
		if cfg.Email.SMTPPort != expectedDefault {
			t.Errorf("SMTPPort = %d; want default %d because input was invalid", cfg.Email.SMTPPort, expectedDefault)
		}
	})
}

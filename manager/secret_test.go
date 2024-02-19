package manager

import (
	"os"
	"testing"
)

// Helper function to create a .secret file and write a secret to it
func createDotSecretFile(secretValue string) error {
	return os.WriteFile(DOT_SECRET, []byte(secretValue), 0644)
}

// Helper function to clean up by removing the .secret file
func cleanupDotSecretFile() {
	os.Remove(DOT_SECRET)
}

// TestInitSecretFromFile tests InitSecret function when the secret is read from a file
func TestInitSecretFromFile(t *testing.T) {
	expectedSecret := "fileSecret"
	err := createDotSecretFile(expectedSecret)
	if err != nil {
		t.Fatalf("Unable to create .secret file: %v", err)
	}
	defer cleanupDotSecretFile()

	s := InitSecret()

	if s.GetSecret() != expectedSecret {
		t.Errorf("Expected secret %s, got %s", expectedSecret, s.GetSecret())
	}
}

// TestInitSecretFromEnv tests InitSecret function when the secret is read from an environment variable
func TestInitSecretFromEnv(t *testing.T) {
	expectedSecret := "envSecret"
	os.Setenv("ENV_SECRET", expectedSecret)
	defer os.Unsetenv("ENV_SECRET")

	// Ensure no .secret file exists to test environment variable functionality
	cleanupDotSecretFile()

	s := InitSecret()

	if s.GetSecret() != expectedSecret {
		t.Errorf("Expected secret %s, got %s", expectedSecret, s.GetSecret())
	}
}

// Additional tests can be added to cover more scenarios, such as when neither the file nor the environment variable is set.

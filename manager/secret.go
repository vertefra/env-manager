package manager

import (
	"io"
	"os"
	"strings"
)

const DOT_SECRET = ".secret"
const ENV_SECRET = "ENV_MANAGER_SECRET"

// Search the secret to encypt / decrypt the files
// It can be in a file `.secret` or it can be passed
type secret struct {
	secret string
}

func (s *secret) GetSecret() string {
	return s.secret
}

func getSecretFromFile() *string {
	if _, err := os.Stat(DOT_SECRET); os.IsNotExist(err) {
		return nil
	}

	/// If found, read the file
	/// and set the secret
	f, err := os.Open(DOT_SECRET)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	fileContent, err := io.ReadAll(f)

	if err != nil {
		panic(err)
	}

	secret := string(fileContent)
	secret = strings.Trim(secret, " ")

	if secret == "" {
		return nil
	}

	return &secret
}

func getSecretFromEnv() *string {
	secret := os.Getenv(ENV_SECRET)

	if secret == "" {
		return nil
	}

	return &secret
}

func (s *secret) findSecret() {
	/// Only support secret from file for now
	/// If not found, search if a file is present
	_secret := getSecretFromEnv()

	if _secret == nil {
		_secret = getSecretFromFile()
	}

	if _secret == nil {
		panic("No secret found")
	}

	s.secret = *_secret
}

func InitSecret() secret {
	s := secret{}
	s.findSecret()
	return s
}

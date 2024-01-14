package manager

import (
	"io"
	"os"
	"strings"
)

const DOT_SECRET = ".secret"

// Search the secret to encypt / decrypt the files
// It can be in a file `.secret` or it can be passed
type secret struct {
	secret string
}

func (s *secret) GetSecret() string {
	return s.secret
}

func (s *secret) findSecret() {
	/// Only support secret from file for now
	/// If not found, search if a file is present
	if _, err := os.Stat(DOT_SECRET); os.IsNotExist(err) {
		panic("Secret not found")
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

	s.secret = string(fileContent)
	s.secret = strings.Trim(s.secret, " ")
}

func InitSecret() secret {
	s := secret{}
	s.findSecret()
	return s
}

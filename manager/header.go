package manager

import (
	"errors"
	"strings"
)

type Header struct {
	Identifier string
	RestoreAs  string
}

func (h *Header) String() []string {
	return []string{
		h.Identifier,
		h.RestoreAs,
	}
}

func InitHeader(text string) (*Header, error) {
	var identifier string = ""
	var restoreAs string = ""

	lines := strings.Split(text, "\n")

	for _, line := range lines {
		// Sanitize line
		line = strings.Trim(line, " ")

		if strings.HasPrefix(line, IDENTIFIER_HEADER) {
			identifier = strings.TrimPrefix(line, IDENTIFIER_HEADER)
		}
		if strings.HasPrefix(line, RESTORE_AS_HEADER) {
			restoreAs = strings.TrimPrefix(line, RESTORE_AS_HEADER)
		}
	}

	if identifier == "" {
		return nil, errors.New("invalid header: identifier not found")
	}

	if restoreAs == "" {
		return nil, errors.New("invalid header: restore-as not found")
	}

	return &Header{
		Identifier: identifier,
		RestoreAs:  restoreAs,
	}, nil
}

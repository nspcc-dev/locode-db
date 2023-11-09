package locodedb

import (
	"errors"
	"fmt"
)

// LocationCodeLen is the length of the location code.
const LocationCodeLen = 3

// ErrInvalidString is returned when the string is not a valid location code.
var ErrInvalidString = errors.New("invalid string format in UN/Locode")

// LocationCode represents a location code for
// the storage in the location database.
type LocationCode [LocationCodeLen]uint8

// LocationCodeFromString parses a string and returns the location code.
func LocationCodeFromString(s string) (*LocationCode, error) {
	if l := len(s); l != LocationCodeLen {
		return nil, fmt.Errorf("incorrect location code length: expect: %d, got: %d",
			LocationCodeLen,
			l,
		)
	}

	for i := range s {
		if !isUpperAlpha(s[i]) && !isDigit(s[i]) {
			return nil, ErrInvalidString
		}
	}

	lc := LocationCode{}
	copy(lc[:], s)

	return &lc, nil
}

func isDigit(sym uint8) bool {
	return sym >= '0' && sym <= '9'
}

func isUpperAlpha(sym uint8) bool {
	return sym >= 'A' && sym <= 'Z'
}

// String returns a string representation of the location code.
func (l *LocationCode) String() string {
	syms := l.Symbols()
	return string(syms[:])
}

// Symbols returns the location code as a slice of symbols.
func (l *LocationCode) Symbols() [LocationCodeLen]uint8 {
	return *l
}

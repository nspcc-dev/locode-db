package locodedb

import (
	"errors"
	"fmt"
)

// CountryCodeLen is the length of the country code.
const CountryCodeLen = 2

// LocationCodeLen is the length of the location code.
const LocationCodeLen = 3

// ErrInvalidString is returned when the string is not a valid location code.
var ErrInvalidString = errors.New("invalid string format in UN/Locode")

// countryCode represents ISO 3166 alpha-2 Country Code.
type countryCode [CountryCodeLen]uint8

// countryCodeFromString parses a string and returns the country code.
func countryCodeFromString(s string) (*countryCode, error) {
	if l := len(s); l != CountryCodeLen {
		return nil, fmt.Errorf("incorrect country code length: expect: %d, got: %d",
			CountryCodeLen,
			l,
		)
	}

	for i := range s {
		if !isUpperAlpha(s[i]) {
			return nil, ErrInvalidString
		}
	}

	cc := countryCode{}
	copy(cc[:], s)

	return &cc, nil
}

func isUpperAlpha(sym uint8) bool {
	return sym >= 'A' && sym <= 'Z'
}

func isDigit(sym uint8) bool {
	return sym >= '0' && sym <= '9'
}

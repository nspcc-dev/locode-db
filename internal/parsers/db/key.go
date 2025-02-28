package locodedb

import (
	"fmt"

	"github.com/nspcc-dev/locode-db/pkg/locodedb"
)

// Key represents the Key in location database. It contains the country code and the location code.
type Key struct {
	cc string
	lc string
}

// NewKey returns a new Key from a country code and a location code string pair (e.g. "USNYC") or an error if
// the string is invalid. The country code must be 2 letters long and the location code 3 letters long.
func NewKey(country, location string) (*Key, error) {
	err := validateCode(country, locodedb.CountryCodeLen)
	if err != nil {
		return nil, fmt.Errorf("could not parse country: %w", err)
	}

	err = validateCode(location, locodedb.LocationCodeLen)
	if err != nil {
		return nil, fmt.Errorf("could not parse location: %w", err)
	}

	return &Key{
		cc: country,
		lc: location,
	}, nil
}

// CountryCode returns the location's country code in string representation.
func (k *Key) CountryCode() string {
	return k.cc
}

// LocationCode returns the location code in string representation.
func (k *Key) LocationCode() string {
	return k.lc
}

// validateCode validates if code is.
func validateCode(s string, codeLen int) error {
	if l := len(s); l != codeLen {
		return fmt.Errorf("incorrect location code length: expect: %d, got: %d",
			codeLen,
			l,
		)
	}

	for i := range s {
		if !isUpperAlpha(s[i]) && !isDigit(s[i]) {
			return locodedb.ErrInvalidString
		}
	}

	return nil
}

func isUpperAlpha(sym uint8) bool {
	return sym >= 'A' && sym <= 'Z'
}

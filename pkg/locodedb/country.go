package locodedb

import (
	"fmt"
)

// CountryCodeLen is the length of the country code.
const CountryCodeLen = 2

// CountryCode represents ISO 3166 alpha-2 Country Code.
type CountryCode [CountryCodeLen]uint8

// CountryCodeFromString parses a string and returns the country code.
func CountryCodeFromString(s string) (*CountryCode, error) {
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

	cc := CountryCode{}
	copy(cc[:], s)

	return &cc, nil
}

// String returns a string representation of the country code.
func (c *CountryCode) String() string {
	syms := c.Symbols()
	return string(syms[:])
}

// Symbols returns the country code as a slice of symbols.
func (c *CountryCode) Symbols() [CountryCodeLen]uint8 {
	return *c
}

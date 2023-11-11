package locodedb

import (
	"fmt"
)

// Key represents the key in location database. It contains the country code and the location code.
type Key struct {
	cc *CountryCode
	lc *LocationCode
}

// NewKey returns a new Key from a country code and a location code string pair (e.g. "USNYC") or an error if
// the string is invalid. The country code must be 2 letters long and the location code 3 letters long.
func NewKey(country, location string) (*Key, error) {
	countryCode, err := CountryCodeFromString(country)
	if err != nil {
		return nil, fmt.Errorf("could not parse country: %w", err)
	}

	locationCode, err := LocationCodeFromString(location)
	if err != nil {
		return nil, fmt.Errorf("could not parse location: %w", err)
	}

	return &Key{
		cc: countryCode,
		lc: locationCode,
	}, nil
}

// CountryCode returns the location's country code.
func (k *Key) CountryCode() *CountryCode {
	return k.cc
}

// LocationCode returns the location code.
func (k *Key) LocationCode() *LocationCode {
	return k.lc
}

// Record represents a record in the location database (resulting CSV files). It contains all the
// information about the location. It is used to fill the database. Country, Location are full names, codes are in Key.
type Record struct {
	Country    string
	Location   string
	SubDivName string
	SubDivCode string
	Point      Point
	Cont       Continent
}

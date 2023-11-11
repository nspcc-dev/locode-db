package locodedb

import (
	"errors"
)

// ErrNotFound is returned when the record is not found in the location database.
var ErrNotFound = errors.New("record not found")

var (
	// mCountries is a map of country codes to country names.
	mCountries map[CountryCode]string

	// mLocodes is a map of location codes to location records. The location code is a concatenation of the country code
	// and the location code. The Record contains country name, the location name, the subdivision name, the subdivision
	// code, the Point and the Continent.
	mLocodes map[string]locodesCSV
)

// Get returns a record for a given locode string. The string must be 5 letters long. The first 2 letters are the country
// code and the last 3 letters are the location code.
func Get(locodeStr string) (Record, error) {
	if err := initLocodeData(); err != nil {
		return Record{}, err
	}

	if len(locodeStr) == CountryCodeLen+LocationCodeLen+1 && locodeStr[CountryCodeLen] == ' ' {
		locodeStr = locodeStr[:CountryCodeLen] + locodeStr[CountryCodeLen+1:]
	}
	if len(locodeStr) != CountryCodeLen+LocationCodeLen {
		return Record{}, ErrInvalidString
	}

	for i := range locodeStr[:CountryCodeLen] {
		if !isUpperAlpha(locodeStr[i]) {
			return Record{}, ErrInvalidString
		}
	}
	for i := range locodeStr[CountryCodeLen:] {
		if !isUpperAlpha(locodeStr[CountryCodeLen+i]) && !isDigit(locodeStr[CountryCodeLen+i]) {
			return Record{}, ErrInvalidString
		}
	}

	locodeCSV, found := mLocodes[locodeStr]
	if !found {
		return Record{}, ErrNotFound
	}

	cc := CountryCode{}
	copy(cc[:], locodeStr[:2])
	country, countryFound := mCountries[cc]
	if !countryFound {
		return Record{}, ErrNotFound
	}

	return Record{
		Country:    country,
		Location:   locodeCSV.locationName,
		SubDivName: locodeCSV.subDivName,
		SubDivCode: locodeCSV.subDivCode,
		Point:      locodeCSV.point,
		Cont:       locodeCSV.continent,
	}, nil
}

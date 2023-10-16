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

	key, err := NewKey(locodeStr[:2], locodeStr[2:])
	if err != nil {
		return Record{}, err
	}
	return getFromKey(*key)
}

// getFromKey returns a record for a given Key.
func getFromKey(key Key) (Record, error) {
	newlocode := key.CountryCode().String() + key.LocationCode().String()

	locodeCSV, found := mLocodes[newlocode]
	if !found {
		return Record{}, ErrNotFound
	}

	country, countryFound := mCountries[*key.CountryCode()]
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

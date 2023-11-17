package locodedb

import (
	"errors"
	"sort"
	"strings"
)

// ErrNotFound is returned when the record is not found in the location database.
var ErrNotFound = errors.New("record not found")

var (
	// locodeStrings is a string containing all substrings of locode data.
	locodeStrings string

	// mCountries is a map of country codes to country names and locodes.
	mCountries map[CountryCode]countryData
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

	cc := CountryCode{}
	copy(cc[:], locodeStr[:2])
	cd, countryFound := mCountries[cc]
	if !countryFound {
		return Record{}, ErrNotFound
	}

	code := locodeStr[CountryCodeLen:]
	n := sort.Search(len(cd.locodes), func(i int) bool {
		cmp := strings.Compare(codeFromCSV(&cd.locodes[i]), code)
		return cmp >= 0
	})
	if n == len(cd.locodes) || strings.Compare(codeFromCSV(&cd.locodes[n]), code) != 0 {
		return Record{}, ErrNotFound
	}

	return Record{
		Country:    cd.name,
		Location:   locFromCSV(&cd.locodes[n]),
		SubDivName: divNameFromCSV(&cd.locodes[n]),
		SubDivCode: divCodeFromCSV(&cd.locodes[n]),
		Point:      cd.locodes[n].point,
		Cont:       cd.locodes[n].continent,
	}, nil
}

func codeFromCSV(c *locodesCSV) string {
	return string(c.code[:])
}

func locFromCSV(c *locodesCSV) string {
	return locodeStrings[c.offset : c.offset+uint32(c.locationLen)]
}

func divCodeFromCSV(c *locodesCSV) string {
	return locodeStrings[c.offset+uint32(c.locationLen) : c.offset+uint32(c.locationLen)+uint32(c.subDivCodeLen)]
}

func divNameFromCSV(c *locodesCSV) string {
	return locodeStrings[c.offset+uint32(c.locationLen)+uint32(c.subDivCodeLen) : c.offset+uint32(c.locationLen)+uint32(c.subDivCodeLen)+uint32(c.subDivNameLen)]
}

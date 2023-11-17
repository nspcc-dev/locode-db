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

	n := sort.Search(len(cd.locodes), func(i int) bool {
		cmp := strings.Compare(codeFromCSV(&cd.locodes[i]), locodeStr[CountryCodeLen:])
		return cmp >= 0
	})
	if n == len(cd.locodes) || strings.Compare(codeFromCSV(&cd.locodes[n]), locodeStr[CountryCodeLen:]) != 0 {
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
	return locodeStrings[int(c.offset) : int(c.offset)+LocationCodeLen]
}

func locFromCSV(c *locodesCSV) string {
	return locodeStrings[int(c.offset)+LocationCodeLen : int(c.offset)+LocationCodeLen+int(c.locationLen)]
}

func divCodeFromCSV(c *locodesCSV) string {
	return locodeStrings[int(c.offset)+LocationCodeLen+int(c.locationLen) : int(c.offset)+LocationCodeLen+int(c.locationLen)+int(c.subDivCodeLen)]
}

func divNameFromCSV(c *locodesCSV) string {
	return locodeStrings[int(c.offset)+LocationCodeLen+int(c.locationLen)+int(c.subDivCodeLen) : int(c.offset)+LocationCodeLen+int(c.locationLen)+int(c.subDivCodeLen)+int(c.subDivNameLen)]
}

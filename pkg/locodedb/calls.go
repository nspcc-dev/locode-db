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

	// mCountries is a map of country codes to country names.
	mCountries map[CountryCode]string

	// mLocodes is a slice of location code data.
	mLocodes []locodesCSV
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

	n := sort.Search(len(mLocodes), func(i int) bool {
		cmp := strings.Compare(nameFromCSV(&mLocodes[i]), locodeStr)
		return cmp >= 0
	})
	if n == len(mLocodes) || strings.Compare(nameFromCSV(&mLocodes[n]), locodeStr) != 0 {
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
		Location:   locFromCSV(&mLocodes[n]),
		SubDivName: divNameFromCSV(&mLocodes[n]),
		SubDivCode: divCodeFromCSV(&mLocodes[n]),
		Point:      mLocodes[n].point,
		Cont:       mLocodes[n].continent,
	}, nil
}

func nameFromCSV(c *locodesCSV) string {
	return locodeStrings[int(c.offset) : int(c.offset)+CountryCodeLen+LocationCodeLen]
}

func locFromCSV(c *locodesCSV) string {
	return locodeStrings[int(c.offset)+CountryCodeLen+LocationCodeLen : int(c.offset)+CountryCodeLen+LocationCodeLen+int(c.locationLen)]
}

func divCodeFromCSV(c *locodesCSV) string {
	return locodeStrings[int(c.offset)+CountryCodeLen+LocationCodeLen+int(c.locationLen) : int(c.offset)+CountryCodeLen+LocationCodeLen+int(c.locationLen)+int(c.subDivCodeLen)]
}

func divNameFromCSV(c *locodesCSV) string {
	return locodeStrings[int(c.offset)+CountryCodeLen+LocationCodeLen+int(c.locationLen)+int(c.subDivCodeLen) : int(c.offset)+CountryCodeLen+LocationCodeLen+int(c.locationLen)+int(c.subDivCodeLen)+int(c.subDivNameLen)]
}

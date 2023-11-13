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
		cmp := strings.Compare(locodeSubstr(mLocodes[i].locode), locodeStr)
		return cmp >= 0
	})
	if n == len(mLocodes) || strings.Compare(locodeSubstr(mLocodes[n].locode), locodeStr) != 0 {
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
		Location:   locodeSubstr(mLocodes[n].locationName),
		SubDivName: locodeSubstr(mLocodes[n].subDivName),
		SubDivCode: locodeSubstr(mLocodes[n].subDivCode),
		Point:      mLocodes[n].point,
		Cont:       mLocodes[n].continent,
	}, nil
}

func locodeSubstr(ol offLen) string {
	return locodeStrings[int(ol.offset) : int(ol.offset)+int(ol.length)]
}

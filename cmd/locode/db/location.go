package locodedb

import (
	"fmt"

	locodecolumn "github.com/nspcc-dev/locode-db/cmd/locode/column"
)

// LocationCode represents a location code for
// the storage in the NeoFS location database.
type LocationCode locodecolumn.LocationCode

// LocationCodeFromString parses a string UN/LOCODE location code
// and returns a LocationCode.
func LocationCodeFromString(s string) (*LocationCode, error) {
	lc, err := locodecolumn.LocationCodeFromString(s)
	if err != nil {
		return nil, fmt.Errorf("could not parse location code: %w", err)
	}

	return LocationFromColumn(lc)
}

// LocationFromColumn converts a UN/LOCODE country code to a LocationCode.
func LocationFromColumn(cc *locodecolumn.LocationCode) (*LocationCode, error) {
	return (*LocationCode)(cc), nil
}

func (l *LocationCode) String() string {
	syms := (*locodecolumn.LocationCode)(l).Symbols()
	return string(syms[:])
}

package locodedb

import (
	"errors"
	"fmt"

	"github.com/nspcc-dev/locode-db/cmd/locode"
)

// SourceTable is an interface of the UN/LOCODE table.
type SourceTable interface {
	// Must iterate over all entries of the table
	// and pass next entry to the handler.
	//
	// Must return handler's errors directly.
	IterateAll(func(locode.Record) error) error
}

// DB is an interface of NeoFS location database.
type DB interface {
	// Must save the record by key in the database.
	Put(Key, Record) error

	// Must return the record by key from the database.
	Get(Key) (*Record, error)
}

// AirportRecord represents the entry in NeoFS airport database.
type AirportRecord struct {
	// Name of the country where airport is located.
	CountryName string

	// Geo point where airport is located.
	Point *Point
}

// ErrAirportNotFound is returned by AirportRecord readers
// when the required airport is not found.
var ErrAirportNotFound = errors.New("airport not found")

// AirportDB is an interface of NeoFS airport database.
type AirportDB interface {
	// Must return the record by UN/LOCODE table record.
	//
	// Must return ErrAirportNotFound if there is no
	// related airport in the database.
	Get(locode.Record) (*AirportRecord, error)
}

// ContinentsDB is an interface of NeoFS continent database.
type ContinentsDB interface {
	// Must return continent of the geo point.
	PointContinent(*Point) (*Continent, error)
}

var ErrSubDivNotFound = errors.New("subdivision not found")

var ErrCountryNotFound = errors.New("country not found")

// NamesDB is an interface of the NeoFS location namespace.
type NamesDB interface {
	// Must resolve a country code to a country name.
	//
	// Must return ErrCountryNotFound if there is no
	// country with the provided code.
	CountryName(*CountryCode) (string, error)

	// Must resolve (country code, subdivision code) to
	// a subdivision name.
	//
	// Must return ErrSubDivNotFound if either country or
	// subdivision is not presented in database.
	SubDivName(*CountryCode, string) (string, error)
}

// FillDatabase generates the NeoFS location database based on the UN/LOCODE table.
func FillDatabase(table SourceTable, airports AirportDB, continents ContinentsDB, names NamesDB, db DB) error {
	return table.IterateAll(func(tableRecord locode.Record) error {
		if tableRecord.LOCODE.LocationCode() == "" {
			return nil
		}

		dbKey, err := NewKey(tableRecord.LOCODE)
		if err != nil {
			return err
		}

		dbRecord, err := NewRecord(tableRecord)
		if err != nil {
			if errors.Is(err, errParseCoordinates) {
				return nil
			}

			return err
		}

		geoPoint := dbRecord.GeoPoint()
		countryName := ""

		if geoPoint == nil {
			airportRecord, err := airports.Get(tableRecord)
			if err != nil {
				if errors.Is(err, ErrAirportNotFound) {
					return nil
				}

				return err
			}

			geoPoint = airportRecord.Point
			countryName = airportRecord.CountryName
		}

		dbRecord.SetGeoPoint(geoPoint)

		if countryName == "" {
			countryName, err = names.CountryName(dbKey.CountryCode())
			if err != nil {
				if errors.Is(err, ErrCountryNotFound) {
					return nil
				}

				return err
			}
		}

		dbRecord.SetCountryName(countryName)

		if subDivCode := dbRecord.SubDivCode(); subDivCode != "" {
			subDivName, err := names.SubDivName(dbKey.CountryCode(), subDivCode)
			if err != nil {
				if errors.Is(err, ErrSubDivNotFound) {
					return nil
				}

				return err
			}

			dbRecord.SetSubDivName(subDivName)
		}

		continent, err := continents.PointContinent(geoPoint)
		if err != nil {
			return fmt.Errorf("could not calculate continent geo point: %w", err)
		} else if continent.Is(ContinentUnknown) {
			return nil
		}

		dbRecord.SetContinent(continent)

		return db.Put(*dbKey, *dbRecord)
	})
}

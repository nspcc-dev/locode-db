package locodedb

import (
	"errors"
	"fmt"

	"github.com/nspcc-dev/locode-db/pkg/locodedb"
)

// SourceTable is an interface of the UN/LOCODE table.
type SourceTable interface {
	// IterateAll must iterate over all entries of the table
	// and pass next entry to the handler.
	//
	// Must return handler's errors directly.
	IterateAll(func(Record) error) error
}

// AirportRecord represents the entry in airport database.
type AirportRecord struct {
	// Name of the country where airport is located.
	CountryName string

	// Geo point where airport is located.
	Point locodedb.Point
}

// ErrAirportNotFound is returned by AirportRecord readers
// when the required airport is not found.
var ErrAirportNotFound = errors.New("airport not found")

// AirportDB is an interface of airport database.
type AirportDB interface {
	// Get must return the record by UN/LOCODE table record.
	//
	// Must return ErrAirportNotFound if there is no
	// related airport in the database.
	Get(Record) (*AirportRecord, error)
}

// ContinentsDB is an interface of continent database.
type ContinentsDB interface {
	// PointContinent must return continent of the geo point.
	PointContinent(locodedb.Point) (*locodedb.Continent, error)
}

var ErrSubDivNotFound = errors.New("subdivision not found")

var ErrCountryNotFound = errors.New("country not found")

// NamesDB is an interface of the location namespace.
type NamesDB interface {
	// CountryName must resolve a country code to a country name.
	//
	// Must return ErrCountryNotFound if there is no
	// country with the provided code.
	CountryName(string) (string, error)

	// SubDivName must resolve (country code, subdivision code) to
	// a subdivision name.
	//
	// Must return ErrSubDivNotFound if either country or
	// subdivision is not presented in database.
	SubDivName(string, string) (string, error)
}

// FillDatabase generates the location database based on the UN/LOCODE table.
func FillDatabase(table SourceTable, airports AirportDB, continents ContinentsDB, names NamesDB, db CsvDB) error {
	var newData []Data
	if err := table.IterateAll(func(tableRecord Record) error {
		if tableRecord.LOCODE[1] == "" {
			return nil
		}

		dbKey, err := NewKey(tableRecord.LOCODE[0], tableRecord.LOCODE[1])
		if err != nil {
			return err
		}

		crd, err := CoordinatesFromString(tableRecord.Coordinates)
		if err != nil {
			if errors.Is(err, locodedb.ErrInvalidString) {
				return nil
			}

			return err
		}

		geoPoint, err := PointFromCoordinates(crd)
		if err != nil {
			return fmt.Errorf("could not parse geo point: %w", err)
		}

		countryName := ""

		if geoPoint == (locodedb.Point{}) {
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

		dbRecord := locodedb.Record{
			Location:   tableRecord.NameWoDiacritics,
			SubDivCode: tableRecord.SubDiv,
			Point:      geoPoint,
		}

		if countryName == "" {
			countryName, err = names.CountryName(dbKey.CountryCode())
			if err != nil {
				if errors.Is(err, ErrCountryNotFound) {
					return nil
				}

				return err
			}
		}

		dbRecord.Country = countryName

		if subDivCode := dbRecord.SubDivCode; subDivCode != "" {
			subDivName, err := names.SubDivName(dbKey.CountryCode(), subDivCode)
			if err != nil {
				if errors.Is(err, ErrSubDivNotFound) {
					return nil
				}

				return err
			}

			dbRecord.SubDivName = subDivName
		}

		continent, err := continents.PointContinent(geoPoint)
		if err != nil {
			return fmt.Errorf("could not calculate continent geo point: %w", err)
		} else if *continent == locodedb.ContinentUnknown {
			return nil
		}

		dbRecord.Cont = *continent

		newData = append(newData, Data{*dbKey, dbRecord})

		return nil
	}); err != nil {
		return err
	}

	if err := db.Put(newData); err != nil {
		return err
	}

	return nil
}

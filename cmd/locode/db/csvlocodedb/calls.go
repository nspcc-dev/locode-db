package csvlocodedb

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"

	locodedb "github.com/nspcc-dev/locode-db/cmd/locode/db"
)

var errRecordNotFound = errors.New("record not found")

// Get returns a record by key. locodedb.Key is a struct that contains CountryCode and LocationCode.
// locodedb.Record is a struct that contains CountryName, LocationName, SubDivName, SubDivCode, GeoPoint and Continent.
func (db *DB) Get(key locodedb.Key) (rec *locodedb.Record, err error) {
	rec = &locodedb.Record{}

	if err := db.getLocodeFromCSV(key, *rec); err != nil {
		return nil, err
	}
	countryName, err := db.getCountryFromCode(key.CountryCode().String())
	if err != nil {
		return nil, err
	}
	rec.SetCountryName(countryName)

	continentName, err := db.getContinentFromCode(*rec.Continent())
	if err != nil {
		return nil, err
	}
	rec.SetContinent(continentName)

	return
}

func (db *DB) getCountryFromCode(code string) (countryName string, err error) {
	file, err := os.Open(db.pathCSVCountries)
	if err != nil {
		return countryName, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return countryName, err
	}

	for _, record := range records {
		if record[0] == code {
			return record[1], nil
		}
	}
	return countryName, errRecordNotFound
}

func (db *DB) getContinentFromCode(continent locodedb.Continent) (c *locodedb.Continent, err error) {
	file, err := os.Open(db.pathCSVContinents)
	if err != nil {
		return c, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return c, err
	}
	for _, record := range records {
		if record[0] == string(continent) {
			cont := locodedb.ContinentFromString(record[1])
			return &cont, nil
		}
	}
	return
}

func (db *DB) getLocodeFromCSV(key locodedb.Key, rec locodedb.Record) (err error) {
	file, err := os.Open(db.pathCSVLocode)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for _, record := range records {
		if record[0] == key.CountryCode().String() && record[1] == key.LocationCode().String() {
			rec.SetLocationName(record[2])

			cont := locodedb.ContinentFromString(record[3])
			rec.SetContinent(&cont)

			rec.SetSubDivName(record[4])
			rec.SetSubDivCode(record[5])

			lat, _ := strconv.ParseFloat(record[6], 64)
			lon, _ := strconv.ParseFloat(record[7], 64)
			rec.SetGeoPoint(locodedb.NewPoint(lat, lon))

		}
	}
	return

}

// Put puts a record by key. locodedb.Key is a struct that contains CountryCode and LocationCode.
// locodedb.Record is a struct that contains CountryName, LocationName, SubDivName, SubDivCode, GeoPoint and Continent.
func (db *DB) Put(key locodedb.Key, rec locodedb.Record) error {
	if err := db.putLocodeToCSV(key, rec); err != nil {
		return err
	}

	if err := db.putCountryToCSV(key, rec); err != nil {
		return err
	}

	if err := db.putContinentToCSV(rec); err != nil {
		return err
	}
	return nil
}

func (db *DB) putLocodeToCSV(key locodedb.Key, rec locodedb.Record) error {
	file, err := os.OpenFile(db.pathCSVLocode, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		key.CountryCode().String(),
		key.LocationCode().String(),
		rec.LocationName(),
		strconv.Itoa(int(*rec.Continent())),
		rec.SubDivCode(),
		rec.SubDivName(),
		strconv.FormatFloat(rec.GeoPoint().Latitude(), 'f', -1, 64),
		strconv.FormatFloat(rec.GeoPoint().Longitude(), 'f', -1, 64),
	}

	if err := writer.Write(record); err != nil {
		return err
	}

	return nil
}

func (db *DB) putCountryToCSV(key locodedb.Key, rec locodedb.Record) error {
	file, err := os.OpenFile(db.pathCSVCountries, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		key.CountryCode().String(),
		rec.CountryName(),
	}

	if err := writer.Write(record); err != nil {
		return err
	}

	return nil
}
func (db *DB) putContinentToCSV(rec locodedb.Record) error {
	file, err := os.OpenFile(db.pathCSVContinents, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		strconv.Itoa(int(*rec.Continent())),
		rec.Continent().String(),
	}

	if err := writer.Write(record); err != nil {
		return err
	}

	return nil
}

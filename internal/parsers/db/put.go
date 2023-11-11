package locodedb

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/nspcc-dev/locode-db/pkg/locodedb"
)

const (
	filenameCSVLocode    = "locodes.csv"
	filenameCSVCountries = "countries.csv"
)

// Data is a struct that contains the Key and the Record.
type Data struct {
	Key    locodedb.Key
	Record locodedb.Record
}

// Put writes the []Data to the CSV files.
func (db *CsvDB) Put(data []Data) error {
	newRecordsLocode := make([][]string, 0, len(data))
	newRecordsCountry := make([][]string, 0, 300)

	uniqueKeys := make(map[string]struct{}, len(data))
	uniqueKeysCountry := make(map[string]struct{}, 300)

	for _, row := range data {
		key := row.Key
		rec := row.Record

		// Calculate a unique index for each key
		keyString := key.CountryCode().String() + key.LocationCode().String()

		if _, exists := uniqueKeys[keyString]; exists {
			continue // Skip duplicates
		}

		// Mark the index as seen
		uniqueKeys[keyString] = struct{}{}

		newRecord := []string{
			key.CountryCode().String(),
			key.LocationCode().String(),
			rec.Location,
			strconv.Itoa(int(rec.Cont)),
			rec.SubDivCode,
			rec.SubDivName,
			strconv.FormatFloat(rec.Point.Latitude(), 'f', -1, 64),
			strconv.FormatFloat(rec.Point.Longitude(), 'f', -1, 64),
		}

		newRecordsLocode = append(newRecordsLocode, newRecord)

		if _, exists := uniqueKeysCountry[key.CountryCode().String()]; exists {
			continue
		}

		uniqueKeysCountry[key.CountryCode().String()] = struct{}{}

		newRecordCountry := []string{
			key.CountryCode().String(),
			rec.Country,
		}

		newRecordsCountry = append(newRecordsCountry, newRecordCountry)
	}

	sort.Slice(newRecordsLocode, func(i, j int) bool {
		return newRecordsLocode[i][0]+newRecordsLocode[i][1] < newRecordsLocode[j][0]+newRecordsLocode[j][1]
	})

	sort.Slice(newRecordsCountry, func(i, j int) bool {
		return newRecordsCountry[i][0] < newRecordsCountry[j][0]
	})

	err := writeToCsvFile(newRecordsCountry, filepath.Join(db.path, filenameCSVCountries))
	if err != nil {
		return err
	}
	return writeToCsvFile(newRecordsLocode, filepath.Join(db.path, filenameCSVLocode))
}

func writeToCsvFile(newRecords [][]string, path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	for _, record := range newRecords {
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()

	return nil
}

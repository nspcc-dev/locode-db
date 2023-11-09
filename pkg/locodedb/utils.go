package locodedb

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/csv"
	"io"
	"strconv"
	"sync"
)

var (
	//go:embed data/countries.csv.gz
	countriesData []byte

	//go:embed data/locodes.csv.gz
	locodesData []byte

	locodeDataOnce sync.Once
)

func initLocodeData() (err error) {
	locodeDataOnce.Do(func() {
		mCountries, err = unpackCountriesData(countriesData)
		if err != nil {
			return
		}
		countriesData = nil
		mLocodes, err = unpackLocodesData(locodesData)
		locodesData = nil
	})
	return
}

type locodesCSV struct {
	locationName string
	subDivCode   string
	subDivName   string
	point        *Point
	continent    *Continent
}

func unpackCountriesData(data []byte) (map[CountryCode]string, error) {
	m := make(map[CountryCode]string)

	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return m, err
	}
	defer gzReader.Close()

	reader := csv.NewReader(gzReader)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return m, err
		}
		countryCode, err := CountryCodeFromString(record[0])
		if err != nil {
			return m, err
		}
		m[*countryCode] = record[1]
	}
	return m, nil
}

func unpackLocodesData(data []byte) (map[string]locodesCSV, error) {
	m := make(map[string]locodesCSV)
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return m, err
	}
	defer gzReader.Close()

	reader := csv.NewReader(gzReader)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return m, err
		}

		countryCode, err := CountryCodeFromString(record[0])
		if err != nil {
			return m, err
		}
		locationCode, err := LocationCodeFromString(record[1])
		if err != nil {
			return m, err
		}

		cont, _ := strconv.ParseUint(record[3], 10, 8)
		var continent = Continent(uint8(cont))

		subDivCode := record[4]
		subDivName := record[5]

		lat, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			return m, err
		}
		lon, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			return m, err
		}

		geoPoint := NewPoint(lat, lon)

		m[countryCode.String()+locationCode.String()] = locodesCSV{
			locationName: record[2],
			continent:    &continent,
			subDivCode:   subDivCode,
			subDivName:   subDivName,
			point:        geoPoint,
		}
	}
	return m, nil
}

package locodedb

import (
	"bytes"
	"compress/bzip2"
	_ "embed"
	"encoding/csv"
	"errors"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"
)

var (
	//go:embed data/countries.csv.bz2
	countriesData []byte

	//go:embed data/locodes.csv.bz2
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
		locodeStrings, mLocodes, err = unpackLocodesData(locodesData)
		locodesData = nil
	})
	return
}

type offLen struct {
	offset uint32
	length uint8 // Likely to be aligned to 32, unfortunately.
}

type locodesCSV struct {
	point        Point
	locationName offLen
	subDivCode   offLen
	subDivName   offLen
	continent    Continent
}

func unpackCountriesData(data []byte) (map[CountryCode]string, error) {
	m := make(map[CountryCode]string)

	zReader := bzip2.NewReader(bytes.NewReader(data))
	reader := csv.NewReader(zReader)

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

func unpackLocodesData(data []byte) (string, map[string]locodesCSV, error) {
	var (
		b       strings.Builder
		m       = make(map[string]locodesCSV)
		zReader = bzip2.NewReader(bytes.NewReader(data))
		reader  = csv.NewReader(zReader)
	)
	reader.ReuseRecord = true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", nil, err
		}

		countryCode, err := CountryCodeFromString(record[0])
		if err != nil {
			return "", nil, err
		}
		locationCode, err := LocationCodeFromString(record[1])
		if err != nil {
			return "", nil, err
		}

		if len(record[2]) > math.MaxUint8 || len(record[4]) > math.MaxUint8 || len(record[5]) > math.MaxUint8 {
			return "", nil, errors.New("record string uint8 overflow")
		}

		var location = offLen{offset: uint32(b.Len()), length: uint8(len(record[2]))}
		b.WriteString(record[2])

		var subDivCode = offLen{offset: uint32(b.Len()), length: uint8(len(record[4]))}
		b.WriteString(record[4])

		var subDivName = offLen{offset: uint32(b.Len()), length: uint8(len(record[5]))}
		b.WriteString(record[5])

		if b.Len() > math.MaxInt32 {
			return "", nil, errors.New("string buffer int32 overflow")
		}
		cont, _ := strconv.ParseUint(record[3], 10, 8)
		var continent = Continent(uint8(cont))

		lat, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			return "", nil, err
		}
		lon, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			return "", nil, err
		}

		m[countryCode.String()+locationCode.String()] = locodesCSV{
			locationName: location,
			continent:    continent,
			subDivCode:   subDivCode,
			subDivName:   subDivName,
			point:        Point{lat: lat, lng: lon},
		}
	}
	return b.String(), m, nil
}

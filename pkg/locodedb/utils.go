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
		locodeStrings, err = unpackLocodesData(locodesData, mCountries)
		locodesData = nil
	})
	return
}

type countryData struct {
	name    string
	locodes []locodesCSV
}

type locodesCSV struct {
	point         Point
	offset        uint32
	code          LocationCode
	locationLen   uint8
	subDivCodeLen uint8
	subDivNameLen uint8
	continent     Continent
}

func unpackCountriesData(data []byte) (map[CountryCode]countryData, error) {
	m := make(map[CountryCode]countryData)

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
		m[*countryCode] = countryData{name: record[1]}
	}
	return m, nil
}

func unpackLocodesData(data []byte, mc map[CountryCode]countryData) (string, error) {
	var (
		b       strings.Builder
		zReader = bzip2.NewReader(bytes.NewReader(data))
		reader  = csv.NewReader(zReader)
	)
	reader.ReuseRecord = true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		if len(record[0]) != CountryCodeLen+LocationCodeLen {
			return "", errors.New("bad locode record length")
		}
		if len(record[1]) > math.MaxUint8 || len(record[3]) > math.MaxUint8 || len(record[4]) > math.MaxUint8 {
			return "", errors.New("record string uint8 overflow")
		}
		if b.Len() > math.MaxInt32 {
			return "", errors.New("string buffer int32 overflow")
		}
		var (
			code          LocationCode
			recOffset     = uint32(b.Len())
			locationLen   = uint8(len(record[1]))
			subDivCodeLen = uint8(len(record[3]))
			subDivNameLen = uint8(len(record[4]))
		)

		copy(code[:], record[0][CountryCodeLen:])
		b.WriteString(record[1])
		b.WriteString(record[3])
		b.WriteString(record[4])

		cont, _ := strconv.ParseUint(record[2], 10, 8)
		var continent = Continent(uint8(cont))

		lat, err := strconv.ParseFloat(record[5], 32)
		if err != nil {
			return "", err
		}
		lon, err := strconv.ParseFloat(record[6], 32)
		if err != nil {
			return "", err
		}

		countryCode, err := CountryCodeFromString(record[0][:CountryCodeLen])
		if err != nil {
			return "", err
		}
		rec, ok := mc[*countryCode]
		if !ok {
			return "", errors.New("invalid country in the DB")
		}
		rec.locodes = append(rec.locodes, locodesCSV{
			point:         Point{Latitude: float32(lat), Longitude: float32(lon)},
			code:          code,
			offset:        recOffset,
			locationLen:   locationLen,
			subDivCodeLen: subDivCodeLen,
			subDivNameLen: subDivNameLen,
			continent:     continent,
		})
		mc[*countryCode] = rec
	}
	for k := range mc {
		rec := mc[k]
		rec.locodes = rec.locodes[:len(rec.locodes):len(rec.locodes)]
		mc[k] = rec
	}
	return b.String(), nil
}

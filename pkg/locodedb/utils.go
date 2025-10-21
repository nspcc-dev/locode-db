package locodedb

import (
	"bytes"
	"compress/bzip2"
	_ "embed"
	"encoding/csv"
	"errors"
	"io"
	"math"
	"slices"
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
	locationLen   uint8
	subDivCodeLen uint8
	subDivNameLen uint8
	continent     Continent
}

func unpackCountriesData(data []byte) (map[countryCode]countryData, error) {
	m := make(map[countryCode]countryData)

	zReader := bzip2.NewReader(bytes.NewReader(data))
	reader := csv.NewReader(zReader)

	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return m, err
		}
		cc, err := countryCodeFromString(record[0])
		if err != nil {
			return m, err
		}
		m[*cc] = countryData{name: record[1]}
	}
	return m, nil
}

func unpackLocodesData(data []byte, mc map[countryCode]countryData) (string, error) {
	var (
		b       strings.Builder
		zReader = bzip2.NewReader(bytes.NewReader(data))
		reader  = csv.NewReader(zReader)
	)
	reader.ReuseRecord = true

	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
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
			recOffset     = uint32(b.Len())
			locationLen   = uint8(len(record[1]))
			subDivCodeLen = uint8(len(record[3]))
			subDivNameLen = uint8(len(record[4]))
		)

		b.WriteString(record[0][CountryCodeLen:])
		b.WriteString(record[1])
		b.WriteString(record[3])
		b.WriteString(record[4])

		cont, _ := strconv.ParseUint(record[2], 10, 8)
		var continent = Continent(uint8(cont))

		lat, err := strconv.ParseFloat(record[5], 32)
		if err != nil {
			return "", err
		}
		lng, err := strconv.ParseFloat(record[6], 32)
		if err != nil {
			return "", err
		}

		cc, err := countryCodeFromString(record[0][:CountryCodeLen])
		if err != nil {
			return "", err
		}
		rec, ok := mc[*cc]
		if !ok {
			return "", errors.New("invalid country in the DB")
		}
		rec.locodes = append(rec.locodes, locodesCSV{
			point:         Point{Latitude: float32(lat), Longitude: float32(lng)},
			offset:        recOffset,
			locationLen:   locationLen,
			subDivCodeLen: subDivCodeLen,
			subDivNameLen: subDivNameLen,
			continent:     continent,
		})
		mc[*cc] = rec
	}
	for k := range mc {
		rec := mc[k]
		rec.locodes = slices.Clip(rec.locodes)
		mc[k] = rec
	}
	return b.String(), nil
}

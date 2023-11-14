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

type locodesCSV struct {
	point         Point
	offset        uint32
	locationLen   uint8
	subDivCodeLen uint8
	subDivNameLen uint8
	continent     Continent
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

func unpackLocodesData(data []byte) (string, []locodesCSV, error) {
	var (
		b       strings.Builder
		m       []locodesCSV
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

		if len(record[0]) != CountryCodeLen+LocationCodeLen {
			return "", nil, errors.New("bad locode record length")
		}
		if len(record[1]) > math.MaxUint8 || len(record[3]) > math.MaxUint8 || len(record[4]) > math.MaxUint8 {
			return "", nil, errors.New("record string uint8 overflow")
		}
		if b.Len() > math.MaxInt32 {
			return "", nil, errors.New("string buffer int32 overflow")
		}
		var (
			recOffset     = uint32(b.Len())
			locationLen   = uint8(len(record[1]))
			subDivCodeLen = uint8(len(record[3]))
			subDivNameLen = uint8(len(record[4]))
		)

		b.WriteString(record[0])
		b.WriteString(record[1])
		b.WriteString(record[3])
		b.WriteString(record[4])

		cont, _ := strconv.ParseUint(record[2], 10, 8)
		var continent = Continent(uint8(cont))

		lat, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return "", nil, err
		}
		lon, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			return "", nil, err
		}

		m = append(m, locodesCSV{
			point:         Point{lat: lat, lng: lon},
			offset:        recOffset,
			locationLen:   locationLen,
			subDivCodeLen: subDivCodeLen,
			subDivNameLen: subDivNameLen,
			continent:     continent,
		})
	}
	return b.String(), m, nil
}

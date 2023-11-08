package csvlocode

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	locode "github.com/nspcc-dev/locode-db/internal/parsers/db"
	"github.com/nspcc-dev/locode-db/pkg/locodedb"
	"golang.org/x/text/encoding/charmap"
)

var errInvalidSubName = errors.New("could not convert subDivName to uft-8 valid string")
var errInvalidRecord = errors.New("invalid table record")

// IterateAll scans a table record one-by-one, parses a UN/LOCODE record
// from it and passes it to f.
//
// Returns f's errors directly.
func (t *Table) IterateAll(f func(locode.Record) error) error {
	const wordsPerRecord = 12

	return t.scanWords(t.paths, wordsPerRecord, func(words []string) error {
		lc := [2]string{words[1], words[2]}

		record := locode.Record{
			LOCODE:           lc,
			Name:             words[3],
			NameWoDiacritics: words[4],
			SubDiv:           words[5],
			Function:         words[6],
			Status:           words[7],
			Date:             words[8],
			IATA:             words[9],
			Coordinates:      words[10],
			Remarks:          words[11],
		}

		return f(record)
	})
}

const (
	_ = iota - 1

	subDivCountry
	subDivSubdivision
	subDivName
	_ // subDivLevel

	subDivFldNum
)

type subDivKey struct {
	countryCode,
	subDivCode string
}

type subDivRecord struct {
	name string
}

// SubDivName scans a table record to an in-memory table (once),
// and returns the subdivision name of the country and the subdivision codes match.
//
// Returns locodedb.ErrSubDivNotFound if no entry matches.
func (t *Table) SubDivName(countryCode *locodedb.CountryCode, code string) (string, error) {
	if err := t.initSubDiv(); err != nil {
		return "", err
	}

	rec, ok := t.mSubDiv[subDivKey{
		countryCode: countryCode.String(),
		subDivCode:  code,
	}]
	if !ok {
		return "", locode.ErrSubDivNotFound
	}

	return rec.name, nil
}

func isValidString(s string) bool {
	return !strings.Contains(s, "\uFFFD") && !strings.Contains(s, "\u0000") && !strings.Contains(s, "?")
}

func utf8encoding(subdivCountry string) (string, error) {
	res, err := charmap.Windows1256.NewDecoder().String(subdivCountry)
	if err == nil && isValidString(res) {
		return res, nil
	}

	res, err = charmap.ISO8859_1.NewDecoder().String(subdivCountry)
	if err == nil && isValidString(res) {
		return res, nil
	}

	return subdivCountry, errInvalidSubName
}

func (t *Table) initSubDiv() (err error) {
	t.subDivOnce.Do(func() {
		t.mSubDiv = make(map[subDivKey]subDivRecord)

		err = t.scanWords([]string{t.subDivPath}, subDivFldNum, func(words []string) error {
			subdiv := words[subDivName]
			if !utf8.ValidString(subdiv) {
				subdiv, err = utf8encoding(subdiv)
				if err != nil {
					if errors.Is(err, errInvalidSubName) {
						return nil
					}
					return err
				}
			}
			t.mSubDiv[subDivKey{
				countryCode: words[subDivCountry],
				subDivCode:  words[subDivSubdivision],
			}] = subDivRecord{
				name: subdiv,
			}

			return nil
		})
	})

	return
}

var errScanInt = errors.New("interrupt scan")

func (t *Table) scanWords(paths []string, fpr int, wordsHandler func([]string) error) error {
	var (
		rdrs    = make([]io.Reader, 0, len(t.paths))
		closers = make([]io.Closer, 0, len(t.paths))
	)

	for i := range paths {
		file, err := os.OpenFile(paths[i], os.O_RDONLY, t.mode)
		if err != nil {
			return err
		}

		rdrs = append(rdrs, file)
		closers = append(closers, file)
	}

	defer func() {
		for i := range closers {
			_ = closers[i].Close()
		}
	}()

	r := csv.NewReader(io.MultiReader(rdrs...))
	r.ReuseRecord = true
	r.FieldsPerRecord = fpr

	for {
		words, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		} else if len(words) != fpr {
			return errInvalidRecord
		}

		if err := wordsHandler(words); err != nil {
			if errors.Is(err, errScanInt) {
				break
			}

			return err
		}
	}

	return nil
}

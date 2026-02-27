package main

import (
	"errors"
	"flag"
	"log"
	"slices"
	"strings"

	locode "github.com/nspcc-dev/locode-db/internal/parsers/db"
	airportsdb "github.com/nspcc-dev/locode-db/internal/parsers/db/airports"
	continentsdb "github.com/nspcc-dev/locode-db/internal/parsers/db/continents/geojson"
	csvlocode "github.com/nspcc-dev/locode-db/internal/parsers/table/csv"
)

type namesDB struct {
	*airportsdb.DB
	*csvlocode.Table
}

const (
	locodeGenerateInputFlag      = "in"
	locodeGenerateSubDivFlag     = "subdiv"
	locodeGenerateAirportsFlag   = "airports"
	locodeGenerateCountriesFlag  = "countries"
	locodeGenerateContinentsFlag = "continents"
	locodeGenerateOutputFlag     = "out"
)

var (
	locodeGenerateInPaths        []string
	locodeGenerateSubDivPath     string
	locodeGenerateAirportsPath   string
	locodeGenerateCountriesPath  string
	locodeGenerateContinentsPath string
	locodeGenerateOutPath        string
)

func init() {
	flag.Func(locodeGenerateInputFlag, "List of paths to UN/LOCODE tables (CSV)", func(s string) error {
		splitStr := strings.Split(s, ",")
		if slices.Contains(splitStr, "") {
			return errors.New("path is empty")
		}
		locodeGenerateInPaths = append(locodeGenerateInPaths, splitStr...)
		return nil
	})
	flag.StringVar(&locodeGenerateSubDivPath, locodeGenerateSubDivFlag, "", "Path to UN/LOCODE subdivision database (CSV)")
	flag.StringVar(&locodeGenerateAirportsPath, locodeGenerateAirportsFlag, "", "Path to OpenFlights airport database (CSV)")
	flag.StringVar(&locodeGenerateCountriesPath, locodeGenerateCountriesFlag, "", "Path to OpenFlights country database (CSV)")
	flag.StringVar(&locodeGenerateContinentsPath, locodeGenerateContinentsFlag, "", "Path to continent polygons (GeoJSON)")
	flag.StringVar(&locodeGenerateOutPath, locodeGenerateOutputFlag, "", "Target path for generated database (directory))")
}

func main() {
	flag.Parse()

	if err := validateFlags(); err != nil {
		log.Fatal(err)
	}

	locodeDB := csvlocode.New(
		csvlocode.Prm{
			Path:       locodeGenerateInPaths[0],
			SubDivPath: locodeGenerateSubDivPath,
		},
		csvlocode.WithExtraPaths(locodeGenerateInPaths[1:]...),
	)

	airportDB := airportsdb.New(airportsdb.Prm{
		AirportsPath:  locodeGenerateAirportsPath,
		CountriesPath: locodeGenerateCountriesPath,
	})

	continentsDB := continentsdb.New(continentsdb.Prm{
		Path: locodeGenerateContinentsPath,
	})

	targetDB := locode.New(locodeGenerateOutPath)

	names := &namesDB{
		DB:    airportDB,
		Table: locodeDB,
	}

	err := locode.FillDatabase(locodeDB, airportDB, continentsDB, names, targetDB)
	if err != nil {
		log.Fatal(err)
	}
}

func validateFlags() error {
	switch {
	case len(locodeGenerateInPaths) == 0:
		return errors.New("at least one UN/LOCODE table is required")
	case locodeGenerateSubDivPath == "":
		return errors.New("path to UN/LOCODE subdivision database is required")
	case locodeGenerateAirportsPath == "":
		return errors.New("path to OpenFlights airport database is required")
	case locodeGenerateCountriesPath == "":
		return errors.New("path to OpenFlights country database is required")
	case locodeGenerateContinentsPath == "":
		return errors.New("path to continent polygons is required")
	case locodeGenerateOutPath == "":
		return errors.New("target path for generated database is required")
	}
	return nil
}

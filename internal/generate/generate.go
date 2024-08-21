package generate

import (
	"flag"
	"fmt"
	"log"

	locode "github.com/nspcc-dev/locode-db/internal/parsers/db"
	airportsdb "github.com/nspcc-dev/locode-db/internal/parsers/db/airports"
	continentsdb "github.com/nspcc-dev/locode-db/internal/parsers/db/continents/geojson"
	csvlocode "github.com/nspcc-dev/locode-db/internal/parsers/table/csv"
)

type namesDB struct {
	*airportsdb.DB
	*csvlocode.Table
}

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return fmt.Sprint(*s)
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
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
	locodeGenerateInPaths        stringSliceFlag
	locodeGenerateSubDivPath     string
	locodeGenerateAirportsPath   string
	locodeGenerateCountriesPath  string
	locodeGenerateContinentsPath string
	locodeGenerateOutPath        string
)

func init() {
	flag.Var(&locodeGenerateInPaths, locodeGenerateInputFlag, "List of paths to UN/LOCODE tables (CSV)")
	flag.StringVar(&locodeGenerateSubDivPath, locodeGenerateSubDivFlag, "", "Path to UN/LOCODE subdivision database (CSV)")
	flag.StringVar(&locodeGenerateAirportsPath, locodeGenerateAirportsFlag, "", "Path to OpenFlights airport database (CSV)")
	flag.StringVar(&locodeGenerateCountriesPath, locodeGenerateCountriesFlag, "", "Path to OpenFlights country database (CSV)")
	flag.StringVar(&locodeGenerateContinentsPath, locodeGenerateContinentsFlag, "", "Path to continent polygons (GeoJSON)")
	flag.StringVar(&locodeGenerateOutPath, locodeGenerateOutputFlag, "", "Target path for generated database (directory))")
}

func Execute() {
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
	for _, inPath := range locodeGenerateInPaths {
		if inPath == "" {
			return fmt.Errorf("no UN/LOCODE table specified")
		}
	}
	switch {
	case len(locodeGenerateInPaths) == 0:
		return fmt.Errorf("at least one UN/LOCODE table is required")
	case locodeGenerateSubDivPath == "":
		return fmt.Errorf("path to UN/LOCODE subdivision database is required")
	case locodeGenerateAirportsPath == "":
		return fmt.Errorf("path to OpenFlights airport database is required")
	case locodeGenerateCountriesPath == "":
		return fmt.Errorf("path to OpenFlights country database is required")
	case locodeGenerateContinentsPath == "":
		return fmt.Errorf("path to continent polygons is required")
	case locodeGenerateOutPath == "":
		return fmt.Errorf("target path for generated database is required")
	}
	return nil
}

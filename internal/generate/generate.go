package generate

import (
	locode "github.com/nspcc-dev/locode-db/internal/parsers/db"
	airportsdb "github.com/nspcc-dev/locode-db/internal/parsers/db/airports"
	continentsdb "github.com/nspcc-dev/locode-db/internal/parsers/db/continents/geojson"
	csvlocode "github.com/nspcc-dev/locode-db/internal/parsers/table/csv"
	"github.com/spf13/cobra"
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

	locodeGenerateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate UN/LOCODE database",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
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
				cmd.PrintErrln(err)
			}
		},
	}
)

func initUtilLocodeGenerateCmd() {
	flags := locodeGenerateCmd.Flags()

	flags.StringSliceVar(&locodeGenerateInPaths, locodeGenerateInputFlag, nil, "List of paths to UN/LOCODE tables (CSV)")
	_ = locodeGenerateCmd.MarkFlagRequired(locodeGenerateInputFlag)

	flags.StringVar(&locodeGenerateSubDivPath, locodeGenerateSubDivFlag, "", "Path to UN/LOCODE subdivision database (CSV)")
	_ = locodeGenerateCmd.MarkFlagRequired(locodeGenerateSubDivFlag)

	flags.StringVar(&locodeGenerateAirportsPath, locodeGenerateAirportsFlag, "", "Path to OpenFlights airport database (CSV)")
	_ = locodeGenerateCmd.MarkFlagRequired(locodeGenerateAirportsFlag)

	flags.StringVar(&locodeGenerateCountriesPath, locodeGenerateCountriesFlag, "", "Path to OpenFlights country database (CSV)")
	_ = locodeGenerateCmd.MarkFlagRequired(locodeGenerateCountriesFlag)

	flags.StringVar(&locodeGenerateContinentsPath, locodeGenerateContinentsFlag, "", "Path to continent polygons (GeoJSON)")
	_ = locodeGenerateCmd.MarkFlagRequired(locodeGenerateContinentsFlag)

	flags.StringVar(&locodeGenerateOutPath, locodeGenerateOutputFlag, "", "Target path for generated database (directory))")
	_ = locodeGenerateCmd.MarkFlagRequired(locodeGenerateOutputFlag)
}

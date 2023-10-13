package csvlocodedb

import (
	"io/fs"
	"os"
)

// Option sets an optional parameter of DB.
type Option func(*options)

type options struct {
	mode              fs.FileMode
	PathCSVLocode     string
	PathCSVCountries  string
	PathCSVContinents string
}

func defaultOpts() *options {
	return &options{
		mode:              os.ModePerm, // 0777
		PathCSVLocode:     "locodes.csv",
		PathCSVCountries:  "countries.csv",
		PathCSVContinents: "continents.csv",
	}
}

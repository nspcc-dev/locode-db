package csvlocodedb

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Prm groups the required parameters of the DB's constructor.
//
// All values must comply with the requirements imposed on them.
// Passing incorrect parameter values will result in constructor
// failure (error or panic depending on the implementation).
type Prm struct {
	// Path to result directory with location database.
	//
	// Must not be empty.
	Path              string
	PathCSVLocode     string
	PathCSVCountries  string
	PathCSVContinents string
}

// DB is a descriptor of the location database.
//
// For correct operation, DB must be created
// using the constructor (New) based on the required parameters
// and optional components.
type DB struct {
	mode              fs.FileMode
	path              string
	pathCSVLocode     string
	pathCSVCountries  string
	pathCSVContinents string
}

const invalidPrmValFmt = "invalid parameter %s (%T):%v"

func panicOnPrmValue(n string, v any) {
	panic(fmt.Sprintf(invalidPrmValFmt, n, v, v))
}

// New creates a new instance of the DB.
//
// Panics if at least one value of the parameters is invalid.
//
// The created DB requires calling the Open method in order
// to initialize required resources.
func New(prm Prm, opts ...Option) *DB {
	switch {
	case prm.Path == "":
		panicOnPrmValue("Path", prm.Path)
	}
	fileInfo, err := os.Stat(prm.Path)
	if err != nil {
		panicOnPrmValue("Error checking path: ", err.Error())
	}

	if !fileInfo.IsDir() {
		panicOnPrmValue("path is not a directory: ", prm.Path)
	}
	o := defaultOpts()

	for i := range opts {
		opts[i](o)
	}

	return &DB{
		mode:              o.mode,
		path:              prm.Path,
		pathCSVLocode:     filepath.Join(prm.Path, o.PathCSVLocode),
		pathCSVCountries:  filepath.Join(prm.Path, o.PathCSVCountries),
		pathCSVContinents: filepath.Join(prm.Path, o.PathCSVContinents),
	}
}

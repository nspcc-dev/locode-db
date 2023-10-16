package locodedb

import (
	"fmt"
	"os"
)

func panicOnPrmValue(n string, v any) {
	panic(fmt.Sprintf("invalid parameter %s (%T):%v", n, v, v))
}

// CsvDB is a resulting database in CSV format. Path should be a valid path to the directory.
type CsvDB struct {
	path string
}

func New(path string) CsvDB {
	switch {
	case path == "":
		panicOnPrmValue("Path", path)
	}
	_, err := os.Stat(path)

	if err != nil {
		panicOnPrmValue("output directory path", path)
	}

	return CsvDB{
		path: path,
	}
}

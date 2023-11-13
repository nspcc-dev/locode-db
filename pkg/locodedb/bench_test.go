package locodedb

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
)

// Unfortunately, benchmarks are run after "normal" tests and this means that
// original vars are already destroyed by the time test starts, so we have to
// duplicate.
var (
	//go:embed data/countries.csv.bz2
	testCountriesData []byte

	//go:embed data/locodes.csv.bz2
	testLocodesData []byte
)

func BenchmarkUnpack(b *testing.B) {
	require.NotEmpty(b, testCountriesData)
	require.NotEmpty(b, testLocodesData)
	for i := 0; i < b.N; i++ {
		_, err := unpackCountriesData(testCountriesData)
		require.NoError(b, err)
		_, _, err = unpackLocodesData(testLocodesData)
		require.NoError(b, err)
	}
}

func BenchmarkGet(b *testing.B) {
	require.NoError(b, initLocodeData())
	_, err := Get("RU MOW")
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Get("RU MOW")
		_, _ = Get("AAAAA")
		_, _ = Get("SESTO")
		_, _ = Get("FRXGS")
		_, _ = Get("JOSAH")
	}
}

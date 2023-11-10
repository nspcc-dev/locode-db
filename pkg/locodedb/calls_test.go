package locodedb_test

import (
	"testing"

	"github.com/nspcc-dev/locode-db/pkg/locodedb"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Run("wrong locode", func(t *testing.T) {
		_, err := locodedb.Get("WRONG LOCODE")
		require.Error(t, err)
	})
	t.Run("nonexistent locode", func(t *testing.T) {
		_, err := locodedb.Get("AAAAA")
		require.Error(t, err)
	})

	t.Run("locode", func(t *testing.T) {
		rec, err := locodedb.Get("RU MOW")
		require.NoError(t, err)

		require.Equal(t, rec.Country, "Russia")
		require.Equal(t, rec.Location, "Moskva")
		require.Equal(t, rec.SubDivCode, "MOW")
		require.Equal(t, rec.SubDivName, "Moskva")
		require.Equal(t, rec.Cont.String(), "Europe")
	})
	t.Run("locode", func(t *testing.T) {
		rec, err := locodedb.Get("RUMOW")
		require.NoError(t, err)

		require.Equal(t, rec.Country, "Russia")
		require.Equal(t, rec.Location, "Moskva")
		require.Equal(t, rec.SubDivCode, "MOW")
		require.Equal(t, rec.SubDivName, "Moskva")
		require.Equal(t, rec.Cont.String(), "Europe")
	})

	t.Run("valid key", func(t *testing.T) {
		_, err := locodedb.NewKey("RU", "MOW")
		require.NoError(t, err)

	})
}

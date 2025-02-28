package locodedb_test

import (
	"testing"
	"unicode/utf8"

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

	t.Run("utf-8 subdiv locode", func(t *testing.T) {
		rec, err := locodedb.Get("SESTO")
		require.NoError(t, err)
		require.Equal(t, utf8.ValidString(rec.SubDivName), true)

		t.Run("utf-8 locode", func(t *testing.T) {
			rec, err := locodedb.Get("ARADS")
			require.NoError(t, err)
			require.Equal(t, utf8.ValidString(rec.SubDivName), true)
		})
		t.Run("utf-8 locode", func(t *testing.T) {
			rec, err := locodedb.Get("JOSAH")
			require.NoError(t, err)
			require.Equal(t, utf8.ValidString(rec.SubDivName), true)
		})
		t.Run("utf-8 locode", func(t *testing.T) {
			rec, err := locodedb.Get("FRXGS")
			require.NoError(t, err)
			require.Equal(t, utf8.ValidString(rec.SubDivName), true)
		})
		t.Run("utf-8 locode", func(t *testing.T) {
			rec, err := locodedb.Get("BRSPY")
			require.NoError(t, err)
			require.Equal(t, utf8.ValidString(rec.SubDivName), true)
		})
	})
}

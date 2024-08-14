package locodedb

import (
	"fmt"

	"github.com/nspcc-dev/locode-db/pkg/locodedb"
)

// PointFromCoordinates converts a UN/LOCODE coordinates to a Point.
func PointFromCoordinates(crd *Coordinates) (locodedb.Point, error) {
	if crd == nil {
		return locodedb.Point{}, nil
	}

	lat, err := crd.Latitude().ToDecimalDegrees()
	if err != nil {
		return locodedb.Point{}, fmt.Errorf("could not parse latitude: %w", err)
	}

	lng, err := crd.Longitude().ToDecimalDegrees()
	if err != nil {
		return locodedb.Point{}, fmt.Errorf("could not parse longitude: %w", err)
	}

	return locodedb.Point{
		Latitude:  float32(lat),
		Longitude: float32(lng),
	}, nil
}

package locodedb

import (
	"fmt"
	"strconv"

	"github.com/nspcc-dev/locode-db/pkg/locodedb"
)

// PointFromCoordinates converts a UN/LOCODE coordinates to a Point.
func PointFromCoordinates(crd *Coordinates) (locodedb.Point, error) {
	if crd == nil {
		return locodedb.Point{}, nil
	}

	cLat := crd.Latitude()
	cLatDeg := cLat.Degrees()
	cLatMnt := cLat.Minutes()

	lat, err := toDecimal(cLatDeg[:], cLatMnt[:])
	if err != nil {
		return locodedb.Point{}, fmt.Errorf("could not parse latitude: %w", err)
	}

	if !cLat.Hemisphere().North() {
		lat = -lat
	}

	cLng := crd.Longitude()
	cLngDeg := cLng.Degrees()
	cLngMnt := cLng.Minutes()

	lng, err := toDecimal(cLngDeg[:], cLngMnt[:])
	if err != nil {
		return locodedb.Point{}, fmt.Errorf("could not parse longitude: %w", err)
	}

	if !cLng.Hemisphere().East() {
		lng = -lng
	}

	return locodedb.Point{
		Latitude:  float32(lat),
		Longitude: float32(lng),
	}, nil
}

func toDecimal(intRaw, minutesRaw []byte) (float64, error) {
	integer, err := strconv.ParseFloat(string(intRaw), 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse integer part: %w", err)
	}

	decimal, err := minutesToDegrees(minutesRaw)
	if err != nil {
		return 0, fmt.Errorf("could not parse decimal part: %w", err)
	}

	return integer + decimal, nil
}

// minutesToDegrees converts minutes to decimal part of a degree.
func minutesToDegrees(raw []byte) (float64, error) {
	minutes, err := strconv.ParseFloat(string(raw), 64)
	if err != nil {
		return 0, err
	}

	return minutes / 60, nil
}

package locodedb

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nspcc-dev/locode-db/pkg/locodedb"
)

const (
	minutesDigits     = 2
	hemisphereSymbols = 1
)

const (
	latDegDigits = 2
	lngDegDigits = 3
)

type coordinateCode struct {
	degDigits int
	value     []uint8
}

// LongitudeCode represents the value of the longitude
// of the location conforming to UN/LOCODE specification.
type LongitudeCode coordinateCode

// LongitudeHemisphere represents the hemisphere of the earth
// // along the Greenwich meridian.
type LongitudeHemisphere [hemisphereSymbols]uint8

// LatitudeCode represents the value of the latitude
// of the location conforming to UN/LOCODE specification.
type LatitudeCode coordinateCode

// LatitudeHemisphere represents the hemisphere of the earth
// along the equator.
type LatitudeHemisphere [hemisphereSymbols]uint8

func coordinateFromString(s string, degDigits int, hemisphereAlphabet []uint8) (*coordinateCode, error) {
	if strings.Contains(s, ".") {
		return &coordinateCode{
			degDigits: len(s) - minutesDigits - hemisphereSymbols,
			value:     []uint8(s),
		}, nil
	}

	if len(s) != degDigits+minutesDigits+hemisphereSymbols {
		return nil, locodedb.ErrInvalidString
	}

	for i := range s[:degDigits+minutesDigits] {
		if !isDigit(s[i]) {
			return nil, locodedb.ErrInvalidString
		}
	}

loop:
	for _, sym := range s[degDigits+minutesDigits:] {
		for j := range hemisphereAlphabet {
			if hemisphereAlphabet[j] == uint8(sym) {
				continue loop
			}
		}

		return nil, locodedb.ErrInvalidString
	}

	return &coordinateCode{
		degDigits: degDigits,
		value:     []uint8(s),
	}, nil
}

func isDigit(sym uint8) bool {
	return sym >= '0' && sym <= '9'
}

// LongitudeFromString parses a string and returns the location's longitude.
func LongitudeFromString(s string) (*LongitudeCode, error) {
	cc, err := coordinateFromString(s, lngDegDigits, []uint8{'W', 'E'})
	if err != nil {
		return nil, err
	}

	return (*LongitudeCode)(cc), nil
}

// LatitudeFromString parses a string and returns the location's latitude.
func LatitudeFromString(s string) (*LatitudeCode, error) {
	cc, err := coordinateFromString(s, latDegDigits, []uint8{'N', 'S'})
	if err != nil {
		return nil, err
	}

	return (*LatitudeCode)(cc), nil
}

func (cc *coordinateCode) degrees() []uint8 {
	return cc.value[:cc.degDigits]
}

// Degrees returns the longitude's degrees.
func (lc *LongitudeCode) Degrees() (l [lngDegDigits]uint8) {
	copy(l[:], (*coordinateCode)(lc).degrees())
	return
}

// Degrees returns the latitude's degrees.
func (lc *LatitudeCode) Degrees() (l [latDegDigits]uint8) {
	copy(l[:], (*coordinateCode)(lc).degrees())
	return
}

func (cc *coordinateCode) minutes() (mnt [minutesDigits]uint8) {
	for i := 0; i < minutesDigits; i++ {
		mnt[i] = cc.value[cc.degDigits+i]
	}

	return
}

// Minutes returns the longitude's minutes.
func (lc *LongitudeCode) Minutes() [minutesDigits]uint8 {
	return (*coordinateCode)(lc).minutes()
}

// Minutes returns the latitude's minutes.
func (lc *LatitudeCode) Minutes() [minutesDigits]uint8 {
	return (*coordinateCode)(lc).minutes()
}

// Hemisphere returns the longitude's hemisphere code.
func (lc *LongitudeCode) Hemisphere() LongitudeHemisphere {
	return (*coordinateCode)(lc).hemisphere()
}

// Hemisphere returns the latitude's hemisphere code.
func (lc *LatitudeCode) Hemisphere() LatitudeHemisphere {
	return (*coordinateCode)(lc).hemisphere()
}

func (cc *coordinateCode) hemisphere() (h [hemisphereSymbols]uint8) {
	for i := 0; i < hemisphereSymbols; i++ {
		h[i] = cc.value[cc.degDigits+minutesDigits+i]
	}

	return h
}

// North returns true for the northern hemisphere.
func (h LatitudeHemisphere) North() bool {
	return h[0] == 'N'
}

// East returns true for the eastern hemisphere.
func (h LongitudeHemisphere) East() bool {
	return h[0] == 'E'
}

// Coordinates represents the coordinates of the location from UN/LOCODE table.
type Coordinates struct {
	lat *LatitudeCode

	lng *LongitudeCode
}

// Latitude returns the location's latitude.
func (c *Coordinates) Latitude() *LatitudeCode {
	return c.lat
}

// Longitude returns the location's longitude.
func (c *Coordinates) Longitude() *LongitudeCode {
	return c.lng
}

// CoordinatesFromString parses a string and returns the location's coordinates.
func CoordinatesFromString(s string) (*Coordinates, error) {
	if len(s) == 0 {
		return nil, nil
	}

	strs := strings.Split(s, " ")
	if len(strs) != 2 {
		return nil, locodedb.ErrInvalidString
	}

	lat, err := LatitudeFromString(strs[0])
	if err != nil {
		return nil, fmt.Errorf("could not parse latitude: %w", err)
	}

	lng, err := LongitudeFromString(strs[1])
	if err != nil {
		return nil, fmt.Errorf("could not parse longitude: %w", err)
	}

	return &Coordinates{
		lat: lat,
		lng: lng,
	}, nil
}

// ToDecimalDegrees returns decimal representation of the longitude or an error if the conversion fails.
func (lc *LongitudeCode) ToDecimalDegrees() (float64, error) {
	return decimalDegreesFromCoordinateCode((*coordinateCode)(lc), lc.Hemisphere().East())
}

// ToDecimalDegrees returns the decimal representation of the latitude or an error if the conversion fails.
func (lc *LatitudeCode) ToDecimalDegrees() (float64, error) {
	return decimalDegreesFromCoordinateCode((*coordinateCode)(lc), lc.Hemisphere().North())
}

// decimalDegreesFromCoordinateCode returns the decimal representation of the coordinate
// or an error if the conversion fails.
// Takes a coordinate code and a boolean indicating
// if the coordinate is in the positive hemisphere.
func decimalDegreesFromCoordinateCode(crdCode *coordinateCode, positiveHemisphere bool) (float64, error) {
	var (
		value float64
		err   error
	)

	crdString := string(crdCode.value)
	// checking that the coordinates can be in decimal degrees
	if strings.Contains(crdString, ".") {
		value, err = strconv.ParseFloat(crdString[:len(crdString)-hemisphereSymbols], 64)
		if err != nil {
			return 0, fmt.Errorf("could not parse coordinate: %w", err)
		}
	} else {
		crdDeg := crdCode.degrees()
		crdMnt := crdCode.minutes()

		value, err = bytesToDecimal(crdDeg[:], crdMnt[:])
		if err != nil {
			return 0, fmt.Errorf("could not parse coordinate: %w", err)
		}
	}

	if !positiveHemisphere {
		value = -value
	}

	return value, nil
}

// bytesToDecimal converts degree and minute components to decimal degrees.
func bytesToDecimal(intRaw, minutesRaw []byte) (float64, error) {
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

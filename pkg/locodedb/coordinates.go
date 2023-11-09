package locodedb

import (
	"fmt"
	"strings"
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
	if len(s) != degDigits+minutesDigits+hemisphereSymbols {
		return nil, ErrInvalidString
	}

	for i := range s[:degDigits+minutesDigits] {
		if !isDigit(s[i]) {
			return nil, ErrInvalidString
		}
	}

loop:
	for _, sym := range s[degDigits+minutesDigits:] {
		for j := range hemisphereAlphabet {
			if hemisphereAlphabet[j] == uint8(sym) {
				continue loop
			}
		}

		return nil, ErrInvalidString
	}

	return &coordinateCode{
		degDigits: degDigits,
		value:     []uint8(s),
	}, nil
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
		return nil, ErrInvalidString
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

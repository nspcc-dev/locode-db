package locodedb

import (
	"testing"

	"github.com/nspcc-dev/locode-db/pkg/locodedb"
)

func TestPointFromCoordinates(t *testing.T) {
	testCases := []struct {
		name             string
		coordsGot        string
		latWant, lngWant float32
		wantErr          bool
	}{
		{
			name:      "Valid coordinates, all positive",
			coordsGot: "5915N 01806E",
			latWant:   59.25,
			lngWant:   18.10,
			wantErr:   false,
		},
		{
			name:      "Valid coordinates, lat negative, lng positive",
			coordsGot: "1000S 02030E",
			latWant:   -10.00,
			lngWant:   20.50,
			wantErr:   false,
		},
		{
			name:      "Valid coordinates, all negative",
			coordsGot: "0145S 03512W",
			latWant:   -01.75,
			lngWant:   -35.20,
			wantErr:   false,
		},
		{
			name:      "Valid coordinates in decimal degree, all positive",
			coordsGot: "26.8618N 89.3748E",
			latWant:   26.8618,
			lngWant:   89.3748,
			wantErr:   false,
		},
		{
			name:      "Valid coordinates in decimal degree, all negative",
			coordsGot: "26.9374S 89.0233W",
			latWant:   -26.9374,
			lngWant:   -89.0233,
			wantErr:   false,
		},
		{
			name:      "Valid coordinates in decimal degree, all positive",
			coordsGot: "26.8618123N 89.37E",
			latWant:   26.8618123,
			lngWant:   89.37,
			wantErr:   false,
		},
		{
			name:      "No lat and lng, so point is empty",
			coordsGot: "",
			latWant:   0,
			lngWant:   0,
			wantErr:   false,
		},
		{
			name:      "Invalid coordinates, in lat to many digits",
			coordsGot: "01451S 03512W",
			latWant:   0,
			lngWant:   0,
			wantErr:   true,
		},
		{
			name:      "Invalid coordinates, in lng to many digits",
			coordsGot: "0145S 035112W",
			latWant:   0,
			lngWant:   0,
			wantErr:   true,
		},
		{
			name:      "Invalid coordinates, no lat",
			coordsGot: "035112W",
			latWant:   0,
			lngWant:   0,
			wantErr:   true,
		},
		{
			name:      "Invalid coordinates, lat haven't got hemisphere",
			coordsGot: "0145 035112W",
			latWant:   0,
			lngWant:   0,
			wantErr:   true,
		},
		{
			name:      "Invalid coordinates, lng haven't got hemisphere",
			coordsGot: "0145N 035112",
			latWant:   0,
			lngWant:   0,
			wantErr:   true,
		},
		{
			name:      "Invalid coordinates, lat have invalid hemisphere",
			coordsGot: "0145Q 035112W",
			latWant:   0,
			lngWant:   0,
			wantErr:   true,
		},
		{
			name:      "Invalid coordinates, lng have invalid hemisphere",
			coordsGot: "0145N 035112P",
			latWant:   0,
			lngWant:   0,
			wantErr:   true,
		},
		{
			name:      "Invalid coordinates, lat have invalid integer part",
			coordsGot: "0a45S 035112W",
			latWant:   0,
			lngWant:   0,
			wantErr:   true,
		},
		{
			name:      "Invalid coordinates, lng have invalid decimal part",
			coordsGot: "0145N 035t2W",
			latWant:   0,
			lngWant:   0,
			wantErr:   true,
		},
		{
			name:      "Invalid coordinates in decimal degree, in lat two dots",
			coordsGot: "26.8618.123N 89.37E",
			latWant:   0,
			lngWant:   0,
			wantErr:   true,
		},
		{
			name:      "Invalid coordinates in decimal degree, in lng two dots",
			coordsGot: "26.8618123N 89.3.7E",
			latWant:   0,
			lngWant:   0,
			wantErr:   true,
		},
	}

	var (
		point locodedb.Point
		err   error
	)

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			point, err = pointFromString(test.coordsGot)
			if (err != nil) != test.wantErr {
				t.Errorf("got error = %v, wantErr %v", err, test.wantErr)
			}
			if !(test.latWant == point.Latitude && test.lngWant == point.Longitude) {
				t.Errorf("got = %v, want %f, %f", point, test.latWant, test.lngWant)
			}
		})
	}
}

func pointFromString(s string) (locodedb.Point, error) {
	crd, err := CoordinatesFromString(s)
	if err != nil {
		return locodedb.Point{}, err
	}

	point, err := PointFromCoordinates(crd)
	if err != nil {
		return locodedb.Point{}, err
	}

	return point, nil
}

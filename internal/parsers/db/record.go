package locodedb

// Record represents a single record of the UN/LOCODE table.
type Record struct {
	// Combination of a 2-character country code and a 3-character location code.
	LOCODE [2]string

	// Name of the locations which has been allocated a UN/LOCODE.
	Name string

	// Names of the locations which have been allocated a UN/LOCODE without diacritic signs.
	NameWoDiacritics string

	// ISO 1-3 character alphabetic and/or numeric code for the administrative division of the country concerned.
	SubDiv string

	// 8-digit function classifier code for the location.
	Function string

	// Status of the entry by a 2-character code.
	Status string

	// Last date when the location was updated/entered.
	Date string

	// The IATA code for the location if different from location code in column LOCODE.
	IATA string

	// Geographical coordinates (latitude/longitude) of the location, if there is any.
	Coordinates string

	// Some general remarks regarding the UN/LOCODE in question.
	Remarks string
}

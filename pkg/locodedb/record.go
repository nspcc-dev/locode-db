package locodedb

// Record represents a record in the location database (resulting CSV files). It contains all the
// information about the location. It is used to fill the database. Country, Location are full names, codes are in Key.
type Record struct {
	Country    string
	Location   string
	SubDivName string
	SubDivCode string
	Point      Point
	Cont       Continent
}

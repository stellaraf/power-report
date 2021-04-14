package data

// LocationData represents summary data for a location.
type LocationData struct {
	current float64
	voltage float64
}

// DataPoint represents a single PDU and its name, location, and power data.
type DataPoint struct {
	device   string
	location string
	current  float64
	voltage  float64
}

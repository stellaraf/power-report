package data

type LocationData struct {
	current float64
	voltage float64
}

type DataPoint struct {
	device   string
	location string
	current  float64
	voltage  float64
}

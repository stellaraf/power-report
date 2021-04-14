package data

import (
	"fmt"
	"log"

	"github.com/tidwall/gjson"
)

/* filterDataPoints is like Array.filter from JS, probably. Blatantly stolen from:
https://stackoverflow.com/questions/37562873/most-idiomatic-way-to-select-elements-from-an-array-in-golang
*/
func filterDataPoints(ss []DataPoint, test func(DataPoint) bool) (ret []DataPoint) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

// getPDUs fetches the device IDs of all PDUs. Requires a device group called "PDUs".
func getPDUs() (ids []string) {
	b, err := request(baseURL, "api", "v0", "devicegroups", "PDUs")
	if err != nil {
		panic(err)
	}

	d := gjson.Parse(b)
	devices := d.Get("devices")
	devices.ForEach(func(idx, value gjson.Result) bool {
		ids = append(ids, value.Get("device_id").String())
		return true
	})
	return ids
}

// getSensors fetches current and voltage sensor data for each device ID.
func getSensors(deviceIDs []string) []DataPoint {
	var data []DataPoint
	for _, id := range deviceIDs {
		// Get the device object.
		d, err := request(baseURL, "api", "v0", "devices", id)
		if err != nil {
			panic(err)
		}
		device := gjson.Parse(d)
		deviceName := device.Get("devices.0.sysName").Str
		deviceLoc := device.Get("devices.0.location").Str

		// Get current sensor IDs.
		cSIDs, err := request(baseURL, "api", "v0", "devices", id, "health", "device_current")
		if err != nil {
			panic(err)
		}
		cData := gjson.Parse(cSIDs)
		cID := cData.Get("graphs.0.sensor_id").Int()

		// Get voltage sensor IDs.
		vSIDs, err := request(baseURL, "api", "v0", "devices", id, "health", "device_voltage")
		if err != nil {
			panic(err)
		}
		vData := gjson.Parse(vSIDs)
		vID := vData.Get("graphs.0.sensor_id").Int()

		// Get current sensor data by sensor ID.
		current, err := request(baseURL, "api", "v0", "devices", id, "health", "device_current", fmt.Sprintf("%d", cID))
		if err != nil {
			panic(err)
		}
		currentValue := gjson.Parse(current).Get("graphs.0.sensor_current").Num

		// Get voltage sensor data by sensor ID.
		voltage, err := request(baseURL, "api", "v0", "devices", id, "health", "device_voltage", fmt.Sprintf("%d", vID))
		if err != nil {
			panic(err)
		}
		voltageValue := gjson.Parse(voltage).Get("graphs.0.sensor_current").Num

		point := DataPoint{device: deviceName, location: deviceLoc, current: currentValue, voltage: voltageValue}
		data = append(data, point)
	}
	return data
}

// GetData finds all PDUs and queries their current and voltage sensors, grouping by location.
func GetData() (map[string]LocationData, map[string][]DataPoint) {
	ids := getPDUs()
	data := getSensors(ids)
	// Map of location → summary data.
	locations := make(map[string]LocationData)
	// Map of location → detail data.
	locationData := make(map[string][]DataPoint)

	for _, pdu := range data {
		value := LocationData{current: 0, voltage: 0}
		locations[pdu.location] = value
	}

	for loc := range locations {
		// Group detail by location.
		points := filterDataPoints(data, func(d DataPoint) bool {
			if d.location == loc {
				return true
			}
			return false
		})
		// Add summary data to summary map.
		for _, p := range points {
			value := LocationData{current: locations[loc].current + p.current, voltage: locations[loc].voltage + p.voltage}
			locations[loc] = value
		}
		locationData[loc] = points
	}
	log.Println("Finished gathering data from LibreNMS")
	return locations, locationData
}

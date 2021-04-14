package data

import (
	"fmt"
	"log"

	"github.com/tidwall/gjson"
)

/*Filter is like Array.filter from JS, probably. Blatantly stolen from:
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

func getPDUs() (ids []string) {
	headers = make(map[string]string)
	headers["X-Auth-Token"] = apiKey
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

func getSensors(deviceIDs []string) []DataPoint {
	var data []DataPoint
	for _, id := range deviceIDs {

		d, err := request(baseURL, "api", "v0", "devices", id)
		if err != nil {
			panic(err)
		}
		device := gjson.Parse(d)
		deviceName := device.Get("devices.0.sysName").Str
		deviceLoc := device.Get("devices.0.location").Str

		bCurrent, err := request(baseURL, "api", "v0", "devices", id, "health", "device_current")
		if err != nil {
			panic(err)
		}
		bVoltage, err := request(baseURL, "api", "v0", "devices", id, "health", "device_voltage")
		if err != nil {
			panic(err)
		}

		dCurrent := gjson.Parse(bCurrent)
		dVoltage := gjson.Parse(bVoltage)

		sCurrent := dCurrent.Get("graphs.0.sensor_id").Int()
		sVoltage := dVoltage.Get("graphs.0.sensor_id").Int()

		sCurrentData, err := request(baseURL, "api", "v0", "devices", id, "health", "device_current", fmt.Sprintf("%d", sCurrent))
		if err != nil {
			panic(err)
		}

		sVoltageData, err := request(baseURL, "api", "v0", "devices", id, "health", "device_voltage", fmt.Sprintf("%d", sVoltage))
		if err != nil {
			panic(err)
		}

		currentValue := gjson.Parse(sCurrentData).Get("graphs.0.sensor_current").Num
		voltageValue := gjson.Parse(sVoltageData).Get("graphs.0.sensor_current").Num
		point := DataPoint{device: deviceName, location: deviceLoc, current: currentValue, voltage: voltageValue}
		data = append(data, point)
	}
	return data
}

func GetData() (map[string]LocationData, map[string][]DataPoint) {
	ids := getPDUs()
	data := getSensors(ids)
	locations := make(map[string]LocationData)
	locationData := make(map[string][]DataPoint)
	for _, pdu := range data {
		value := LocationData{current: 0, voltage: 0}
		locations[pdu.location] = value
	}

	for loc := range locations {
		points := filterDataPoints(data, func(d DataPoint) bool {
			if d.location == loc {
				return true
			}
			return false
		})
		for _, p := range points {
			value := LocationData{current: locations[loc].current + p.current, voltage: locations[loc].voltage + p.voltage}
			locations[loc] = value
		}
		locationData[loc] = points
	}
	log.Println("Finished gathering data from LibreNMS")
	return locations, locationData
}

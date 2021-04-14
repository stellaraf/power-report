package main

import (
	"os"

	data "github.com/stellaraf/power-report/data"
)

func main() {
	// Query LibreNMS for the raw data.
	s, d := data.GetData()
	// Create a spreadsheet with per-PDU, per-location details.
	f := data.CreateSpreadsheet(d)
	// Send an email with the summary in the body and the spreadsheet attached.
	data.SendEmail(s, d, f)
	os.Exit(0)
}

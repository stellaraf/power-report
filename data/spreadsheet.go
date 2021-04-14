package data

import (
	"fmt"
	"log"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func CreateSpreadsheet(d map[string][]DataPoint) (fn string) {
	now := time.Now()
	f := excelize.NewFile()
	for loc := range d {
		end := len(d[loc]) + 2
		f.NewSheet(loc)
		f.SetCellValue(loc, "A1", "Device")
		f.SetCellValue(loc, "B1", "Current")
		f.SetCellValue(loc, "C1", "Voltage")
		for i, p := range d[loc] {
			c := i + 2
			f.SetCellValue(loc, fmt.Sprintf("A%d", c), p.device)
			f.SetCellValue(loc, fmt.Sprintf("B%d", c), p.current)
			f.SetCellValue(loc, fmt.Sprintf("C%d", c), p.voltage)
		}
		f.SetCellFormula(loc, fmt.Sprintf("B%d", end), fmt.Sprintf("=SUM(B2:B%d)", end-1))
		f.SetCellFormula(loc, fmt.Sprintf("C%d", end), fmt.Sprintf("=SUM(C2:C%d)", end-1))

	}
	f.DeleteSheet("Sheet1")
	fn = fmt.Sprintf("%d%02d%02d-stellar-power-report.xlsx", now.Year(), now.Month(), now.Day())
	err := f.SaveAs(fn)
	if err != nil {
		panic(err)
	}

	log.Printf("Wrote to file '%s'\n", fn)
	return fn
}

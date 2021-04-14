package data

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

func CreateEmailBody(date string, sum map[string]LocationData) (b string) {
	style := `<style>
table {
	width: 100%;
	border-collapse: collapse;
}
th, td {
	padding: 10px;
	text-align: left;
	vertical-align: middle;
}
table, th, td {
	border: 1px solid #ddd;
}
.time {
	opacity: 0.6;
}
</style>
`
	intro := fmt.Sprintf(`<h3>Data Center Power Utilization Report</h3>
<p class="time">%s</p>
`, date)
	tableBase := `<table>
<tr>
  <th>Location</th>
  <th>Total Current</th>
  <th>Total Voltage</th>
</tr>
%s
</table>
`
	var rows string
	for loc := range sum {
		rows += fmt.Sprintf(`<tr>
  <td>%s</td>
  <td>%.2f</td>
  <td>%.2f</td>
</tr>`, loc, sum[loc].current, sum[loc].voltage)
	}
	table := fmt.Sprintf(tableBase, rows)
	b = strings.Join([]string{style, intro, table}, "\n")
	return b
}

func SendEmail(sum map[string]LocationData, detail map[string][]DataPoint, f string) {
	now := time.Now()
	date := fmt.Sprintf("%s %02d, %d", now.Month(), now.Day(), now.Year())
	b := CreateEmailBody(date, sum)
	m := gomail.NewMessage()
	m.SetHeader("From", emailFrom)
	m.SetHeader("To", emailTo...)
	m.SetHeader("Subject", fmt.Sprintf("Data Center Power Utilization - %s", date))
	m.SetBody("text/html", b)
	m.Attach(f)
	port, err := strconv.Atoi(emailPort)
	if err != nil {
		panic(err)
	}
	dialer := gomail.Dialer{Host: emailHost, Port: port}

	err = dialer.DialAndSend(m)
	if err != nil {
		panic(err)
	}
	for _, r := range emailTo {
		log.Printf("Sent email to %s\n", r)
	}
}

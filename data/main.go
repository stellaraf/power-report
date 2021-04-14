package data

import (
	"log"
	"net/http"
	"os"
	"strings"
)

const requestTimeout int = 15

var httpClient *http.Client
var headers map[string]string
var baseURL string
var apiKey string
var emailTo []string
var emailFrom string
var emailHost string
var emailPort string

func init() {
	baseURL = os.Getenv("NMS_URL")
	if baseURL == "" {
		log.Fatalln("Unable to find required NMS_URL environment variable.")
		os.Exit(1)
	}
	apiKey = os.Getenv("NMS_API_KEY")
	if apiKey == "" {
		log.Fatalln("Unable to find required NMS_API_KEY environment variable.")
		os.Exit(1)
	}
	emailToEnv := os.Getenv("POWER_REPORT_EMAIL_TO")
	if emailToEnv == "" {
		log.Fatalln("Unable to find required POWER_REPORT_EMAIL_TO environment variable.")
		os.Exit(1)
	} else {
		emailTo = strings.Split(emailToEnv, ",")
	}
	emailFrom = os.Getenv("POWER_REPORT_EMAIL_FROM")
	if emailFrom == "" {
		log.Fatalln("Unable to find required POWER_REPORT_EMAIL_FROM environment variable.")
		os.Exit(1)
	}
	emailHost = os.Getenv("POWER_REPORT_EMAIL_HOST")
	if emailHost == "" {
		log.Fatalln("Unable to find required POWER_REPORT_EMAIL_HOST environment variable.")
		os.Exit(1)
	}
	emailPort = os.Getenv("POWER_REPORT_EMAIL_PORT")
	if emailPort == "" {
		log.Fatalln("Unable to find required POWER_REPORT_EMAIL_PORT environment variable.")
		os.Exit(1)
	}
	httpClient = createClient()
}

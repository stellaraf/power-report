package data

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

func createClient() *http.Client {
	jar, err := cookiejar.New(nil)

	if err != nil {
		panic(err)
	}

	client := &http.Client{Jar: jar, Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}, Timeout: time.Duration(requestTimeout) * time.Second}

	return client
}

// Perform an HTTP GET request and return the body as a string.
func request(u string, p ...string) (b string, e error) {
	// Ensure the base URL has a trailing slash.
	if string(u[len(u)-1:]) != "/" {
		u = u + "/"
	}
	reqURL, err := url.Parse(u + strings.Join(p, "/"))

	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest(http.MethodGet, reqURL.String(), nil)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err := httpClient.Do(req)

	if err != nil {
		return "", fmt.Errorf("Request to %s failed:\n%s", reqURL, err.Error())
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("Error requesting data from %s - %s", reqURL, res.Status)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", fmt.Errorf("Unable to parse response from %s\n%s", reqURL, err.Error())
	}

	return string(body), nil
}

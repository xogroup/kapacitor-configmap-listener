package configuration

import (
	"strings"

	"github.com/influxdata/influxdb/client/v2"
)

// NewInfluxClient create a Influx Client based on flags
func NewInfluxClient(url string, username string, password string, ssl bool, unsafeSsl bool) (client.Client, error) {

	if url == "" {
		url = "localhost:8086"
	}

	if strings.HasPrefix(url, "http") == false {
		if ssl == true {
			url = "https://" + url
		} else {
			url = "http://" + url
		}
	}

	return client.NewHTTPClient(client.HTTPConfig{
		Addr:               url,
		InsecureSkipVerify: unsafeSsl,
		Username:           username,
		Password:           password,
	})
}

package configuration

import (
	"strings"

	"github.com/influxdata/kapacitor/client/v1"
)

// NewKapacitorClient create a Kapacitor Client based on flags
func NewKapacitorClient(url string) (*client.Client, error) {

	if url == "" {
		url = "localhost:9092"
	}

	if strings.HasPrefix(url, "http") == false {
		url = "http://" + url
	}

	var config = client.Config{URL: url}

	return client.New(config)
}

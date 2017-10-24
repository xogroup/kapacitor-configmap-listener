package influx

import (
	"fmt"

	"github.com/influxdata/influxdb/client/v2"

	log "github.com/sirupsen/logrus"
)

// Subscriptions is a model for handling influx subscription cleanup
type Subscriptions struct {
	influxClient client.Client
}

type deleteSubscriptionCommand struct {
	name          string
	rp            string
	database      string
	deleteCommand string
}

// NewSubscriptions instantiates a new object of that type
func NewSubscriptions(influxClient client.Client) *Subscriptions {

	return &Subscriptions{
		influxClient: influxClient,
	}
}

// RemoveAll subscriptions on the server
func (subscriptions *Subscriptions) RemoveAll() error {

	response, err := subscriptions.influxClient.Query(buildQuery("show subscriptions"))

	if err != nil {
		return err
	}

	deleteSubscriptionCommands := createDeleteCommandsFromSubscriptionList(*response)

	for _, deleteSubscriptionCommand := range deleteSubscriptionCommands {
		log.Infof("Removing %s from %s", deleteSubscriptionCommand.name, deleteSubscriptionCommand.database)
		_, err := subscriptions.influxClient.Query(buildQuery(deleteSubscriptionCommand.deleteCommand))
		if err != nil {
			return err
		}
	}

	return nil
}

func buildQuery(command string) client.Query {
	return client.Query{
		Database: "kubernetes.2-weeks",
		Command:  command,
	}
}

func createDeleteCommandsFromSubscriptionList(response client.Response) []deleteSubscriptionCommand {
	slice := make([]deleteSubscriptionCommand, 0)

	for _, serie := range response.Results[0].Series {
		for _, value := range serie.Values {
			deleteSubscriptionCommand := deleteSubscriptionCommand{
				database:      serie.Name,
				rp:            value[0].(string),
				name:          value[1].(string),
				deleteCommand: fmt.Sprintf("drop subscription \"%s\" on \"%s\".\"%s\"", value[1].(string), serie.Name, value[0].(string)),
			}
			slice = append(slice, deleteSubscriptionCommand)
		}
	}

	return slice
}

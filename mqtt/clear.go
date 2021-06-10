package mqtt

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	mqttclient "github.com/eclipse/paho.mqtt.golang"
	"github.com/fatih/structtag"
	"github.com/koesie10/ws-upload/wsupload"
)

func DeleteAllDevices(options PublisherOptions) error {
	hostname, _ := os.Hostname()

	connOpts := mqttclient.NewClientOptions().SetClientID(fmt.Sprintf("%s-%d", hostname, time.Now().Unix())).SetCleanSession(true)

	for _, broker := range options.Brokers {
		connOpts.AddBroker(broker)
	}

	if options.Username != "" {
		connOpts.SetUsername(options.Username)
		if options.Password != "" {
			connOpts.SetPassword(options.Password)
		}
	}
	connOpts.SetAutoReconnect(true)
	connOpts.SetConnectRetry(true)

	client := mqttclient.NewClient(connOpts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	reflectType := reflect.TypeOf(wsupload.Observation{})

	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)

		tag, err := structtag.Parse(string(field.Tag))
		if err != nil {
			return fmt.Errorf("failed to parse struct tag for %s: %w", field.Name, err)
		}

		jsonTag, err := tag.Get("json")
		if err != nil {
			continue
		}

		v := struct{}{}

		topic := fmt.Sprintf("%s/sensor/%s%s/config", options.HomeAssistant.DiscoveryPrefix, options.HomeAssistant.DevicePrefix, jsonTag.Name)

		data, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal observation to JSON: %w", err)
		}

		client.Publish(topic, byte(options.HomeAssistant.DiscoveryQoS), true, string(data))
	}

	return nil
}

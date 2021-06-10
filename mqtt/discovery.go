package mqtt

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/fatih/structtag"
	"github.com/koesie10/ws-upload/wsupload"
	"github.com/koesie10/ws-upload/x"
	"github.com/sirupsen/logrus"
)

type homeAssistantConfig struct {
	DeviceClass       string `json:"device_class,omitempty"`
	Name              string `json:"name"`
	StateTopic        string `json:"state_topic"`
	StateClass        string `json:"state_class,omitempty"`
	UnitOfMeasurement string `json:"unit_of_measurement,omitempty"`
	ValueTemplate     string `json:"value_template"`
}

func (p *publisher) publishDiscovery() error {
	if !p.options.HomeAssistant.DiscoveryEnabled {
		return nil
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
		homeAssistantTag, err := tag.Get("homeassistant")
		if err != nil {
			logrus.Warnf("Field %s is missing homeassistant tag", field.Name)
			continue
		}

		options := x.ParseStructTagOptions(homeAssistantTag.Options)

		config := homeAssistantConfig{
			DeviceClass:       options["device_class"],
			Name:              homeAssistantTag.Name,
			StateTopic:        p.options.Topic,
			StateClass:        options["state_class"],
			UnitOfMeasurement: options["unit_of_measurement"],
			ValueTemplate:     fmt.Sprintf("{{ value.%s }}", jsonTag.Name),
		}

		topic := fmt.Sprintf("%s/sensor/%s%s/config", p.options.HomeAssistant.DiscoveryPrefix, p.options.HomeAssistant.DevicePrefix, jsonTag.Name)

		data, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal observation to JSON: %w", err)
		}

		token := p.client.Publish(topic, byte(p.options.HomeAssistant.DiscoveryQoS), true, string(data))
		go func(topic string) {
			token.Wait()
			if err := token.Error(); err != nil {
				logrus.WithError(err).Warn("Failed to publish config %s to MQTT", topic)
			}
		}(topic)
	}

	return nil
}

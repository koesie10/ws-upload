package mqtt

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/fatih/structtag"
	"github.com/koesie10/ws-upload/wsupload"
	"github.com/koesie10/ws-upload/x"
	"go.uber.org/zap"
)

type homeAssistantDevice struct {
	Identifiers  []string `json:"identifiers,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	Model        string   `json:"model,omitempty"`
	Name         string   `json:"name,omitempty"`
}

type homeAssistantConfig struct {
	DeviceClass       string `json:"device_class,omitempty"`
	Name              string `json:"name"`
	StateTopic        string `json:"state_topic"`
	StateClass        string `json:"state_class,omitempty"`
	UnitOfMeasurement string `json:"unit_of_measurement,omitempty"`
	ValueTemplate     string `json:"value_template"`

	UniqueID string              `json:"unique_id,omitempty"`
	Device   homeAssistantDevice `json:"device"`
}

func (p *publisher) publishDiscovery() error {
	if !p.options.HomeAssistant.DiscoveryEnabled {
		return nil
	}

	reflectType := reflect.TypeOf(wsupload.Observation{})

	device := homeAssistantDevice{
		Identifiers:  p.options.HomeAssistant.DeviceIdentifiers,
		Manufacturer: p.options.HomeAssistant.DeviceManufacturer,
		Model:        p.options.HomeAssistant.DeviceModel,
		Name:         p.options.HomeAssistant.DeviceName,
	}

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
			p.logger.Warn("Field is missing homeassistant tag", zap.String("discovery.field", field.Name))
			continue
		}

		options := x.ParseStructTagOptions(homeAssistantTag.Options)

		config := homeAssistantConfig{
			DeviceClass:       options["device_class"],
			Name:              homeAssistantTag.Name,
			StateTopic:        p.options.Topic,
			StateClass:        options["state_class"],
			UnitOfMeasurement: options["unit_of_measurement"],
			ValueTemplate:     fmt.Sprintf("{{ value_json.%s }}", jsonTag.Name),

			UniqueID: fmt.Sprintf("%s%s", p.options.HomeAssistant.UniqueIDPrefix, jsonTag.Name),
			Device:   device,
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
				p.logger.Warn("Failed to publish config to MQTT", zap.String("mqtt.topic", topic), zap.Error(err))
			}
		}(topic)
	}

	return nil
}

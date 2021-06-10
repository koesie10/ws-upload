package mqtt

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	mqttclient "github.com/eclipse/paho.mqtt.golang"
	"github.com/koesie10/ws-upload/wsupload"
	"github.com/sirupsen/logrus"
)

var _ wsupload.Publisher = (*publisher)(nil)

type publisher struct {
	client mqttclient.Client

	options PublisherOptions

	done chan struct{}
}

func NewPublisher(options PublisherOptions) (wsupload.Publisher, error) {
	if options.Debug {
		setupDebugLogs()
	}

	hostname, _ := os.Hostname()

	if options.ClientID == "" {
		options.ClientID = fmt.Sprintf("%s-%d", hostname, time.Now().Unix())
	}

	connOpts := mqttclient.NewClientOptions().SetClientID(options.ClientID).SetCleanSession(true)

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

	p := &publisher{
		client:  client,
		options: options,

		done: make(chan struct{}),
	}

	go p.watchdog()

	return p, nil
}

type PublisherOptions struct {
	Brokers  []string `env:"MQTT_BROKERS" flag:"brokers" desc:"MQTT broker addresses, leave empty to disable"`
	ClientID string   `env:"MQTT_CLIENT_ID" flag:"client-id" desc:"MQTT client ID, default will be autogenerated based on the client hostname"`
	Username string   `env:"MQTT_USERNAME" flag:"username" desc:"MQTT username"`
	Password string   `env:"MQTT_PASSWORD" flag:"password" desc:"MQTT password"`

	Topic string `env:"MQTT_TOPIC" flag:"topic" desc:"topic to publish to"`
	QoS   int    `env:"MQTT_QOS" flag:"qos" desc:"the QoS to send the messages at"`

	HomeAssistant HomeAssistantOptions `env:",squash"`

	Debug bool `env:"MQTT_DEBUG" flag:"debug" desc:"whether to enable debug logging"`
}

type HomeAssistantOptions struct {
	DiscoveryEnabled  bool          `env:"MQTT_HOMEASSISTANT_DISCOVERY_ENABLED" flag:"discovery-enabled" desc:"whether HomeAssistant MQTT discovery is enabled"`
	DiscoveryPrefix   string        `env:"MQTT_HOMEASSISTANT_DISCOVERY_PREFIX" flag:"discovery-prefix" desc:"HomeAssistant MQTT discovery prefix"`
	DiscoveryQoS      int           `env:"MQTT_HOMEASSISTANT_DISCOVERY_QOS" flag:"discovery-qos" desc:"HomeAssistant MQTT discovery QoS"`
	DiscoveryInterval time.Duration `env:"MQTT_HOMEASSISTANT_DISCOVERY_INTERVAL" flag:"discovery-interval" desc:"HomeAssistant MQTT discovery interval"`
	DevicePrefix      string        `env:"MQTT_HOMEASSISTANT_DEVICE_PREFIX" flag:"device-prefix" desc:"HomeAssistant device prefix"`
}

func (p *publisher) Publish(obs *wsupload.Observation) error {
	data, err := json.Marshal(obs)
	if err != nil {
		return fmt.Errorf("failed to marshal observation to JSON: %w", err)
	}

	token := p.client.Publish(p.options.Topic, byte(p.options.QoS), true, string(data))
	go func() {
		token.Wait()
		if err := token.Error(); err != nil {
			logrus.WithError(err).Warn("Failed to publish observation to MQTT")
		}
	}()

	return nil
}

func (p *publisher) Close() error {
	close(p.done)

	return nil
}

func (p *publisher) watchdog() {
	token := p.client.Connect()

	token.Wait()

	if token.Error() != nil {
		logrus.WithError(token.Error()).Errorf("Failed to connect to MQTT broker")
	}

	discoveryInterval := p.options.HomeAssistant.DiscoveryInterval
	if discoveryInterval == 0 {
		discoveryInterval = 30 * time.Second
	}

	t := time.NewTicker(discoveryInterval)

	logrus.Infof("Connected to MQTT broker")

	if err := p.publishDiscovery(); err != nil {
		logrus.Warnf("Failed to publish discovery message")
	}

	select {
	case <-p.done:
		p.client.Disconnect(250)

		return
	case <-t.C:
		if err := p.publishDiscovery(); err != nil {
			logrus.Warnf("Failed to publish discovery message")
		}
	}
}

func setupDebugLogs() {
	mqttclient.DEBUG = &logger{
		level: logrus.TraceLevel,
		entry: logrus.WithField("system", "mqtt"),
	}
	mqttclient.WARN = &logger{
		level: logrus.InfoLevel,
		entry: logrus.WithField("system", "mqtt"),
	}
	mqttclient.ERROR = &logger{
		level: logrus.WarnLevel,
		entry: logrus.WithField("system", "mqtt"),
	}
	mqttclient.CRITICAL = &logger{
		level: logrus.ErrorLevel,
		entry: logrus.WithField("system", "mqtt"),
	}
}

type logger struct {
	level logrus.Level
	entry *logrus.Entry
}

func (l *logger) Println(v ...interface{}) {
	l.entry.Log(l.level, v...)
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.entry.Logf(l.level, format, v...)
}

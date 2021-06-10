package influx

import (
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/koesie10/ws-upload/wsupload"
)

var _ wsupload.Publisher = (*publisher)(nil)

type publisher struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI

	options PublisherOptions
}

func NewPublisher(options PublisherOptions) (wsupload.Publisher, error) {
	influxOptions := influxdb2.DefaultOptions()
	influxOptions.SetPrecision(time.Second)
	client := influxdb2.NewClientWithOptions(options.Addr, options.AuthToken, influxOptions)

	writeAPI := client.WriteAPI(options.Organization, options.Bucket)

	return &publisher{
		client:   client,
		writeAPI: writeAPI,
		options:  options,
	}, nil
}

type PublisherOptions struct {
	Addr            string `env:"INFLUX_ADDR" flag:"addr" desc:"InfluxDB HTTP address, set empty to disable"`
	AuthToken       string `env:"INFLUX_AUTH_TOKEN" flag:"auth-token" desc:"InfluxDB auth token, use username:password for InfluxDB 1.8"`
	Organization    string `env:"INFLUX_ORGANIZATION" flag:"organization" desc:"InfluxDB organization, do not set if using InfluxDB 1.8"`
	Bucket          string `env:"INFLUX_BUCKET" flag:"bucket" desc:"InfluxDB bucket, set to database/retention-policy or database for InfluxDB 1.8"`
	MeasurementName string `env:"MEASUREMENT_NAME" flag:"measurement-name" desc:"InfluxDB measurement name"`
}

func (p *publisher) Publish(obs *wsupload.Observation) error {
	point, err := CreatePoint(obs, p.options.MeasurementName)
	if err != nil {
		return fmt.Errorf("failed to create point: %w", err)
	}

	p.writeAPI.WritePoint(point)

	return nil
}

func (p *publisher) Close() error {
	p.client.Close()

	return nil
}

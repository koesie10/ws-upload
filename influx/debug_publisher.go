package influx

import (
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/koesie10/ws-upload/wsupload"
	"time"
)

var _ wsupload.Publisher = (*debugPublisher)(nil)

type debugPublisher struct {
	options DebugPublisherOptions
}

func NewDebugPublisher(options DebugPublisherOptions) (wsupload.Publisher, error) {
	return &debugPublisher{
		options: options,
	}, nil
}

type DebugPublisherOptions struct {
	MeasurementName string `env:"MEASUREMENT_NAME" flag:"measurement-name" desc:"InfluxDB measurement name"`
}

func (p *debugPublisher) Publish(obs *wsupload.Observation) error {
	point, err := CreatePoint(obs, p.options.MeasurementName)
	if err != nil {
		return fmt.Errorf("failed to create point: %w", err)
	}

	fmt.Printf("INFLUX DEBUG: %s", write.PointToLineProtocol(point, time.Millisecond))

	return nil
}

func (p *debugPublisher) Close() error {
	return nil
}

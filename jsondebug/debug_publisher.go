package jsondebug

import (
	"encoding/json"
	"fmt"

	"github.com/koesie10/ws-upload/wsupload"
)

var _ wsupload.Publisher = (*debugPublisher)(nil)

type debugPublisher struct{}

func NewDebugPublisher() (wsupload.Publisher, error) {
	return &debugPublisher{}, nil
}

func (p *debugPublisher) Publish(obs *wsupload.Observation) error {
	data, err := json.Marshal(obs)
	if err != nil {
		return fmt.Errorf("failed to marshal observation to JSON: %w", err)
	}

	fmt.Printf("JSON DEBUG: %s\n", string(data))

	return nil
}

func (p *debugPublisher) Close() error {
	return nil
}

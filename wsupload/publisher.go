package wsupload

type Publisher interface {
	Publish(obs *Observation) error

	Close() error
}

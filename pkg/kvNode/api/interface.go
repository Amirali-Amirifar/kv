package api

type TransportHandler interface {
	Serve(port int) error
}

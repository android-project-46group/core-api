package config

const (
	grpcPort = "9876"
	service  = "saka-core-api"
)

// Configuration values used throughout the application.
type Config struct {
	// The port number on which gRPC is running.
	GrpcPort string

	// Service name (also used for tracer)
	Service string
}

// New returns the config
func New() (Config, error) {
	c := Config{
		GrpcPort: grpcPort,
		Service:  service,
	}

	return c, nil
}

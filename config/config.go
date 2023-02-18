package config

const (
	grpcPort = "9876"
	service  = "saka-core-api"
	logPath  = "./core-api-log.txt"
)

// Configuration values used throughout the application.
type Config struct {
	// The port number on which gRPC is running.
	GrpcPort string

	// Service name (also used for tracer)
	Service string

	// File path to log output.
	LogPath string
}

// New returns the config.
func New() (Config, error) {
	cfg := Config{
		GrpcPort: grpcPort,
		Service:  service,
		LogPath:  logPath,
	}

	return cfg, nil
}

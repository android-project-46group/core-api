package config

const (
	grpcPort = "9876"
	service  = "saka-core-api"
	logPath  = "./core-api-log.txt"

	dbDriver   = "postgres"
	dbUser     = "ubuntu"
	dbPassword = "sakamichi"
	dbName     = "sakamichi"
	dbHost     = "localhost"
	dbPort     = "5432"
)

// Configuration values used throughout the application.
type Config struct {
	// The port number on which gRPC is running.
	GrpcPort string

	// Service name (also used for tracer)
	Service string

	// File path to log output.
	LogPath string

	// Driver name of the Database.
	DBDriver string

	// User name of the Database.
	DBUser string

	// Password of the Database.
	DBPassword string

	// Database name.
	DBName string

	// Host name of the Database.
	DBHost string

	// Port number of the Database.
	DBPort string
}

// New returns the config.
func New() (Config, error) {
	cfg := Config{
		GrpcPort: grpcPort,
		Service:  service,
		LogPath:  logPath,

		DBDriver:   dbDriver,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBHost:     dbHost,
		DBPort:     dbPort,
	}

	return cfg, nil
}

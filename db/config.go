package db

import (
	"fmt"
)

// Config holds the database configuration from the environment
type Config struct {
	User       string `envconfig:"db_user"`
	Host       string `envconfig:"db_host"`
	Port       string `envconfig:"db_port"`
	Database   string `envconfig:"db_name"`
	Password   string `envconfig:"db_pass"`
	SSLEnabled string `envconfig:"db_ssl_enabled" default:"disable"`
}

// ConnectionString gets the connection string from the environment variables
func (c *Config) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Database,
		c.SSLEnabled,
	)
}

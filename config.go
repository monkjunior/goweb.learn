package main

import (
	"fmt"
	"strings"
)

type Config struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Pepper   string         `json:"pepper"`
	HMACKey  string         `json:"hmac_key"`
	Database PostgresConfig `json:"database"`
	Mailgun  MailgunConfig  `json:"mailgun"`
}

func DefaultConfig() Config {
	return Config{
		Port:     8080,
		Env:      "dev",
		Pepper:   "ted-is-so-handsome",
		HMACKey:  "secret-hmac-key",
		Database: DefaultPostgresConfig(),
	}
}

func (c Config) IsProd() bool {
	env := strings.ToLower(c.Env)
	if env == "prod" || env == "production" {
		return true
	}
	return false
}

type PostgresConfig struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
}

func (c *PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Name,
		)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name,
	)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "ted",
		Password: "your-password",
		Name:     "goweb_dev",
	}
}

type MailgunConfig struct {
	ApiKey       string `json:"api_key"`
	PublicApiKey string `json:"public_api_key"`
	Domain       string `json:"domain"`
}

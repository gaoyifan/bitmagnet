package server

import "time"

type Config struct {
	LocalAddress string
	Port         uint16
	QueryTimeout time.Duration
}

func NewDefaultConfig() Config {
	return Config{
		LocalAddress: "0.0.0.0",
		Port:         3334,
		QueryTimeout: time.Second * 4,
	}
}

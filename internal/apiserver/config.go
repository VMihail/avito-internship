package apiserver

import (
	_ "encoding/json"
	_ "os"
)

type Config struct {
	BindAddr string `json:"bind_addr"`
	LogLevel string `json:"log_level"`
}

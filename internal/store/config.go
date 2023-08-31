package store

type Config struct {
	DbUrl      string `json:"database_url"`
	DriverName string `json:"driver_name"`
	LogLevel   string `json:"log_level"`
}

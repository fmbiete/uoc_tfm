package config

type Config struct {
	Database ConfigDatabase `json:"database"`
	Server   ConfigServer   `json:"server"`
}

type ConfigDatabase struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

type ConfigServer struct {
	Port      int    `json:"port"`
	JWTSecret string `json:"jwt_secret"`
}

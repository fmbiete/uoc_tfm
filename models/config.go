package models

type Config struct {
	Database   ConfigDatabase `json:"database"`
	Server     ConfigServer   `json:"server"`
	SiteAdmin  User           `json:"site_admin"`
	SiteConfig Configuration  `json:"site_config"`
}

type ConfigDatabase struct {
	Host         string `json:"host"`
	User         string `json:"user"`
	Password     string `json:"password"`
	Port         int    `json:"port"`
	Database     string `json:"database"`
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
	Reset        bool   `json:"reset"`
}

type ConfigServer struct {
	Port      int    `json:"port"`
	JWTSecret string `json:"jwt_secret"`
}

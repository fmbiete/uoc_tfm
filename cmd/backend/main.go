package main

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"tfm_backend/api"
	"tfm_backend/models"
	"tfm_backend/orm"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	var cfg models.Config
	raw, err := os.ReadFile("config.json")
	if err != nil {
		log.Error().Err(err).Msg("Failed to read config file")
		return
	}
	json.Unmarshal(raw, &cfg)

	// Railway - read db credentials from environment variables
	if len(cfg.Database.Host) == 0 {
		cfg.Database.Host = os.Getenv("DATABASE_HOST")
		cfg.Database.User = os.Getenv("DATABASE_USER")
		cfg.Database.Password = os.Getenv("DATABASE_PASSWORD")
	}

	database := orm.NewDatabase(&cfg)
	server := api.NewServer(cfg.Server, database)

	err = database.Setup()
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize database")
		return
	}

	err = server.Listen()
	if err != nil {
		log.Error().Err(err).Msg("Faile to listen REST API")
		return
	}
}

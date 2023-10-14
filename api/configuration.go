package api

import (
	"net/http"
	"tfm_backend/orm"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) ConfigurationDetails(c echo.Context) error {
	config, err := s.db.ConfigurationDetails()
	if err != nil {
		log.Error().Err(err).Msg("Failed to read configuracion")
		return err
	}

	return c.JSON(http.StatusOK, config)
}

func (s *Server) ConfigurationModify(c echo.Context) error {
	var config orm.Configuration
	err := c.Bind(&config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind configuracion")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	config, err = s.db.ConfigurationModify(config)
	if err != nil {
		log.Error().Err(err).Interface("config", config).Msg("Failed to modify configuracion")
		return err
	}

	return c.JSON(http.StatusOK, config)
}

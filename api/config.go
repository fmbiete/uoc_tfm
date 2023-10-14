package api

import (
	"net/http"
	"tfm_backend/orm"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) ConfiguracionDetails(c echo.Context) error {
	config, err := s.db.ConfiguracionDetails()
	if err != nil {
		log.Error().Err(err).Msg("Failed to read configuracion")
		return err
	}

	return c.JSON(http.StatusOK, config)
}

func (s *Server) ConfiguracionModify(c echo.Context) error {
	var config orm.Configuracion
	err := c.Bind(&config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind configuracion")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	config, err = s.db.ConfiguracionModify(config)
	if err != nil {
		log.Error().Err(err).Interface("config", config).Msg("Failed to modify configuracion")
		return err
	}

	return c.JSON(http.StatusOK, config)
}

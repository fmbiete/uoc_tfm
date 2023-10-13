package api

import (
	"net/http"
	"strconv"
	"tfm_backend/orm"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) PlatoCreate(c echo.Context) error {
	var plato orm.Plato
	err := c.Bind(&plato)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind plato")
		return c.NoContent(http.StatusBadRequest)
	}

	plato, err = s.db.PlatoCreate(plato)
	if err != nil {
		log.Error().Err(err).Interface("plato", plato).Msg("Failed to create plato")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, plato)
}

func (s *Server) PlatoDelete(c echo.Context) error {
	platoId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	err = s.db.PlatoDelete(platoId)
	if err != nil {
		log.Error().Err(err).Int64("id", platoId).Msg("Failed to delete plato")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) PlatoDetails(c echo.Context) error {
	platoId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	plato, err := s.db.PlatoDetails(platoId)
	if err != nil {
		log.Error().Err(err).Interface("id", platoId).Msg("Failed to read plato")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, plato)
}

func (s *Server) PlatoList(c echo.Context) error {
	var userId int64 = -1
	var err error
	if len(c.QueryParam("usuarioId")) > 0 {
		userId, err = strconv.ParseInt(c.QueryParam("usuarioId"), 10, 64)
		if err != nil {
			log.Error().Err(err).Str("id", c.QueryParam("usuarioId")).Msg("Failed to convert usuarioId to int64")
			return c.NoContent(http.StatusBadRequest)
		}
	}

	platos, err := s.db.PlatoList(userId)
	if err != nil {
		log.Error().Err(err).Interface("usuarioId", userId).Msg("Failed to list platos")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, platos)
}

func (s *Server) PlatoModify(c echo.Context) error {
	platoId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	var plato orm.Plato
	err = c.Bind(&plato)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind plato")
		return c.NoContent(http.StatusBadRequest)
	}
	plato.ID = uint(platoId)

	plato, err = s.db.PlatoModify(plato)
	if err != nil {
		log.Error().Err(err).Interface("plato", plato).Msg("Failed to modify plato")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, plato)
}

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
	platoId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	err = s.db.PlatoDelete(platoId)
	if err != nil {
		log.Error().Err(err).Uint64("id", platoId).Msg("Failed to delete plato")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) PlatoDetails(c echo.Context) error {
	platoId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	plato, err := s.db.PlatoDetails(platoId)
	if err != nil {
		log.Error().Err(err).Uint64("id", platoId).Msg("Failed to read plato")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, plato)
}

func (s *Server) PlatoList(c echo.Context) error {
	var userId int64 = -1
	if authenticated(c) {
		userId = int64(authenticatedUserId(c))
	}

	platos, err := s.db.PlatoList(userId)
	if err != nil {
		log.Error().Err(err).Int64("usuarioId", userId).Msg("Failed to list platos")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, platos)
}

func (s *Server) PlatoModify(c echo.Context) error {
	platoId, err := strconv.ParseUint(c.Param("id"), 10, 64)
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
	plato.ID = platoId

	plato, err = s.db.PlatoModify(plato)
	if err != nil {
		log.Error().Err(err).Interface("plato", plato).Msg("Failed to modify plato")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, plato)
}

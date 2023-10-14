package api

import (
	"net/http"
	"strconv"
	"tfm_backend/orm"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) PromocionCreate(c echo.Context) error {
	var promocion orm.Promocion
	err := c.Bind(&promocion)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind promocion")
		return c.NoContent(http.StatusBadRequest)
	}

	promocion, err = s.db.PromocionCreate(promocion)
	if err != nil {
		log.Error().Err(err).Interface("promocion", promocion).Msg("Failed to create promocion")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, promocion)
}

func (s *Server) PromocionDelete(c echo.Context) error {
	promocionId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	err = s.db.PromocionDelete(promocionId)
	if err != nil {
		log.Error().Err(err).Uint64("id", promocionId).Msg("Failed to delete promocion")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) PromocionDetails(c echo.Context) error {
	promocionId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	promocion, err := s.db.PromocionDetails(promocionId)
	if err != nil {
		log.Error().Err(err).Uint64("id", promocionId).Msg("Failed to read promocion")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, promocion)
}

func (s *Server) PromocionList(c echo.Context) error {
	promocions, err := s.db.PromocionList()
	if err != nil {
		log.Error().Err(err).Msg("Failed to list promociones")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, promocions)
}

func (s *Server) PromocionModify(c echo.Context) error {
	promocionId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	var promocion orm.Promocion
	err = c.Bind(&promocion)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind promocion")
		return c.NoContent(http.StatusBadRequest)
	}
	promocion.ID = uint(promocionId)

	promocion, err = s.db.PromocionModify(promocion)
	if err != nil {
		log.Error().Err(err).Interface("promocion", promocion).Msg("Failed to modify promocion")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, promocion)
}

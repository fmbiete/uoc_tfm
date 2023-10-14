package api

import (
	"net/http"
	"strconv"
	"tfm_backend/orm"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) CarritoDelete(c echo.Context) error {
	userId, err := strconv.ParseUint(c.Param("usuarioid"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("usuarioid", c.Param("usuarioid")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	carrito, err := s.db.CarritoDelete(userId)
	if err != nil {
		log.Error().Err(err).Interface("carrito", carrito).Msg("Failed to clear carrito")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, carrito)
}

func (s *Server) CarritoDetails(c echo.Context) error {
	userId, err := strconv.ParseUint(c.Param("usuarioid"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("usuarioid", c.Param("usuarioid")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	carrito, err := s.db.CarritoDetails(userId)
	if err != nil {
		log.Error().Err(err).Interface("carrito", carrito).Msg("Failed to read carrito")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, carrito)
}

func (s *Server) CarritoSave(c echo.Context) error {
	userId, err := strconv.ParseUint(c.Param("usuarioid"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("usuarioid", c.Param("usuarioid")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	var carrito orm.Carrito
	err = c.Bind(&carrito)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind carrito")
		return c.NoContent(http.StatusBadRequest)
	}
	carrito.UsuarioID = userId

	carrito, err = s.db.CarritoSave(carrito)
	if err != nil {
		log.Error().Err(err).Interface("carrito", carrito).Msg("Failed to save carrito")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, carrito)
}

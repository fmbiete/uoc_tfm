package api

import (
	"net/http"
	"tfm_backend/orm"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) CarritoDelete(c echo.Context) error {
	carrito, err := s.db.CarritoDelete(authenticatedUserId(c))
	if err != nil {
		log.Error().Err(err).Interface("carrito", carrito).Msg("Failed to clear carrito")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, carrito)
}

func (s *Server) CarritoDetails(c echo.Context) error {
	carrito, err := s.db.CarritoDetails(authenticatedUserId(c))
	if err != nil {
		log.Error().Err(err).Interface("carrito", carrito).Msg("Failed to read carrito")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, carrito)
}

func (s *Server) CarritoSave(c echo.Context) error {
	var carrito orm.Carrito
	err := c.Bind(&carrito)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind carrito")
		return c.NoContent(http.StatusBadRequest)
	}
	carrito.UsuarioID = authenticatedUserId(c)

	carrito, err = s.db.CarritoSave(carrito)
	if err != nil {
		log.Error().Err(err).Interface("carrito", carrito).Msg("Failed to save carrito")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, carrito)
}

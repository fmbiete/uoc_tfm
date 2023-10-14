package api

import (
	"net/http"
	"tfm_backend/orm"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) CartDelete(c echo.Context) error {
	cart, err := s.db.CartDelete(authenticatedUserId(c))
	if err != nil {
		log.Error().Err(err).Interface("cart", cart).Msg("Failed to clear cart")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, cart)
}

func (s *Server) CartDetails(c echo.Context) error {
	cart, err := s.db.CartDetails(authenticatedUserId(c))
	if err != nil {
		log.Error().Err(err).Interface("cart", cart).Msg("Failed to read cart")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, cart)
}

func (s *Server) CartSave(c echo.Context) error {
	var cart orm.Cart
	err := c.Bind(&cart)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind cart")
		return c.NoContent(http.StatusBadRequest)
	}
	cart.UserID = authenticatedUserId(c)

	cart, err = s.db.CartSave(cart)
	if err != nil {
		log.Error().Err(err).Interface("cart", cart).Msg("Failed to save cart")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, cart)
}

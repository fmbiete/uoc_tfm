package api

import (
	"net/http"
	"strconv"
	"tfm_backend/orm"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) UsuarioCrear(c echo.Context) error {
	var user orm.Usuario
	err := c.Bind(&user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind usuario")
		return c.NoContent(http.StatusBadRequest)
	}

	if len(user.Password) > 0 {
		user.Password = stringToSha512(user.Password)
	}
	user, err = s.db.UsuarioCrear(user)
	if err != nil {
		log.Error().Err(err).Interface("user", user).Msg("Failed to create user")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, user)
}

func (s *Server) UsuarioEliminar(c echo.Context) error {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg("Failed to convert ID to int64")
		return c.NoContent(http.StatusBadRequest)
	}

	err = s.db.UsuarioEliminar(userId)
	if err != nil {
		log.Error().Err(err).Int64("id", userId).Msg("Failed to delete user")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) UsuarioGet(c echo.Context) error {
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg("Failed to convert ID to int64")
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := s.db.UsuarioGet(userId)
	if err != nil {
		log.Error().Err(err).Int64("id", userId).Msg("Failed to read user")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, user)
}

func (s *Server) UsuarioModificar(c echo.Context) error {
	var user orm.Usuario
	err := c.Bind(&user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind usuario")
		return c.NoContent(http.StatusBadRequest)
	}

	// If we have a new password, we generate the hash
	if len(user.Password) > 0 {
		user.Password = stringToSha512(user.Password)
	}

	user, err = s.db.UsuarioModificar(user)
	if err != nil {
		log.Error().Err(err).Interface("user", user).Msg("Failed to modify user")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, user)
}

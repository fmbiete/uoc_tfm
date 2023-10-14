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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(user.Password) > 0 {
		user.Password = stringToSha512(user.Password)
	}
	user, err = s.db.UsuarioCrear(user)
	if err != nil {
		log.Error().Err(err).Interface("user", user).Msg("Failed to create user")
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func (s *Server) UsuarioEliminar(c echo.Context) error {
	var userId = authenticatedUserId(c)
	var err error

	if authenticatedIsRestaurador(c) {
		// Only a restaurador can delete other users
		userId, err = strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	err = s.db.UsuarioEliminar(userId)
	if err != nil {
		log.Error().Err(err).Uint64("id", userId).Msg("Failed to delete user")
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) UsuarioGet(c echo.Context) error {
	var userId = authenticatedUserId(c)
	var err error

	if authenticatedIsRestaurador(c) {
		// Only a restaurador can read other users
		userId, err = strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	user, err := s.db.UsuarioGet(userId)
	if err != nil {
		log.Error().Err(err).Uint64("id", userId).Msg("Failed to read user")
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (s *Server) UsuarioModificar(c echo.Context) error {
	var userId = authenticatedUserId(c)
	var err error

	if authenticatedIsRestaurador(c) {
		// Only a restaurador can modify other users
		userId, err = strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	var user orm.Usuario
	err = c.Bind(&user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind usuario")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user.ID = userId

	// If we have a new password, we generate the hash
	if len(user.Password) > 0 {
		user.Password = stringToSha512(user.Password)
	}

	user, err = s.db.UsuarioModificar(user)
	if err != nil {
		log.Error().Err(err).Interface("user", user).Msg("Failed to modify user")
		return err
	}

	return c.JSON(http.StatusOK, user)
}

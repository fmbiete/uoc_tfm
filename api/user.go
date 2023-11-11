package api

import (
	"net/http"
	"strconv"
	"tfm_backend/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) UserCreate(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind user")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(user.Password) > 0 {
		user.Password = stringToSha512(user.Password)
	}
	user, err = s.db.UserCreate(user)
	if err != nil {
		log.Error().Err(err).Interface("user", user).Msg("Failed to create user")
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func (s *Server) UserDelete(c echo.Context) error {
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

	// we don't allow deletion of "admin" user
	if userId == 1 {
		return echo.NewHTTPError(http.StatusForbidden, "Initial User cannot be deleted")
	}

	err = s.db.UserDelete(userId)
	if err != nil {
		log.Error().Err(err).Uint64("id", userId).Msg("Failed to delete user")
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) UserDetails(c echo.Context) error {
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

	user, err := s.db.UserDetails(userId)
	if err != nil {
		log.Error().Err(err).Uint64("id", userId).Msg("Failed to read user")
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (s *Server) UserList(c echo.Context) error {
	limit, page, offset := parsePagination(c)

	users, err := s.db.UserList(limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list users")
		return err
	}

	return c.JSON(http.StatusOK, models.PaginationUsers{Limit: limit, Page: page, Users: users})
}

func (s *Server) UserModify(c echo.Context) error {
	var authUserId = authenticatedUserId(c)
	var err error

	// user
	userId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var user models.User
	err = c.Bind(&user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind user")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user.ID = userId

	if authenticatedIsRestaurador(c) {
		// An administrator cannot remove its own admin access (mistake protection)
		if authUserId == userId && !user.IsAdmin {
			log.Warn().Uint64("authUserId", authUserId).Msg(`An Admin cannot remove its own admin access, ask another admin`)
			return echo.NewHTTPError(http.StatusForbidden, `An Administrator cannot remove its own administrative access, please ask another Administrator`)
		}
	} else {
		// Only a restaurador can modify other users
		if authUserId != userId {
			log.Warn().Uint64("authUserId", authUserId).Uint64("userId", userId).Msg(`A Non-Admin user is trying to modify another user`)
			return echo.NewHTTPError(http.StatusForbidden, `Only an Administrator can modify another user`)
		}
	}

	// If we have a new password, we generate the hash
	if len(user.Password) > 0 {
		user.Password = stringToSha512(user.Password)
	}

	user, err = s.db.UserModify(user)
	if err != nil {
		log.Error().Err(err).Interface("user", user).Msg("Failed to modify user")
		return err
	}

	return c.JSON(http.StatusOK, user)
}

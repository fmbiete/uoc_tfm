package api

import (
	"net/http"
	"strconv"
	"tfm_backend/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) DishCreate(c echo.Context) error {
	var dish models.Dish
	err := c.Bind(&dish)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind dish")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dish, err = s.db.DishCreate(dish)
	if err != nil {
		log.Error().Err(err).Interface("dish", dish).Msg("Failed to create dish")
		return err
	}

	return c.JSON(http.StatusCreated, dish)
}

func (s *Server) DishDelete(c echo.Context) error {
	dishId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = s.db.DishDelete(dishId)
	if err != nil {
		log.Error().Err(err).Uint64("id", dishId).Msg("Failed to delete dish")
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) DishDetails(c echo.Context) error {
	dishId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dish, err := s.db.DishDetails(dishId)
	if err != nil {
		log.Error().Err(err).Uint64("id", dishId).Msg("Failed to read dish")
		return err
	}

	return c.JSON(http.StatusOK, dish)
}

func (s *Server) DishDislike(c echo.Context) error {
	dishId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = s.db.DishDislike(authenticatedUserId(c), dishId)
	if err != nil {
		log.Error().Err(err).Uint64("id", dishId).Msg("Failed to dislike dish")
		return err
	}

	return c.JSON(http.StatusOK, true)
}

func (s *Server) DishFavourites(c echo.Context) error {
	var userId int64 = -1
	if authenticated(c) {
		userId = int64(authenticatedUserId(c))
	}

	limit, page, offset := parsePagination(c)

	dishes, err := s.db.DishFavourites(userId, limit, offset)
	if err != nil {
		log.Error().Err(err).Int64("userId", userId).Msg("Failed to list favourite dishes")
		return err
	}

	return c.JSON(http.StatusOK, models.PaginationDishes{Limit: limit, Page: page, Dishes: dishes})
}

func (s *Server) DishLike(c echo.Context) error {
	dishId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = s.db.DishLike(authenticatedUserId(c), dishId)
	if err != nil {
		log.Error().Err(err).Uint64("id", dishId).Msg("Failed to like dish")
		return err
	}

	return c.JSON(http.StatusOK, true)
}

func (s *Server) DishList(c echo.Context) error {
	limit, page, offset := parsePagination(c)

	dishes, err := s.db.DishList(limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list dishes")
		return err
	}

	return c.JSON(http.StatusOK, models.PaginationDishes{Limit: limit, Page: page, Dishes: dishes})
}

func (s *Server) DishModify(c echo.Context) error {
	dishId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var dish models.Dish
	err = c.Bind(&dish)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind dish")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	dish.ID = dishId

	dish, err = s.db.DishModify(dish)
	if err != nil {
		log.Error().Err(err).Interface("dish", dish).Msg("Failed to modify dish")
		return err
	}

	return c.JSON(http.StatusOK, dish)
}

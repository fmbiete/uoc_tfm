package api

import (
	"net/http"
	"strconv"
	"tfm_backend/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) IngredientCreate(c echo.Context) error {
	var ingredient models.Ingredient
	err := c.Bind(&ingredient)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind Ingredient")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ingredient, err = s.db.IngredientCreate(ingredient)
	if err != nil {
		log.Error().Err(err).Interface("ingredient", ingredient).Msg("Failed to create Ingredient")
		return err
	}

	return c.JSON(http.StatusCreated, ingredient)
}

func (s *Server) IngredientDetails(c echo.Context) error {
	ingredientId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ingredient, err := s.db.IngredientDetails(ingredientId)
	if err != nil {
		log.Error().Err(err).Uint64("id", ingredientId).Msg("Failed to read Ingredient")
		return err
	}

	return c.JSON(http.StatusOK, ingredient)
}

func (s *Server) IngredientDelete(c echo.Context) error {
	ingredientId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = s.db.IngredientDelete(ingredientId)
	if err != nil {
		log.Error().Err(err).Uint64("id", ingredientId).Msg("Failed to delete Ingredient")
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) IngredientDishes(c echo.Context) error {
	ingredientId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	limit, page, offset := parsePagination(c)

	dishes, err := s.db.IngredientDishes(ingredientId, limit, offset)
	if err != nil {
		log.Error().Err(err).Uint64("id", ingredientId).Msg("Failed to read Ingredient Dishes")
		return err
	}

	return c.JSON(http.StatusOK, models.PaginationDishes{Limit: limit, Page: page, Dishes: dishes})
}

func (s *Server) IngredientList(c echo.Context) error {
	categories, err := s.db.IngredientList()
	if err != nil {
		log.Error().Err(err).Msg("Failed to list Ingredient")
		return err
	}

	return c.JSON(http.StatusOK, categories)
}

func (s *Server) IngredientModify(c echo.Context) error {
	ingredientId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var ingredient models.Ingredient
	err = c.Bind(&ingredient)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind Ingredient")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ingredient.ID = ingredientId

	ingredient, err = s.db.IngredientModify(ingredient)
	if err != nil {
		log.Error().Err(err).Interface("ingredient", ingredient).Msg("Failed to modify Ingredient")
		return err
	}

	return c.JSON(http.StatusOK, ingredient)
}

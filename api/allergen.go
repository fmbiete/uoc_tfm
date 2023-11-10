package api

import (
	"net/http"
	"strconv"
	"tfm_backend/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) AllergenCreate(c echo.Context) error {
	var allergen models.Allergen
	err := c.Bind(&allergen)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind Allergen")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	allergen, err = s.db.AllergenCreate(allergen)
	if err != nil {
		log.Error().Err(err).Interface("allergen", allergen).Msg("Failed to create Allergen")
		return err
	}

	return c.JSON(http.StatusCreated, allergen)
}

func (s *Server) AllergenDetails(c echo.Context) error {
	allergenId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	allergen, err := s.db.AllergenDetails(allergenId)
	if err != nil {
		log.Error().Err(err).Uint64("id", allergenId).Msg("Failed to read Allergen")
		return err
	}

	return c.JSON(http.StatusOK, allergen)
}

func (s *Server) AllergenDelete(c echo.Context) error {
	allergenId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = s.db.AllergenDelete(allergenId)
	if err != nil {
		log.Error().Err(err).Uint64("id", allergenId).Msg("Failed to delete Allergen")
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) AllergenDishes(c echo.Context) error {
	allergenId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	limit, page, offset := parsePagination(c)

	dishes, err := s.db.AllergenDishes(allergenId, limit, offset)
	if err != nil {
		log.Error().Err(err).Uint64("id", allergenId).Msg("Failed to read Allergen Dishes")
		return err
	}

	return c.JSON(http.StatusOK, models.PaginationDishes{Limit: limit, Page: page, Dishes: dishes})
}

func (s *Server) AllergenList(c echo.Context) error {
	categories, err := s.db.AllergenList()
	if err != nil {
		log.Error().Err(err).Msg("Failed to list Allergen")
		return err
	}

	return c.JSON(http.StatusOK, categories)
}

func (s *Server) AllergenModify(c echo.Context) error {
	allergenId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var allergen models.Allergen
	err = c.Bind(&allergen)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind Allergen")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	allergen.ID = allergenId

	allergen, err = s.db.AllergenModify(allergen)
	if err != nil {
		log.Error().Err(err).Interface("allergen", allergen).Msg("Failed to modify Allergen")
		return err
	}

	return c.JSON(http.StatusOK, allergen)
}

package api

import (
	"net/http"
	"strconv"
	"tfm_backend/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) CategoryCreate(c echo.Context) error {
	var category models.Category
	err := c.Bind(&category)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind Category")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	category, err = s.db.CategoryCreate(category)
	if err != nil {
		log.Error().Err(err).Interface("category", category).Msg("Failed to create Category")
		return err
	}

	return c.JSON(http.StatusCreated, category)
}

func (s *Server) CategoryDetails(c echo.Context) error {
	categoryId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	category, err := s.db.CategoryDetails(categoryId)
	if err != nil {
		log.Error().Err(err).Uint64("id", categoryId).Msg("Failed to read Category")
		return err
	}

	return c.JSON(http.StatusOK, category)
}

func (s *Server) CategoryList(c echo.Context) error {
	categories, err := s.db.CategoryList()
	if err != nil {
		log.Error().Err(err).Msg("Failed to list Category")
		return err
	}

	return c.JSON(http.StatusOK, categories)
}

func (s *Server) CategoryModify(c echo.Context) error {
	categoryId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var category models.Category
	err = c.Bind(&category)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind Category")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	category.ID = categoryId

	category, err = s.db.CategoryModify(category)
	if err != nil {
		log.Error().Err(err).Interface("category", category).Msg("Failed to modify Category")
		return err
	}

	return c.JSON(http.StatusOK, category)
}

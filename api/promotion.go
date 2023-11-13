package api

import (
	"net/http"
	"strconv"
	"tfm_backend/models"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) PromotionCreate(c echo.Context) error {
	var promotion models.Promotion
	err := c.Bind(&promotion)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind promotion")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	promotion, err = s.db.PromotionCreate(promotion)
	if err != nil {
		log.Error().Err(err).Interface("promotion", promotion).Msg("Failed to create promotion")
		return err
	}

	return c.JSON(http.StatusCreated, promotion)
}

func (s *Server) PromotionDelete(c echo.Context) error {
	promotionId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = s.db.PromotionDelete(promotionId)
	if err != nil {
		log.Error().Err(err).Uint64("id", promotionId).Msg("Failed to delete promotion")
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) PromotionDetails(c echo.Context) error {
	promotionId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	promotion, err := s.db.PromotionDetails(promotionId)
	if err != nil {
		log.Error().Err(err).Uint64("id", promotionId).Msg("Failed to read promotion")
		return err
	}

	return c.JSON(http.StatusOK, promotion)
}

func (s *Server) PromotionList(c echo.Context) error {
	limit, page, offset := parsePagination(c)

	activeOnly, _ := strconv.ParseBool(c.QueryParam("activeOnly"))

	promotions, err := s.db.PromotionList(activeOnly, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list promotiones")
		return err
	}

	return c.JSON(http.StatusOK, models.PaginationPromotion{Limit: limit, Page: page, Promotions: promotions})
}

func (s *Server) PromotionModify(c echo.Context) error {
	promotionId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var promotion models.Promotion
	err = c.Bind(&promotion)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind promotion")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	promotion.ID = promotionId

	promotion, err = s.db.PromotionModify(promotion)
	if err != nil {
		log.Error().Err(err).Interface("promotion", promotion).Msg("Failed to modify promotion")
		return err
	}

	return c.JSON(http.StatusOK, promotion)
}

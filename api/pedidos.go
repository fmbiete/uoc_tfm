package api

import (
	"net/http"
	"strconv"
	"tfm_backend/orm"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) PedidoCreate(c echo.Context) error {
	var pedido orm.Pedido
	err := c.Bind(&pedido)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind pedido")
		return c.NoContent(http.StatusBadRequest)
	}

	pedido, err = s.db.PedidoCreate(pedido)
	if err != nil {
		log.Error().Err(err).Interface("pedido", pedido).Msg("Failed to create pedido")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, pedido)
}

func (s *Server) PedidoDelete(c echo.Context) error {
	pedidoId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	err = s.db.PedidoDelete(pedidoId)
	if err != nil {
		log.Error().Err(err).Uint64("id", pedidoId).Msg("Failed to delete pedido")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) PedidoDetails(c echo.Context) error {
	pedidoId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	pedido, err := s.db.PedidoDetails(pedidoId)
	if err != nil {
		log.Error().Err(err).Uint64("id", pedidoId).Msg("Failed to read pedido")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pedido)
}

func (s *Server) PedidoList(c echo.Context) error {
	var userId int64 = -1
	var err error
	if len(c.QueryParam("usuarioId")) > 0 {
		userId, err = strconv.ParseInt(c.QueryParam("usuarioId"), 10, 64)
		if err != nil {
			log.Error().Err(err).Str("id", c.QueryParam("usuarioId")).Msg("Failed to convert usuarioId to int64")
			return c.NoContent(http.StatusBadRequest)
		}
	}

	pedidos, err := s.db.PedidoList(userId)
	if err != nil {
		log.Error().Err(err).Int64("usuario", userId).Msg("Failed to list pedidos")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pedidos)
}

func (s *Server) PedidoModify(c echo.Context) error {
	pedidoId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	var pedido orm.Pedido
	err = c.Bind(&pedido)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind pedido")
		return c.NoContent(http.StatusBadRequest)
	}
	pedido.ID = uint(pedidoId)

	pedido, err = s.db.PedidoModify(pedido)
	if err != nil {
		log.Error().Err(err).Interface("pedido", pedido).Msg("Failed to modify pedido")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pedido)
}

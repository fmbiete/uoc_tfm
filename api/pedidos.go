package api

import (
	"net/http"
	"strconv"
	"tfm_backend/orm"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) PedidoCancel(c echo.Context) error {
	pedidoId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return c.NoContent(http.StatusBadRequest)
	}

	err = s.db.PedidoDelete(pedidoId)
	if err != nil {
		log.Error().Err(err).Uint64("id", pedidoId).Msg("Failed to delete pedido")
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) PedidoCreateFromCarrito(c echo.Context) error {
	pedido, err := s.db.PedidoCreateFromCarrito(authenticatedUserId(c))
	if err != nil {
		log.Error().Err(err).Interface("pedido", pedido).Msg("Failed to create pedido from carrito")
		return err
	}

	return c.JSON(http.StatusCreated, pedido)
}

func (s *Server) PedidoDetails(c echo.Context) error {
	pedidoId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pedido, err := s.db.PedidoDetails(pedidoId)
	if err != nil {
		log.Error().Err(err).Uint64("id", pedidoId).Msg("Failed to read pedido")
		return err
	}

	return c.JSON(http.StatusOK, pedido)
}

func (s *Server) PedidoList(c echo.Context) error {
	var userId int64 = int64(authenticatedUserId(c))
	if authenticatedIsRestaurador(c) {
		userId = -1
	}

	var dayFilter string = c.QueryParam("day")

	pedidos, err := s.db.PedidoList(userId, dayFilter)
	if err != nil {
		log.Error().Err(err).Int64("usuario", userId).Str("day", dayFilter).Msg("Failed to list pedidos")
		return err
	}

	return c.JSON(http.StatusOK, pedidos)
}

func (s *Server) PedidoLineaCreate(c echo.Context) error {
	pedidoId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var linea orm.CarritoLinea
	err = c.Bind(&linea)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind carrito linea")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pedido, err := s.db.PedidoLineaCreate(pedidoId, linea)
	if err != nil {
		log.Error().Err(err).Interface("pedido", pedido).Msg("Failed to create pedido linea")
		return err
	}

	return c.JSON(http.StatusOK, pedido)
}
func (s *Server) PedidoLineaDelete(c echo.Context) error {
	pedidoId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	lineaId, err := strconv.ParseUint(c.Param("lineaid"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("lineaid", c.Param("lineaid")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pedido, err := s.db.PedidoLineaDelete(pedidoId, lineaId)
	if err != nil {
		log.Error().Err(err).Interface("pedido", pedido).Msg("Failed to delete pedido linea")
		return err
	}

	return c.JSON(http.StatusOK, pedido)
}

func (s *Server) PedidoLineaModify(c echo.Context) error {
	pedidoId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("id", c.Param("id")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	lineaId, err := strconv.ParseUint(c.Param("lineaid"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("lineaid", c.Param("lineaid")).Msg(msgErrorIdToInt)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var linea orm.PedidoLinea
	err = c.Bind(&linea)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind pedidoLinea")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	linea.PedidoID = pedidoId
	linea.ID = lineaId

	pedido, err := s.db.PedidoLineaModify(linea)
	if err != nil {
		log.Error().Err(err).Interface("pedido", pedido).Msg("Failed to modify pedido linea")
		return err
	}

	return c.JSON(http.StatusOK, pedido)
}

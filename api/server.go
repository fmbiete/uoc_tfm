package api

import (
	"fmt"
	"net/http"
	"os"
	"tfm_backend/config"
	"tfm_backend/orm"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ziflex/lecho/v3"
)

type Server struct {
	e   *echo.Echo
	db  *orm.Database
	cfg *config.ConfigServer
}

const msgErrorIdToInt = "Failed to convert ID to int64"

func NewServer(cfg config.ConfigServer, db *orm.Database) *Server {
	s := Server{e: echo.New(), cfg: &cfg, db: db}
	s.e.HideBanner = true
	s.e.Logger = lecho.New(os.Stdout, lecho.WithLevel(log.DEBUG), lecho.WithTimestamp(), lecho.WithCaller())

	// s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())

	return &s
}

func (s *Server) Listen() error {
	var requiresLogin = echojwt.WithConfig(echojwt.Config{SigningKey: []byte(s.cfg.JWTSecret)})

	s.e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "TFM Backend API")
	})

	// Config API
	gConfig := s.e.Group("/config")
	gConfig.GET("/", s.ConfiguracionDetails, requiresLogin, requiresRestaurador)
	gConfig.PATCH("/", s.ConfiguracionModify, requiresLogin, requiresRestaurador)

	// Usuario API
	gUsuario := s.e.Group("/usuarios")
	gUsuario.POST("/login", s.Login)
	gUsuario.POST("/", s.UsuarioCrear)
	gUsuario.GET("/:id", s.UsuarioGet)
	gUsuario.PATCH("/:id", s.UsuarioModificar, requiresLogin)
	gUsuario.DELETE("/:id", s.UsuarioEliminar, requiresLogin)

	// Platos API
	gPlatos := s.e.Group("/platos")
	gPlatos.GET("/", s.PlatoList)
	gPlatos.GET("/:id", s.PlatoDetails)
	gPlatos.POST("/", s.PlatoCreate, requiresLogin, requiresRestaurador)
	gPlatos.PATCH("/:id", s.PlatoModify, requiresLogin, requiresRestaurador)
	gPlatos.DELETE("/:id", s.PlatoDelete, requiresLogin, requiresRestaurador)

	// Promociones API
	gPromociones := s.e.Group("/promociones")
	gPromociones.GET("/", s.PromocionList)
	gPromociones.GET("/:id", s.PromocionDetails)
	gPromociones.POST("/", s.PromocionCreate, requiresLogin, requiresRestaurador)
	gPromociones.PATCH("/:id", s.PromocionModify, requiresLogin, requiresRestaurador)
	gPromociones.DELETE("/:id", s.PromocionDelete, requiresLogin, requiresRestaurador)

	// Carrito API
	gCarritos := s.e.Group("/carritos")
	gCarritos.GET("/", s.CarritoDetails, requiresLogin)
	gCarritos.POST("/", s.CarritoSave, requiresLogin)
	gCarritos.DELETE("/", s.CarritoDelete, requiresLogin)

	// Pedidos API
	gPedidos := s.e.Group("/pedidos")
	gPedidos.POST("/", s.PedidoCreateFromCarrito, requiresLogin)
	gPedidos.GET("/", s.PedidoList, requiresLogin)
	gPedidos.GET("/:id", s.PedidoDetails, requiresLogin)
	gPedidos.DELETE("/:id", s.PedidoCancel, requiresLogin)
	gPedidos.POST("/:id/linea/", s.PedidoLineaCreate, requiresLogin)
	gPedidos.PATCH("/:id/linea/:lineaid", s.PedidoLineaModify, requiresLogin)
	gPedidos.DELETE("/:id/linea/:lineaid", s.PedidoLineaDelete, requiresLogin)

	return s.e.Start(fmt.Sprintf(`:%d`, s.cfg.Port))
}

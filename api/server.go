package api

import (
	"errors"
	"fmt"
	"net/http"
	"tfm_backend/config"
	"tfm_backend/orm"

	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v3"
)

type Server struct {
	e             *echo.Echo
	db            *orm.Database
	cfg           *config.ConfigServer
	requiresLogin echo.MiddlewareFunc
	optionalLogin echo.MiddlewareFunc
}

const msgErrorIdToInt = "Failed to convert ID to int64"

func NewServer(cfg config.ConfigServer, db *orm.Database) *Server {
	s := Server{e: echo.New(), cfg: &cfg, db: db}

	s.requiresLogin = echojwt.WithConfig(echojwt.Config{SigningKey: []byte(s.cfg.JWTSecret)})
	s.optionalLogin = echojwt.WithConfig(
		echojwt.Config{
			SigningKey:             []byte(s.cfg.JWTSecret),
			ContinueOnIgnoredError: true,
			ErrorHandler: func(c echo.Context, err error) error {
				fmt.Println(err)
				if errors.Is(err, echojwt.ErrJWTMissing) {
					return nil
				}
				return err
			},
		},
	)

	s.e.HideBanner = true
	s.e.Logger = lecho.From(log.Logger)

	// s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())

	s.e.HTTPErrorHandler = customHTTPErrorHandler

	return &s
}

func customHTTPErrorHandler(err error, c echo.Context) {
	uuid := uuid.NewString()
	log.Error().Err(err).Str("uuid", uuid).Msg("Reflection")
	if he, ok := err.(*echo.HTTPError); ok {
		c.JSON(he.Code, map[string]interface{}{"message": he.Message, "reflection": uuid})
	} else {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err.Error(), "reflection": uuid})
	}
}

func (s *Server) Listen() error {
	s.e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "TFM Backend API")
	})

	// Config API
	gConfig := s.e.Group("/config")
	gConfig.GET("/", s.ConfiguracionDetails, s.requiresLogin, requiresRestaurador)
	gConfig.PATCH("/", s.ConfiguracionModify, s.requiresLogin, requiresRestaurador)

	// Usuario API
	gUsuario := s.e.Group("/usuarios")
	gUsuario.POST("/login", s.Login)
	gUsuario.POST("/", s.UsuarioCrear)
	gUsuario.GET("/:id", s.UsuarioGet)
	gUsuario.PATCH("/:id", s.UsuarioModificar, s.requiresLogin)
	gUsuario.DELETE("/:id", s.UsuarioEliminar, s.requiresLogin)

	// Platos API
	gPlatos := s.e.Group("/platos")
	// platos is authenticated (show list of favourite platos for user) and unauthenticated (show list of favourite platos for everybody)
	gPlatos.GET("/", s.PlatoList, s.optionalLogin)
	gPlatos.GET("/:id", s.PlatoDetails)
	gPlatos.POST("/", s.PlatoCreate, s.requiresLogin, requiresRestaurador)
	gPlatos.PATCH("/:id", s.PlatoModify, s.requiresLogin, requiresRestaurador)
	gPlatos.DELETE("/:id", s.PlatoDelete, s.requiresLogin, requiresRestaurador)

	// Promociones API
	gPromociones := s.e.Group("/promociones")
	gPromociones.GET("/", s.PromocionList)
	gPromociones.GET("/:id", s.PromocionDetails)
	gPromociones.POST("/", s.PromocionCreate, s.requiresLogin, requiresRestaurador)
	gPromociones.PATCH("/:id", s.PromocionModify, s.requiresLogin, requiresRestaurador)
	gPromociones.DELETE("/:id", s.PromocionDelete, s.requiresLogin, requiresRestaurador)

	// Carrito API
	gCarritos := s.e.Group("/carritos")
	gCarritos.GET("/", s.CarritoDetails, s.requiresLogin)
	gCarritos.POST("/", s.CarritoSave, s.requiresLogin)
	gCarritos.DELETE("/", s.CarritoDelete, s.requiresLogin)

	// Pedidos API
	gPedidos := s.e.Group("/pedidos")
	gPedidos.POST("/", s.PedidoCreateFromCarrito, s.requiresLogin)
	gPedidos.GET("/", s.PedidoList, s.requiresLogin)
	gPedidos.GET("/:id", s.PedidoDetails, s.requiresLogin)
	gPedidos.DELETE("/:id", s.PedidoCancel, s.requiresLogin)
	gPedidos.POST("/:id/linea/", s.PedidoLineaCreate, s.requiresLogin)
	gPedidos.PATCH("/:id/linea/:lineaid", s.PedidoLineaModify, s.requiresLogin)
	gPedidos.DELETE("/:id/linea/:lineaid", s.PedidoLineaDelete, s.requiresLogin)

	return s.e.Start(fmt.Sprintf(`:%d`, s.cfg.Port))
}

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

func NewServer(cfg config.ConfigServer, db *orm.Database) *Server {
	s := Server{e: echo.New(), cfg: &cfg, db: db}
	s.e.HideBanner = true
	s.e.Logger = lecho.New(os.Stdout, lecho.WithLevel(log.DEBUG), lecho.WithTimestamp(), lecho.WithCaller())
	return &s
}

func (s *Server) Listen() error {
	var requiresLogin = echojwt.WithConfig(echojwt.Config{SigningKey: []byte(s.cfg.JWTSecret)})
	s.e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "TFM Backend API")
	})

	s.e.Use(middleware.Recover())

	// Usuario API
	gUsuario := s.e.Group("/usuarios")
	gUsuario.POST("/login", s.Login)
	gUsuario.POST("/", s.UsuarioCrear)
	gUsuario.GET("/:id", s.UsuarioGet)
	gUsuario.PATCH("/:id", s.UsuarioModificar, requiresLogin)
	gUsuario.DELETE("/:id", s.UsuarioEliminar, requiresLogin)

	gPlatos := s.e.Group("/platos")
	gPlatos.GET("/", s.PlatoList)
	gPlatos.GET("/:id", s.PlatoDetails)
	gPlatos.POST("/", s.PlatoCreate, requiresLogin, requiresRestaurador)
	gPlatos.PATCH("/:id", s.PlatoModify, requiresLogin, requiresRestaurador)
	gPlatos.DELETE("/:id", s.PlatoDelete, requiresLogin, requiresRestaurador)

	s.e.GET("/todo", func(c echo.Context) error { return c.String(http.StatusOK, "OK") }, requiresLogin, requiresRestaurador)

	return s.e.Start(fmt.Sprintf(`:%d`, s.cfg.Port))
}

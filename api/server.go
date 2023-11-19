package api

import (
	"errors"
	"fmt"
	"net/http"
	"tfm_backend/models"
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
	cfg           *models.ConfigServer
	requiresLogin echo.MiddlewareFunc
	optionalLogin echo.MiddlewareFunc
}

const msgErrorIdToInt = "Failed to convert ID to int64"

func NewServer(cfg models.ConfigServer, db *orm.Database) *Server {
	s := Server{e: echo.New(), cfg: &cfg, db: db}

	s.requiresLogin = echojwt.WithConfig(echojwt.Config{SigningKey: []byte(s.cfg.JWTSecret)})
	s.optionalLogin = echojwt.WithConfig(
		echojwt.Config{
			SigningKey:             []byte(s.cfg.JWTSecret),
			ContinueOnIgnoredError: true,
			ErrorHandler: func(c echo.Context, err error) error {
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

	s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},
	}))

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

	// Configuration API
	gConfiguration := s.e.Group("/configuration")
	gConfiguration.GET("/", s.ConfigurationDetails, s.requiresLogin, requiresAdministrator)
	gConfiguration.PATCH("/", s.ConfigurationModify, s.requiresLogin, requiresAdministrator)

	// User API
	gUser := s.e.Group("/user")
	gUser.POST("/login", s.Login)
	gUser.POST("/", s.UserCreate)
	gUser.GET("/:id", s.UserDetails, s.requiresLogin)
	gUser.PATCH("/:id", s.UserModify, s.requiresLogin)
	gUser.DELETE("/:id", s.UserDelete, s.requiresLogin)
	s.e.GET("/users", s.UserList, s.requiresLogin, requiresAdministrator)
	s.e.GET("/users/count", s.UserCount, s.requiresLogin, requiresAdministrator)

	// Allergens API
	gAllergen := s.e.Group("/allergen")
	gAllergen.POST("/", s.AllergenCreate, s.requiresLogin, requiresAdministrator)
	gAllergen.GET("/:id", s.AllergenDetails)
	gAllergen.PATCH("/:id", s.AllergenModify, s.requiresLogin, requiresAdministrator)
	gAllergen.DELETE("/:id", s.AllergenDelete, s.requiresLogin, requiresAdministrator)
	gAllergen.GET("/:id/dishes", s.AllergenDishes)
	s.e.GET("/allergens", s.AllergenList)

	// Categories API
	gCategory := s.e.Group("/category")
	gCategory.POST("/", s.CategoryCreate, s.requiresLogin, requiresAdministrator)
	gCategory.GET("/:id", s.CategoryDetails)
	gCategory.PATCH("/:id", s.CategoryModify, s.requiresLogin, requiresAdministrator)
	gCategory.DELETE("/:id", s.CategoryDelete, s.requiresLogin, requiresAdministrator)
	gCategory.GET("/:id/dishes", s.CategoryDishes)
	s.e.GET("/categories", s.CategoryList)

	// Ingredients API
	gIngredient := s.e.Group("/ingredient")
	gIngredient.POST("/", s.IngredientCreate, s.requiresLogin, requiresAdministrator)
	gIngredient.GET("/:id", s.IngredientDetails)
	gIngredient.PATCH("/:id", s.IngredientModify, s.requiresLogin, requiresAdministrator)
	gIngredient.DELETE("/:id", s.IngredientDelete, s.requiresLogin, requiresAdministrator)
	gIngredient.GET("/:id/dishes", s.IngredientDishes)
	s.e.GET("/ingredients", s.IngredientList)

	// Dishes API
	gDishes := s.e.Group("/dish")
	// /favourites is authenticated (show list of favourite dishes for user) and unauthenticated (show list of favourite dishes for everybody)
	gDishes.GET("/favourites", s.DishFavourites, s.optionalLogin)
	// /:id is authenticated (show like/dislike for user) and authenticated (don't show like/dislike)
	gDishes.GET("/:id", s.DishDetails, s.optionalLogin)
	gDishes.POST("/", s.DishCreate, s.requiresLogin, requiresAdministrator)
	gDishes.PATCH("/:id", s.DishModify, s.requiresLogin, requiresAdministrator)
	gDishes.DELETE("/:id", s.DishDelete, s.requiresLogin, requiresAdministrator)
	s.e.GET("/dishes", s.DishList, s.optionalLogin)
	s.e.GET("/dishes/count", s.DishCount, s.requiresLogin, requiresAdministrator)
	gDishes.POST("/:id/like", s.DishLike, s.requiresLogin)
	gDishes.POST("/:id/dislike", s.DishDislike, s.requiresLogin)

	// Promotions API
	gPromotions := s.e.Group("/promotion")
	gPromotions.GET("/:id", s.PromotionDetails)
	gPromotions.POST("/", s.PromotionCreate, s.requiresLogin, requiresAdministrator)
	gPromotions.PATCH("/:id", s.PromotionModify, s.requiresLogin, requiresAdministrator)
	gPromotions.DELETE("/:id", s.PromotionDelete, s.requiresLogin, requiresAdministrator)
	s.e.GET("/promotions", s.PromotionList)
	s.e.GET("/promotions/count", s.PromotionCount, s.requiresLogin, requiresAdministrator)

	// Orders API
	gOrders := s.e.Group("/order")
	gOrders.GET("/subvention", s.OrderSubvention, s.requiresLogin)
	gOrders.POST("/", s.OrderCreate, s.requiresLogin)
	gOrders.GET("/:id", s.OrderDetails, s.requiresLogin)
	gOrders.DELETE("/:id", s.OrderDelete, s.requiresLogin)
	gOrders.POST("/:id/line/", s.OrderLineCreate, s.requiresLogin)
	gOrders.PATCH("/:id/line/:lineid", s.OrderLineModify, s.requiresLogin)
	gOrders.DELETE("/:id/line/:lineid", s.OrderLineDelete, s.requiresLogin)
	s.e.GET("/orders", s.OrderList, s.requiresLogin)
	s.e.GET("/orders/count", s.OrderCount, s.requiresLogin, requiresAdministrator)

	return s.e.Start(fmt.Sprintf(`:%d`, s.cfg.Port))
}

package api

import (
	"crypto/sha512"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (s *Server) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := s.db.Login(username, stringToSha512(password))
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to find user")
		return echo.ErrUnauthorized
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["nombre"] = user.Name
	claims["apellidos"] = user.Surname
	claims["restaurador"] = user.IsAdmin
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to sign JWT token")
		return echo.ErrUnauthorized
	}

	log.Info().Str("username", username).Msg("User logged in")
	return c.JSON(http.StatusOK, map[string]string{"token": t})
}

// Assert that the JWT token is from a restaurador user
func requiresRestaurador(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		flag := claims["restaurador"].(bool)
		if !flag {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

func authenticated(c echo.Context) bool {
	return c.Get("user") != nil
}

func authenticatedIsRestaurador(c echo.Context) bool {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["restaurador"].(bool)
}

func authenticatedUserId(c echo.Context) uint64 {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return uint64(claims["id"].(float64))
}

// Generates SHA512 from a string
func stringToSha512(s string) string {
	h := sha512.New()
	h.Write([]byte(s))
	return fmt.Sprintf(`%x`, h.Sum(nil))
}

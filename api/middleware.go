package api

import (

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func (c echo.Context) jwt.Claims  {
			return new(JWTCustomClaims)
		},
		SigningKey:	[]byte(GetJwtKet()),

	}
	return echojwt.WithConfig(config)
}

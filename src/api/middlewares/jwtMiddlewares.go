package middlewares

import (
        "github.com/labstack/echo"
        "github.com/labstack/echo/middleware"
)

func SetJwtMiddlewares(jwtGroup *echo.Group) {
        jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
                SigningMethod: "HS512",
                SigningKey: []byte("mySecret"),
                TokenLookup: "cookie:JWTToken",
                // AuthScheme: "iLoveDogs",
        }))

}

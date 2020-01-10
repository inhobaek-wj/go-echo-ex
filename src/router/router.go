package router

import (
        "api"
        "api/middlewares"

        "github.com/labstack/echo"
)

func New() *echo.Echo {
        e := echo.New()

        // create groups.
        adminGroup := e.Group("/admin")
        cookieGroup := e.Group("/cookie")
        jwtGroup := e.Group("/jwt")

        // set all middlewares.
        middlewares.SetMainMiddlewares(e)
        middlewares.SetAdminMiddlewares(adminGroup)
        middlewares.SetCookieMiddlewares(cookieGroup)
        middlewares.SetJwtMiddlewares(jwtGroup)

        // set main group.
        api.MainGroup(e)

        // set group routes.
        api.AdminGroup(adminGroup)
        api.CookieGroup(cookieGroup)
        api.JwtGroup(jwtGroup)

        return e
}

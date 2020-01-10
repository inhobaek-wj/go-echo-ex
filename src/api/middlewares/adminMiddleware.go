package middlewares

import (
        "github.com/labstack/echo"
        "github.com/labstack/echo/middleware"
)

func SetAdminMiddlewares(adminGroup *echo.Group) {

        // this logs the server interaction.
        // g := e.Group("/admin", middleware.Logger())
        // g.Use(middleware.Logger())
        adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
                Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
        }))

        // Authentication.
        adminGroup.Use(middleware.BasicAuth(
                func(username, password string, c echo.Context) (bool, error) {
                        // check in the DB.
                        if username == "inho" && password == "1234" {
                                return true, nil
                        }
                        return false, nil
                }))
}

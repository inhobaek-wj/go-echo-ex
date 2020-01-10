package middlewares

import (
        "log"
        "strings"
        "net/http"

        "github.com/labstack/echo"
        "github.com/labstack/echo/middleware"
)

func SetCookieMiddlewares(cookieGroup *echo.Group) {
        cookieGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
                Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
        }))
        cookieGroup.Use(checkCookie)

}


func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
                cookie, err := c.Cookie("sessionId")

                if err != nil {
                        if strings.Contains(err.Error(),"named cookie not present") {
                                return c.String(
                                        http.StatusUnauthorized,
                                        "You don't have the right cookie, cookie1")
                        }
                        log.Println(err)
                        return err
                }

                if cookie.Value == "some_string" {
                        return next(c)
                }

                return c.String(
                        http.StatusUnauthorized, "You don't have the right cookie, cookie2")
        }
}

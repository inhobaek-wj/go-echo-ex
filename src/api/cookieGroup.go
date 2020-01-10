package api

import (
        "net/http"

        "github.com/labstack/echo"
)

func CookieGroup(cookieGroup *echo.Group) {
        cookieGroup.GET("/main", mainCookie)
}


func mainCookie(c echo.Context) error {
        return c.String(http.StatusOK, "You are on the secret cookie page")
}

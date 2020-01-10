package api

import (
        "net/http"

        "github.com/labstack/echo"
)

func AdminGroup(adminGroup *echo.Group) {
        adminGroup.GET("/main", mainAdmin)
}

func mainAdmin(c echo.Context) error {
        return c.String(http.StatusOK, "Hooray you are on the secret admin main page!")
}

package api

import (
        "log"
        "net/http"

        jwt "github.com/dgrijalva/jwt-go"
        "github.com/labstack/echo"
)

func JwtGroup(jwtGroup *echo.Group) {
        jwtGroup.GET("/main", mainJwt)
}


func mainJwt(c echo.Context) error {
        user := c.Get("user") // this is interface.
        token := user.(*jwt.Token)

        if claim, ok := token.Claims.(jwt.MapClaims); ok {
                log.Println("User name: ",claim["name"], "User Id: ", claim["jti"])
        }

        return c.String(http.StatusOK, "You are on the top secret jwt page!")
}

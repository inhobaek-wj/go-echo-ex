package handlers

import (
        "log"
        "net/http"
        "time"

        jwt "github.com/dgrijalva/jwt-go"
        "github.com/labstack/echo"
)

type JwtClaims struct {
        Name string `json:"name"`
        jwt.StandardClaims
}

func Login(c echo.Context) error {
        username := c.QueryParam("username")
        password := c.QueryParam("password")

        // check username and password on DB after hasing the password.
        if username == "inho" && password == "1234" {
                cookie := &http.Cookie{}

                // this is the same.
                // coocke := new(http.Cookie)

                cookie.Name = "sessionId"
                cookie.Value = "some_string"
                cookie.Expires = time.Now().Add(1 * time.Hour)

                c.SetCookie(cookie)

                // create jwt token.
                token, err := createJwtToken(username)
                if err != nil {
                        log.Println("Error Creating JWT token", err)
                        return c.String(http.StatusInternalServerError, "something went wrong")
                }

                jwtCookie := &http.Cookie{}

                jwtCookie.Name = "JWTToken"
                jwtCookie.Value = token
                jwtCookie.Expires = time.Now().Add(1 * time.Hour)

                c.SetCookie(jwtCookie)

                return c.JSON(http.StatusOK, map[string]string{
                        "message": "You were logged in!",
                        "token": token,
                })
        }
        return c.String(http.StatusOK, "Your username or password were wrong")
}

func createJwtToken(name string) (string, error){
        claims := JwtClaims{
                name,
                jwt.StandardClaims{
                        Id: "main_user_id",
                        ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
                },
        }

        rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
        token, err := rawToken.SignedString([]byte("mySecret"))

        if err != nil {
                return "", err
        }

        return token, nil
}

package main

import (
        "fmt"
        "net/http"
        "io/ioutil"
        "log"
        "encoding/json"
        "time"
        "strings"

        "github.com/labstack/echo"  // go get github.com/labstack/echo
        "github.com/labstack/echo/middleware" // https://echo.labstack.com/middleware

        jwt "github.com/dgrijalva/jwt-go"
)

type Cat struct {
        Name string `json:"name"`
        Type string `json:"type"`
}

type Dog struct {
        Name string `json:"name"`
        Type string `json:"type"`
}

type Hamster struct {
        Name string `json:"name"`
        Type string `json:"type"`
}

type JwtClaims struct {
        Name string `json:"name"`
        jwt.StandardClaims
}

func yallo (c echo.Context) error {
        return c.String(http.StatusOK, "yallo from the web site")
}

func getCats (c echo.Context) error {
        catName := c.QueryParam("name")
        catType := c.QueryParam("type")

        dataType := c.Param("data")

        if dataType == "string" {
                return c.String(
                        http.StatusOK,
                        fmt.Sprintf("Your cat name is %s and its type is %s\n",
                                catName,
                                catType))
        }

        if dataType == "json" {
                return c.JSON(http.StatusOK,map[string]string {
                        "name": catName,
                                "type": catType,
                        })
        }

        return c.JSON(http.StatusBadRequest, map[string]string{
                "error": "You need to let us know if you want string or json data",
        })
}

// the fastest.
func addCat(c echo.Context) error {
        cat := Cat{}

        defer c.Request().Body.Close()

        b, err := ioutil.ReadAll(c.Request().Body)
        if err != nil {
                log.Printf("Failed reading the request body")
                return c.String(http.StatusInternalServerError, "")
        }

        err = json.Unmarshal(b, &cat)
        if err != nil {
                log.Printf("Failed unmarshaling in addCat: %s", err)
                return c.String(http.StatusInternalServerError, "")
        }

        log.Printf("This is your cat: %#v", cat)

        return c.String(http.StatusOK, "We got your cat!")
}

func addDog(c echo.Context) error {
        dog := Dog{}

        defer c.Request().Body.Close()

        err := json.NewDecoder(c.Request().Body).Decode(&dog)
        if err != nil {
                log.Printf("Failed processing addDog request: %s", err)
                return c.String(http.StatusInternalServerError, "")
        }

        log.Printf("This is your dog: %#v", dog)
        return c.String(http.StatusOK, "We got your dog!")
}

// the slowest.
// it belongs to echo.
func addHamster(c echo.Context) error {
        hamster := Hamster{}

        err := c.Bind(&hamster)
        if err != nil {
                log.Printf("Failed processing addHamster request: %s", err)
                return echo.NewHTTPError(http.StatusInternalServerError)
        }

        log.Printf("This is your hamster: %#v", hamster)
        return c.String(http.StatusOK, "We got your hamster!")
}

func mainAdmin(c echo.Context) error {
        return c.String(http.StatusOK, "Hooray you are on the secret admin main page!")
}

func mainCookie(c echo.Context) error {
        return c.String(http.StatusOK, "You are on the secret cookie page")
}

func mainJwt(c echo.Context) error {
        user := c.Get("user") // this is interface.
        token := user.(*jwt.Token)

        if claim, ok := token.Claims.(jwt.MapClaims); ok {
                log.Println("User name: ",claim["name"], "User Id: ", claim["jti"])
        }

        return c.String(http.StatusOK, "You are on the top secret jwt page!")
}

func login(c echo.Context) error {
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

////////////////// custom middlewares //////////////////
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
                c.Response().Header().Set(echo.HeaderServer, "BlueBot/1.0")
                c.Response().Header().Set("notReallyHeader", "thisHasNoMeaning")
                return next(c)
        }
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

func main() {
        fmt.Println("Welocom to the server")

        e := echo.New()

        // apply custom middlewares.
        e.Use(ServerHeader)

        e.GET("/login", login)
        e.GET("/", yallo)
        e.GET("/cats/:data", getCats)

        e.POST("/cats", addCat)
        e.POST("/dogs", addDog)
        e.POST("/hamsters", addHamster)

        // add admin group.
        adminGroup := e.Group("/admin")

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

        adminGroup.GET("/main", mainAdmin)

        // add cookie group.
        cookieGroup := e.Group("/cookie")

        cookieGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
                Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
        }))
        cookieGroup.Use(checkCookie)

        cookieGroup.GET("/main", mainCookie)

        // add JWT group.
        jwtGroup := e.Group("/jwt")
        jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
                SigningMethod: "HS512",
                SigningKey: []byte("mySecret"),
        }))
        jwtGroup.GET("/main", mainJwt)

        e.Start(":8080")
}

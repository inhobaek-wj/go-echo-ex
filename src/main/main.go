package main

import (
        "fmt"
        "net/http"
        "io/ioutil"
        "log"
        "encoding/json"

        "github.com/labstack/echo"  // go get github.com/labstack/echo
        "github.com/labstack/echo/middleware" // https://echo.labstack.com/middleware
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

func main() {
        fmt.Println("Welocom to the server")

        e := echo.New()

        g := e.Group("/admin")

        // this logs the server interaction.
        // g := e.Group("/admin", middleware.Logger())
        // g.Use(middleware.Logger())
        g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
                Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${letancy}` + "\n",
        }))

	// Authentication.
	g.Use(middleware.BasicAuth(
		func(username, password string, c echo.Context) (bool, error) {
		// check in the DB.
		if username == "inho" && password == "1234" {
			return true, nil
		}
		return false, nil
	}))

        g.GET("/main", mainAdmin)

        e.GET("/", yallo)
        e.GET("/cats/:data", getCats)

        e.POST("/cats", addCat)
        e.POST("/dogs", addDog)
        e.POST("/hamsters", addHamster)

        e.Start(":8080")
}

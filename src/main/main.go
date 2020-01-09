package main

import (
        "fmt"
        "net/http"

        "github.com/labstack/echo"
)

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
                        fmt.Sprintf("Your cat name is %s and its type is %s\n", catName, catType))
        }

        if dataType == "json" {
                return c.JSON(http.StatusOK, map[string]string {
                        "name": catName,
                                "type": catType,
                        })
        }

        return c.JSON(http.StatusBadRequest, map[string]string{
                "error": "You need to let us know if you want string or json data",
        })
}

func main() {
        fmt.Println("Welocom to the server")

        e := echo.New()

        e.GET("/", yallo)
        e.GET("/cats/:data", getCats)

        e.Start(":8080")
}

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

        return c.String(
                http.StatusOK,
                fmt.Sprintf("Your cat name is %s and its type is %s\n", catName, catType))
}

func main() {
        fmt.Println("Welocom to the server")

        e := echo.New()

        e.GET("/", yallo)
        e.GET("/cats", getCats)
        e.Start(":8080")
}

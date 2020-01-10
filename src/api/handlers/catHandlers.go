package handlers

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"

        "github.com/labstack/echo"
)

type Cat struct {
        Name string `json:"name"`
        Type string `json:"type"`
}

func GetCats (c echo.Context) error {
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
func AddCat(c echo.Context) error {
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

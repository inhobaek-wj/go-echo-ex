package handlers

import (
        "log"
        "net/http"

        "github.com/labstack/echo"
)

type Hamster struct {
        Name string `json:"name"`
        Type string `json:"type"`
}

// the slowest.
// it belongs to echo.
func AddHamster(c echo.Context) error {
        hamster := Hamster{}

        err := c.Bind(&hamster)
        if err != nil {
                log.Printf("Failed processing addHamster request: %s", err)
                return echo.NewHTTPError(http.StatusInternalServerError)
        }

        log.Printf("This is your hamster: %#v", hamster)
        return c.String(http.StatusOK, "We got your hamster!")
}

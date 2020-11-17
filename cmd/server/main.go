package main

import (
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Static("/data", "data")
	e.Logger.Fatal(e.Start(":8080"))
}

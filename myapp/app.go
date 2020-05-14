package myapp

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type User struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewHttpHandler() *echo.Echo {
	e := echo.New()
	e.GET("/", indexHandler)
	e.GET("/bar", barHandler)
	e.GET("/foo", fooHandler)
	return e
}

func indexHandler(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello world")
}
func fooHandler(ctx echo.Context) error {
	user := new(User)
	user.CreatedAt = time.Now()
	if err := ctx.Bind(user); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, user)
}
func barHandler(ctx echo.Context) error {
	name := ctx.FormValue("name")
	if name == "" {
		name = "World"
	}
	return ctx.String(http.StatusOK, fmt.Sprintf("Hello %s", name))
}

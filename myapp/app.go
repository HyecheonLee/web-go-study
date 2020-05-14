package myapp

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
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
	e.Static("", "public/index.html")
	e.GET("/bar", barHandler)
	e.GET("/foo", fooHandler)
	e.POST("/upload", uploadsHandler)
	return e
}
func uploadsHandler(ctx echo.Context) error {
	header, err := ctx.FormFile("upload_file")
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	uploadFile, err := header.Open()
	if err != nil {
		return err
	}
	defer uploadFile.Close()

	dirname := "./uploads"
	os.MkdirAll(dirname, 0777)
	filepath := fmt.Sprintf("%s/%s", dirname, header.Filename)
	file, err := os.Create(filepath)
	defer file.Close()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	if _, err = io.Copy(file, uploadFile); err != nil {
		return err
	}
	return ctx.String(http.StatusOK, filepath)
}

func indexHandler(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello world")
}
func fooHandler(ctx echo.Context) error {
	user := new(User)
	err := json.NewDecoder(ctx.Request().Body).Decode(user)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	user.CreatedAt = time.Now()
	return ctx.JSON(http.StatusCreated, user)
}
func barHandler(ctx echo.Context) error {
	name := ctx.FormValue("name")
	if name == "" {
		name = "World"
	}
	return ctx.String(http.StatusOK, fmt.Sprintf("Hello %s", name))
}

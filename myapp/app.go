package myapp

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	Id        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var userDb map[int]*User
var lastID = 0

func NewHttpHandler() *echo.Echo {
	userDb = make(map[int]*User)
	e := echo.New()
	e.Static("/file", "public")
	e.GET("/", indexHandler)
	e.GET("/bar", barHandler)
	e.GET("/foo", fooHandler)
	e.POST("/upload", uploadsHandler)
	e.GET("/users", usersHandler)
	e.POST("/users", createUsersHandler)
	e.GET("/users/:userId", getUserInfoHandler)

	return e
}

func createUsersHandler(context echo.Context) error {
	user := new(User)
	err := context.Bind(user)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}
	lastID++
	user.Id = lastID
	user.CreatedAt = time.Now()
	userDb[user.Id] = user

	return context.JSON(http.StatusCreated, user)
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
func usersHandler(ctx echo.Context) error {
	userId := ctx.Param("userId")
	return ctx.String(http.StatusOK, fmt.Sprintf("Get UserInfo by /users/%s", userId))
}
func getUserInfoHandler(ctx echo.Context) error {
	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	user, ok := userDb[id]
	if ok {
		return ctx.JSON(http.StatusOK, user)
	} else {
		return ctx.JSON(http.StatusOK, "No User Id")
	}
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

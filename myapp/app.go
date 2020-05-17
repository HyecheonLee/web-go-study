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
	e.DELETE("/users/:userId", deleteUserHandler)
	e.PUT("/users/:userId", updateUserHandler)

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
func usersHandler(ctx echo.Context) error {
	var users []*User
	for _, user := range userDb {
		users = append(users, user)
	}
	if len(users) == 0 {
		return ctx.JSON(http.StatusOK, []*User{})
	}
	return ctx.JSON(http.StatusOK, users)
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
func deleteUserHandler(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	_, ok := userDb[userId]
	if !ok {
		return findNotUserId(ctx, userId)
	}
	delete(userDb, userId)
	return ctx.JSON(http.StatusOK, fmt.Sprint("Delete User Id: ", userId))
}
func updateUserHandler(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	savedUser, ok := userDb[userId]
	if !ok {
		return findNotUserId(ctx, userId)
	}

	updateUser := new(User)
	if err = ctx.Bind(updateUser); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	if updateUser.FirstName != "" {
		savedUser.FirstName = updateUser.FirstName
	}
	if updateUser.LastName != "" {
		savedUser.LastName = updateUser.LastName
	}
	if updateUser.Email != "" {
		savedUser.Email = updateUser.Email
	}
	return ctx.JSON(http.StatusOK, savedUser)
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

func findNotUserId(ctx echo.Context, userId int) error {
	return ctx.JSON(http.StatusOK, fmt.Sprint("No User Id: ", userId))
}

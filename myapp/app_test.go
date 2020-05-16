package myapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIndexPathHandler(t *testing.T) {
	e := echo.New()
	defer e.Close()

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	indexHandler(c)
	// Assertions

	assert.Equal(t, http.StatusOK, rec.Code)
	data, _ := ioutil.ReadAll(rec.Body)
	assert.Equal(t, "Hello world", string(data))
}

func TestBarPathHandler_WithoutName(t *testing.T) {
	e := echo.New()
	defer e.Close()

	req := httptest.NewRequest(http.MethodPost, "/bar", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	barHandler(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	data, _ := ioutil.ReadAll(rec.Body)
	assert.Equal(t, "Hello World", string(data))
}

func TestBarPathHandler_WithName(t *testing.T) {
	e := echo.New()
	defer e.Close()

	q := make(url.Values)
	q.Set("name", "name")

	req := httptest.NewRequest(http.MethodPost, "/bar/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	barHandler(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	data, _ := ioutil.ReadAll(rec.Body)
	assert.Equal(t, "Hello name", string(data))
}
func TestFooHandler_WithoutJson(t *testing.T) {

	e := echo.New()
	defer e.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	c := e.NewContext(req, res)
	fooHandler(c)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}
func TestFooHandler_WithJson(t *testing.T) {

	e := echo.New()
	defer e.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/foo",
		strings.NewReader(`{"first_name":"hyecheon","last_name":"lee","email":"rainbow880616@gmail.com"}`))
	c := e.NewContext(req, res)
	fooHandler(c)

	assert.Equal(t, http.StatusCreated, res.Code)

	user := new(User)
	err := json.NewDecoder(res.Body).Decode(user)
	assert.Nil(t, err)
	assert.Equal(t, "hyecheon", user.FirstName)
	assert.Equal(t, "lee", user.LastName)

}

func TestUploadsHandler(t *testing.T) {

	path := `C:\Users\hyecheon\Downloads\증명-.jpg`

	file, _ := os.Open(path)
	defer file.Close()

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	multi, err := writer.CreateFormFile("upload_file", filepath.Base(path))
	assert.NoError(t, err)
	io.Copy(multi, file)
	writer.Close()

	e := echo.New()
	defer e.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/upload", buf)
	req.Header.Set("Content-type", writer.FormDataContentType())
	c := e.NewContext(req, res)
	err = uploadsHandler(c)
	assert.NoError(t, err)

	uploadFilePath := "./uploads/" + filepath.Base(path)
	_, err = os.Stat(uploadFilePath)
	assert.NoError(t, err)

	uploadFile, _ := os.Open(uploadFilePath)
	originFile, _ := os.Open(path)
	defer uploadFile.Close()
	defer originFile.Close()
	uploadData := []byte{}
	originData := []byte{}
	uploadFile.Read(uploadData)
	originFile.Read(originData)

	assert.Equal(t, originData, uploadData)

}

func TestUsers(t *testing.T) {

	e := echo.New()
	defer e.Close()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := usersHandler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	data, _ := ioutil.ReadAll(rec.Body)
	assert.Contains(t, string(data), "Get UserInfo")
}
func TestGetUsers(t *testing.T) {

	e := echo.New()
	defer e.Close()

	req := httptest.NewRequest(http.MethodGet, "/users/89", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:userId")
	c.SetParamNames("userId")
	c.SetParamValues("89")

	err := usersHandler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	data, _ := ioutil.ReadAll(rec.Body)
	assert.Contains(t, string(data), "/users/89")
}

func TestCreateUsers(t *testing.T) {
	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name": "hyecheon", "last_name": "hclee", "email":"rainbow0616@naver.com"}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	savedUser := User{}
	json.NewDecoder(resp.Body).Decode(savedUser)

	resp, err = http.Get(fmt.Sprint(ts.URL, "/users/", savedUser.Id))

	assert.NoError(t, err)
	findUser := User{}
	json.NewDecoder(resp.Body).Decode(findUser)
	assert.Equal(t, savedUser.Id, findUser.Id)
}

package myapp

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestIndexPathHandler(t *testing.T) {
	e := echo.New()
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

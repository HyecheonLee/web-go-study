package main

import (
	"net/http"
	"web-go-study/myapp"
)

func main() {
	e := myapp.NewHttpHandler()
	e.Logger.Fatal(http.ListenAndServe(":3000", e))
}

package main

import "web-go-study/myapp"

func main() {
	e := myapp.NewHttpHandler()

	e.Logger.Fatal(e.Start(":3000"))
}

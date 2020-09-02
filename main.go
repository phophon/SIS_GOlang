package main

import (
	"app"
	"fmt"
	"time"
)

func main() {
	app.Init()
	fmt.Println("=========Start", time.Now().String())
	StartServer()
}

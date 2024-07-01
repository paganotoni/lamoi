package main

import (
	"fmt"
	"net/http"

	"lamoi/internal"

	_ "github.com/leapkit/leapkit/core/tools/envload" // Load environment variables
	_ "github.com/mattn/go-sqlite3"                   // sqlite3 driver
)

func main() {
	s := internal.New()
	fmt.Println("Server started at", s.Addr())
	err := http.ListenAndServe(s.Addr(), s.Handler())
	if err != nil {
		fmt.Println("[error] starting app:", err)
	}
}

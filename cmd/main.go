package main

import (
	"fmt"
	"net/http"

	httpd "github.com/shanebailey05/future_backend_homework/http"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("[%v]\n", err)
	}
}

func run() error {
	service, err := httpd.New()
	if err != nil {
		return err
	}
	return http.ListenAndServe(":10000", httpd.Routes(service))
}

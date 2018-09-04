package main

import (
	"github.com/KonstantinGrig/goapisql/config"
	"github.com/KonstantinGrig/goapisql/handlers"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		log.Println("Start with custom config file: " + args[0])
	} else {
		log.Println("Start default config file: config.json!")
	}
	os.Setenv("GOAPISQL_ENV", "test")
	config.Init()
	port := config.GetString("server-port")

	fasthttp.ListenAndServe(":"+port, handlers.SqlHandler)
}

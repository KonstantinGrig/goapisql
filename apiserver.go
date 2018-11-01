package main

import (
	"github.com/KonstantinGrig/goapisql/config"
	"github.com/KonstantinGrig/goapisql/handlers"
	"github.com/valyala/fasthttp"
	"os"
)

const defaultPort = "9595"

func main() {
	config.InitConfigFile("./config.json")
	port := os.Getenv("GOAPISQL_PORT")
	if port == "" {
		port = defaultPort
	}

	err := fasthttp.ListenAndServe(":"+port, handlers.FastHttpHandler)
	if err != nil {
		panic(err)
	}
}

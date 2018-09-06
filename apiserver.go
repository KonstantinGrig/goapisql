package main

import (
	"github.com/KonstantinGrig/goapisql/config"
	"github.com/KonstantinGrig/goapisql/handlers"
	"github.com/valyala/fasthttp"
)

func main() {
	config.Init()
	port := config.GetString("server-port")

	fasthttp.ListenAndServe(":"+port, handlers.SqlHandler)
}

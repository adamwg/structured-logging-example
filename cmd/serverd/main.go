package main

import (
	"net/http"

	"github.com/adamwg/structured-logging-example/internal/server"
	"github.com/aybabtme/log"
)

func main() {
	l := log.KV("app", "serverd")

	srv := server.New(l)
	l.Info("starting server")
	err := http.ListenAndServe(":80", srv)
	l.Err(err).Info("server exited")
}

package main

import "github.com/fire9900/auth/internal/app"

func main() {
	app.LoggerRun()
	go app.Run()
	go app.StartGRPCServer()
}

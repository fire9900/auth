package main

import "github.com/fire9900/auth/internal/app"

func main() {
	app.LoggerRun()
	go app.Run()
	go app.StartGRPCServer()
	select {}
	//password, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	//if err != nil {
	//	return
	//}
	//fmt.Println(string(password))
}

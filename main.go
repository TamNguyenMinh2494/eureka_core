package main

import (
	"fmt"
	"main/core"
	"main/core/routers"
)

func main() {
	fmt.Println("Eureka Core...")
	fmt.Println("8==> Start")
	server := core.Server{}
	server.LoadConfig("config.json")
	server.Create()

	server.Routers = []core.Router{
		&routers.AuthRouter{Name: "v1/auth"},
		&routers.UserRouter{Name: "v1/users"},
		&routers.CourseRouter{Name: "v1/courses"},
		&routers.ExamRouter{Name: "v1/exams"},
	}

	server.ConnectRouters()

	server.Start(server.Config.Address)
	defer server.Dispose()
}

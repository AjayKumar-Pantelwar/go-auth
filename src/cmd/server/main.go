package main

import (
	"go-authentication/src/internal/adaptors/persistance"
	userhandler "go-authentication/src/internal/interfaces/input/api/rest/handler"
	"go-authentication/src/internal/interfaces/input/api/rest/routes"
	user "go-authentication/src/internal/usecase"
	"log"
	"net/http"
)

func main() {
	database, err := persistance.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	userRepo := persistance.NewUserRepo(database)
	userService := user.NewUserService(userRepo)
	userHandler := userhandler.NewUserHandler(userService)

	router := routes.InitRoutes(&userHandler)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}

package main

import (
	"go-authentication/src/db"
	"go-authentication/src/internal/adaptors/persistance"
	"go-authentication/src/internal/interfaces/input/api/rest/routes"
	user "go-authentication/src/internal/usecase"
	"log"
	"net/http"
)

func main() {

	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	userrepo := persistance.NewUserRepo(database)

	userService := user.NewUserService(userrepo)

	router := routes.InitRoutes()

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}

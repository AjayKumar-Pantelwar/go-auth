package main

import (
	"go-authentication/src/internal/adaptors/persistance"
	userhandler "go-authentication/src/internal/interfaces/input/api/rest/handler"
	"go-authentication/src/internal/interfaces/input/api/rest/routes"
	user "go-authentication/src/internal/usecase"
	"go-authentication/src/pkg/migrate"
	"log"
	"net/http"
	"os"
)

func main() {
	database, err := persistance.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get cwd: %v", err)
	}

	migrate := migrate.NewMigrate(
		database.GetDB(),
		cwd+"/src/migrations",
	)

	err = migrate.RunMigrations()
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := persistance.NewUserRepo(database)
	userService := user.NewUserService(userRepo)
	userHandler := userhandler.NewUserHandler(userService)

	router := routes.InitRoutes(&userHandler)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}

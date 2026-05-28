package main

import (
	"errors"
	"os"
	"tarefas/internal/application/user"
	"tarefas/internal/infrastructure/db"

	"github.com/go-fuego/fuego"
	"github.com/joho/godotenv"
)

func main() {
	server := fuego.NewServer()

	godotenv.Load()

	databaseConfig := db.Config{
		DSN: os.Getenv("DATABASE_DSN"),
	}

	database, err := db.NewPostgresDatabase(databaseConfig)
	if err != nil {
		err = errors.Join(errors.New("Failed to initialize database"), err)
		panic(err)
	}

	userRepo := db.NewUserRepo(database)

	userService := user.NewService(userRepo)

	fuego.Get(server, "/ping", func(c fuego.ContextNoBody) (string, error) {
		return "pong", nil
	})

	users := fuego.Group(server, "/users")

	fuego.Post(users, "/", func(c fuego.ContextNoBody) (string, error) {
		return "User created", nil
	})

	fuego.Get(users, "/:id", func(c fuego.ContextNoBody) (string, error) {
		return "user", nil
	})

	fuego.Put(users, "/:id", func(c fuego.ContextNoBody) (string, error) {
		return "User updated", nil
	})

	fuego.Delete(users, "/:id", func(c fuego.ContextNoBody) (string, error) {
		return "User deleted", nil
	})

	tasks := fuego.Group(server, "/tasks")

	fuego.Post(tasks, "/", func(c fuego.ContextNoBody) (string, error) {
		return "Task created", nil
	})

	fuego.Get(tasks, "/:id", func(c fuego.ContextNoBody) (string, error) {
		return "task", nil
	})

	fuego.Get(tasks, "", func(c fuego.ContextNoBody) (string, error) {
		return "tasks assigned to user", nil
	})

	fuego.Put(tasks, "/:id", func(c fuego.ContextNoBody) (string, error) {
		return "Task updated", nil
	})

	fuego.Delete(tasks, "/:id", func(c fuego.ContextNoBody) (string, error) {
		return "Task deleted", nil
	})

	auth := fuego.Group(server, "/auth")

	fuego.Post(auth, "/login", func(c fuego.ContextNoBody) (string, error) {
		return "User logged in", nil
	})

	fuego.Post(auth, "/logout", func(c fuego.ContextNoBody) (string, error) {
		return "User logged out", nil
	})

	server.Run()
}

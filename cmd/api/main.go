package main

import (
	"errors"
	"os"
	"tarefas/cmd/api/handler"
	"tarefas/internal/application/task"
	"tarefas/internal/application/user"
	"tarefas/internal/infrastructure/db"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
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
	userHandler := handler.NewUserHandler(userService)

	taskRepo := db.NewTaskRepo(database)
	taskService := task.NewService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	fuego.Get(server, "/ping", func(c fuego.ContextNoBody) (string, error) {
		return "pong", nil
	})

	users := fuego.Group(server, "/users")
	fuego.Post(users, "/", userHandler.CreateUser)
	fuego.Get(users, "/{id}", userHandler.GetUserById)
	fuego.Put(users, "/:id", userHandler.UpdateUser)
	fuego.Delete(users, "/:id", userHandler.DeleteUser)

	tasks := fuego.Group(server, "/tasks")

	fuego.Post(tasks, "/", taskHandler.CreateTask)
	fuego.Get(tasks, "/:id", taskHandler.GetTaskById)
	fuego.Get(
		tasks,
		"/",
		taskHandler.QueryTasks,
		option.Query("assignedTo", "Filter by assigned user", param.Nullable()),
		option.Query("createdBy", "Filter by creator", param.Nullable()),
	)
	fuego.Put(tasks, "/", taskHandler.UpdateTask)
	fuego.Delete(tasks, "/:id", taskHandler.DeleteTask)

	auth := fuego.Group(server, "/auth")

	fuego.Post(auth, "/login", func(c fuego.ContextNoBody) (string, error) {
		return "User logged in", nil
	})

	fuego.Post(auth, "/logout", func(c fuego.ContextNoBody) (string, error) {
		return "User logged out", nil
	})

	server.Run()
}

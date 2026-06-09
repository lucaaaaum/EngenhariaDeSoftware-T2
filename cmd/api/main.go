package main

import (
	"fmt"
	"os"

	"tarefas/cmd/api/handler"
	"tarefas/cmd/api/middleware"
	"tarefas/internal/application/auth"
	"tarefas/internal/application/task"
	"tarefas/internal/application/user"
	"tarefas/internal/application/webhook"
	"tarefas/internal/infrastructure/db"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	database, err := db.NewPostgresDatabase(db.Config{
		DSN: os.Getenv("DATABASE_DSN"),
	})
	if err != nil {
		panic(fmt.Errorf("database init: %w", err))
	}

	userRepo := db.NewUserRepo(database)
	taskRepo := db.NewTaskRepo(database)

	webhookService := webhook.NewService()
	userService := user.NewService(userRepo)
	taskService := task.NewService(taskRepo, webhookService)
	authService := auth.NewService(userRepo)

	userHandler := handler.NewUserHandler(userService)
	taskHandler := handler.NewTaskHandler(taskService)
	authHandler := handler.NewAuthHandler(authService)

	server := fuego.NewServer(fuego.WithAddr(":9999"))
	server.OpenAPI.Description().Servers = openapi3.Servers{{URL: "http://localhost:9999"}}

	// Adiciona o esquema de autenticação Bearer JWT ao Swagger
	server.OpenAPI.Description().Components.SecuritySchemes = openapi3.SecuritySchemes{
		"bearerAuth": &openapi3.SecuritySchemeRef{
			Value: openapi3.NewSecurityScheme().
				WithType("http").
				WithScheme("bearer").
				WithBearerFormat("JWT"),
		},
	}
	server.OpenAPI.Description().Security = openapi3.SecurityRequirements{
		{"bearerAuth": {}},
	}

	// rotas públicas
	authGroup := fuego.Group(server, "/auth")
	fuego.Post(authGroup, "/login", authHandler.Login)
	fuego.Post(authGroup, "/logout", authHandler.Logout)

	// POST /users é público (cadastro não precisa de token)
	publicUsers := fuego.Group(server, "/users")
	fuego.Post(publicUsers, "/", userHandler.CreateUser)

	// rotas protegidas por JWT
	protected := fuego.Group(server, "", fuego.OptionMiddleware(middleware.AuthMiddleware))

	users := fuego.Group(protected, "/users")
	fuego.Get(users, "/{id}", userHandler.GetUserById)
	fuego.Put(users, "/{id}", userHandler.UpdateUser)
	fuego.Delete(users, "/{id}", userHandler.DeleteUser)

	tasks := fuego.Group(protected, "/tasks")
	fuego.Post(tasks, "/", taskHandler.CreateTask)
	fuego.Get(tasks, "/{id}", taskHandler.GetTaskById)
	fuego.Get(
		tasks,
		"/",
		taskHandler.QueryTasks,
		option.Query("assignedTo", "filter by assigned user (UUID)", param.Nullable()),
		option.Query("createdBy", "filter by creator (UUID)", param.Nullable()),
		option.Query("status", "0=Pending, 1=InProgress, 2=Completed", param.Nullable()),
		option.Query("priority", "0=Low, 1=Medium, 2=High", param.Nullable()),
		option.Query("dueBefore", "max due date (YYYY-MM-DD)", param.Nullable()),
	)
	fuego.Put(tasks, "/{id}", taskHandler.UpdateTask)
	fuego.Delete(tasks, "/{id}", taskHandler.DeleteTask)

	if err := server.Run(); err != nil {
		panic(err)
	}
}

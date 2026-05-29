package main

import (
	"errors"
	"os"
	"tarefas/internal/application/user"
	"tarefas/internal/infrastructure/db"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
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

	fuego.Post(users, "/", func(c fuego.ContextWithBody[user.CreateUserCommand]) (*user.UserDto, error) {
		cmd, err := c.Body()
		if err != nil {
			return nil, errors.Join(errors.New("Failed to parse request body"), err)
		}
		createdUser, err := userService.CreateUser(c.Context(), cmd)
		if err != nil {
			return nil, errors.Join(errors.New("Failed to create user"), err)
		}

		return user.NewUserDto(createdUser), nil
	})

	fuego.Get(users, "/{id}", func(c fuego.ContextNoBody) (*user.UserDto, error) {
		stringId := c.PathParam("id")
		id, err := uuid.Parse(stringId)
		if err != nil {
			return nil, errors.Join(errors.New("Invalid user ID"), err)
		}
		userFound, err := userRepo.GetUserById(c.Context(), id)
		if err != nil {
			return nil, errors.Join(errors.New("Failed to get user"), err)
		}
		return user.NewUserDto(userFound), nil
	})

	fuego.Put(users, "/:id", func(c fuego.ContextWithBody[user.UpdateUserCommand]) (any, error) {
		cmd, err := c.Body()
		if err != nil {
			return nil, errors.Join(errors.New("Failed to parse request body"), err)
		}
		id, err := uuid.Parse(c.PathParam("id"))
		if err != nil {
			return nil, errors.Join(errors.New("Invalid user ID"), err)
		}
		cmd.Id = id
		err = userService.UpdateUser(c.Context(), cmd)
		if err != nil {
			return nil, errors.Join(errors.New("Failed to update user"), err)
		}
		c.SetStatus(204)
		return nil, nil
	})

	fuego.Delete(users, "/:id", func(c fuego.ContextNoBody) (any, error) {
		id, err := uuid.Parse(c.PathParam("id"))
		if err != nil {
			return nil, errors.Join(errors.New("Invalid user ID"), err)
		}
		err = userService.DeleteUser(c.Context(), id)
		if err != nil {
			return nil, errors.Join(errors.New("Failed to delete user"), err)
		}
		c.SetStatus(204)
		return nil, nil
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

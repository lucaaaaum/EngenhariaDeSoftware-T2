package user_test

import (
	"tarefas/internal/domain/user"
	"testing"
)

func TestCreateUser_Success(t *testing.T) {
	name := "username"
	email := "email@email.com"
	password := "password"

	user, err := user.NewUser(name, email, password)

	if user == nil {
		t.Fatalf("User was supposed to be created")
	}

	if err != nil {
		t.Fatalf("There should be no error here")
	}

	if user.Name != name {
		t.Fatalf("Expeceted username %q, got %q", name, user.Name)
	}

	if user.Email != email {
		t.Fatalf("Expeceted email %q, got %q", email, user.Email)
	}

	if user.PasswordHash == "" || user.PasswordHash == password {
		t.Fatalf("Password should've been hashed")
	}
}

func TestCreateUser_Fail_NameIsEmpty(t *testing.T) {
	name := ""
	email := "email@email.com"
	password := "password"

	user, err := user.NewUser(name, email, password)

	if user != nil {
		t.Fatalf("User was not supposed to have been created")
	}

	if err == nil {
		t.Fatalf("There should be an error here")
	}
}

func TestCreateUser_Fail_EmailIsEmpty(t *testing.T) {
	name := "username"
	email := ""
	password := "password"

	user, err := user.NewUser(name, email, password)

	if user != nil {
		t.Fatalf("User was not supposed to have been created")
	}

	if err == nil {
		t.Fatalf("There should be an error here")
	}
}

func TestCreateUser_Fail_PasswordIsEmpty(t *testing.T) {
	name := "username"
	email := "email@email.com"
	password := ""

	user, err := user.NewUser(name, email, password)

	if user != nil {
		t.Fatalf("User was not supposed to have been created")
	}

	if err == nil {
		t.Fatalf("There should be an error here")
	}
}

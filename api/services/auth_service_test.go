package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myanimeapi/api/auth"
	"myanimeapi/api/models"
)

func TestAuthService_NewAuthService(t *testing.T) {
	repo := &mockAuthRepo{}
	svc := NewAuthService(repo)
	assert.NotNil(t, svc)
}

func TestAuthService_RegisterUser_Success(t *testing.T) {
	repo := &mockAuthRepo{}
	svc := NewAuthService(repo)
	ctx := context.Background()

	user := &models.User{Username: "testuser", Email: "test@example.com", Password: "plainpass"}

	repo.On("GetUserByUsername", ctx, "testuser").Return(nil, errors.New("not found"))
	repo.On("GetUserByEmail", ctx, "test@example.com").Return(nil, errors.New("not found"))
	repo.On("CreateUser", ctx, user).Return(nil)

	err := svc.RegisterUser(ctx, user)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAuthService_RegisterUser_UsernameConflict(t *testing.T) {
	repo := &mockAuthRepo{}
	svc := NewAuthService(repo)
	ctx := context.Background()

	existing := &models.User{Username: "testuser"}
	user := &models.User{Username: "testuser", Email: "new@example.com", Password: "pass"}

	repo.On("GetUserByUsername", ctx, "testuser").Return(existing, nil)

	err := svc.RegisterUser(ctx, user)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestAuthService_RegisterUser_EmailConflict(t *testing.T) {
	repo := &mockAuthRepo{}
	svc := NewAuthService(repo)
	ctx := context.Background()

	existing := &models.User{Email: "taken@example.com"}
	user := &models.User{Username: "newuser", Email: "taken@example.com", Password: "pass"}

	repo.On("GetUserByUsername", ctx, "newuser").Return(nil, errors.New("not found"))
	repo.On("GetUserByEmail", ctx, "taken@example.com").Return(existing, nil)

	err := svc.RegisterUser(ctx, user)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestAuthService_RegisterUser_CreateError(t *testing.T) {
	repo := &mockAuthRepo{}
	svc := NewAuthService(repo)
	ctx := context.Background()

	user := &models.User{Username: "newuser", Email: "new@example.com", Password: "pass"}

	repo.On("GetUserByUsername", ctx, "newuser").Return(nil, errors.New("not found"))
	repo.On("GetUserByEmail", ctx, "new@example.com").Return(nil, errors.New("not found"))
	repo.On("CreateUser", ctx, user).Return(errors.New("db error"))

	err := svc.RegisterUser(ctx, user)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestAuthService_AuthenticateUser_UserNotFound(t *testing.T) {
	repo := &mockAuthRepo{}
	svc := NewAuthService(repo)
	ctx := context.Background()

	creds := &models.UserCredentials{Username: "ghost", Password: "pass"}
	repo.On("GetUserByUsername", ctx, "ghost").Return(nil, errors.New("not found"))

	token, err := svc.AuthenticateUser(ctx, creds)
	assert.Empty(t, token)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestAuthService_AuthenticateUser_WrongPassword(t *testing.T) {
	repo := &mockAuthRepo{}
	svc := NewAuthService(repo)
	ctx := context.Background()

	hashed, _ := auth.HashPassword("correctpass")
	user := &models.User{Username: "alice", Password: hashed}
	user.ID = 1

	creds := &models.UserCredentials{Username: "alice", Password: "wrongpass"}
	repo.On("GetUserByUsername", ctx, "alice").Return(user, nil)

	token, err := svc.AuthenticateUser(ctx, creds)
	assert.Empty(t, token)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestAuthService_AuthenticateUser_Success(t *testing.T) {
	repo := &mockAuthRepo{}
	svc := NewAuthService(repo)
	ctx := context.Background()

	hashed, _ := auth.HashPassword("correctpass")
	user := &models.User{Username: "alice", Password: hashed, IsAdmin: false, Role: "user"}
	user.ID = 1

	creds := &models.UserCredentials{Username: "alice", Password: "correctpass"}
	repo.On("GetUserByUsername", ctx, "alice").Return(user, nil)

	token, err := svc.AuthenticateUser(ctx, creds)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	repo.AssertExpectations(t)
}

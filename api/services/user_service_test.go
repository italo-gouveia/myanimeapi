package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"myanimeapi/api/auth"
	"myanimeapi/api/models"
)

func newUserSvc(repo *mockUserRepo) *UserService {
	return &UserService{
		userRepo: repo,
		logger:   newTestLogger(),
	}
}

func TestUserService_GetAllUsers_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	u1 := &models.User{Username: "alice", Email: "alice@test.com"}
	u1.ID = 1
	u2 := &models.User{Username: "bob", Email: "bob@test.com"}
	u2.ID = 2
	repo.On("GetAll", ctx, 1, 10).Return([]interface{}{u1, u2}, int64(2), nil)

	result, total, err := svc.GetAllUsers(ctx, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

func TestUserService_GetAllUsers_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	repo.On("GetAll", ctx, 1, 10).Return(nil, int64(0), errors.New("db error"))

	_, _, err := svc.GetAllUsers(ctx, 1, 10)
	require.Error(t, err)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	repo.On("GetByUsername", ctx, "newuser").Return(nil, errors.New("not found"))
	repo.On("GetByEmail", ctx, "new@test.com").Return(nil, errors.New("not found"))
	repo.On("Create", ctx, mock.Anything).Return(nil)

	user := &models.User{Username: "newuser", Email: "new@test.com", Password: "password123"}
	err := svc.CreateUser(ctx, user)
	require.NoError(t, err)
}

func TestUserService_CreateUser_UsernameConflict(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	existing := &models.User{Username: "newuser"}
	repo.On("GetByUsername", ctx, "newuser").Return(existing, nil)

	user := &models.User{Username: "newuser", Email: "new@test.com", Password: "password123"}
	err := svc.CreateUser(ctx, user)
	require.Error(t, err)
}

func TestUserService_CreateUser_EmailConflict(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	repo.On("GetByUsername", ctx, "newuser").Return(nil, errors.New("not found"))
	existing := &models.User{Email: "new@test.com"}
	repo.On("GetByEmail", ctx, "new@test.com").Return(existing, nil)

	user := &models.User{Username: "newuser", Email: "new@test.com", Password: "password123"}
	err := svc.CreateUser(ctx, user)
	require.Error(t, err)
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	user := &models.User{Username: "testuser"}
	user.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(user, nil)
	repo.On("Delete", ctx, uint(1)).Return(nil)

	err := svc.DeleteUser(ctx, 1)
	require.NoError(t, err)
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	repo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	err := svc.DeleteUser(ctx, 1)
	require.Error(t, err)
}

func TestUserService_GetByID_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	user := &models.User{Username: "testuser"}
	user.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(user, nil)

	result, err := svc.GetByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "testuser", result.Username)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	repo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	_, err := svc.GetByID(ctx, 1)
	require.Error(t, err)
}

func TestUserService_GetByEmail_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	user := &models.User{Email: "test@test.com"}
	repo.On("GetByEmail", ctx, "test@test.com").Return(user, nil)

	result, err := svc.GetByEmail(ctx, "test@test.com")
	require.NoError(t, err)
	assert.Equal(t, "test@test.com", result.Email)
}

func TestUserService_GetByUsername_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	user := &models.User{Username: "testuser"}
	repo.On("GetByUsername", ctx, "testuser").Return(user, nil)

	result, err := svc.GetByUsername(ctx, "testuser")
	require.NoError(t, err)
	assert.Equal(t, "testuser", result.Username)
}

func TestUserService_ValidateUser_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	hash, err := auth.HashPassword("password123")
	require.NoError(t, err)

	user := &models.User{Username: "testuser", Password: hash}
	user.ID = 1
	repo.On("GetByUsername", ctx, "testuser").Return(user, nil)

	result, err := svc.ValidateUser(ctx, "testuser", "password123")
	require.NoError(t, err)
	assert.Equal(t, "testuser", result.Username)
}

func TestUserService_ValidateUser_UserNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	repo.On("GetByUsername", ctx, "unknown").Return(nil, errors.New("not found"))

	_, err := svc.ValidateUser(ctx, "unknown", "password123")
	require.Error(t, err)
}

func TestUserService_ValidateUser_WrongPassword(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	hash, err := auth.HashPassword("password123")
	require.NoError(t, err)

	user := &models.User{Username: "testuser", Password: hash}
	user.ID = 1
	repo.On("GetByUsername", ctx, "testuser").Return(user, nil)

	_, err = svc.ValidateUser(ctx, "testuser", "wrongpassword")
	require.Error(t, err)
}

func TestUserService_UpdateUser_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	existing := &models.User{Username: "testuser", Email: "test@test.com"}
	existing.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	repo.On("Update", ctx, mock.Anything).Return(nil)

	user := &models.User{Username: "testuser", Email: "test@test.com"}
	user.ID = 1
	err := svc.UpdateUser(ctx, user)
	require.NoError(t, err)
}

func TestUserService_UpdateUser_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	repo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	user := &models.User{Username: "testuser"}
	user.ID = 1
	err := svc.UpdateUser(ctx, user)
	require.Error(t, err)
}

func TestUserService_ChangePassword_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	hash, err := auth.HashPassword("currentpass")
	require.NoError(t, err)

	user := &models.User{Username: "testuser", Password: hash}
	user.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(user, nil)
	repo.On("Update", ctx, mock.Anything).Return(nil)

	err = svc.ChangePassword(ctx, 1, "currentpass", "newpassword")
	require.NoError(t, err)
}

func TestUserService_ChangePassword_WrongCurrent(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	hash, err := auth.HashPassword("currentpass")
	require.NoError(t, err)

	user := &models.User{Username: "testuser", Password: hash}
	user.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(user, nil)

	err = svc.ChangePassword(ctx, 1, "wrongpass", "newpassword")
	require.Error(t, err)
}

func TestUserService_ChangePassword_UserNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	repo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	err := svc.ChangePassword(ctx, 1, "currentpass", "newpassword")
	require.Error(t, err)
}

func TestUserService_DeactivateAccount_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	hash, err := auth.HashPassword("mypassword")
	require.NoError(t, err)

	user := &models.User{Username: "testuser", Password: hash, IsActive: true}
	user.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(user, nil)
	repo.On("Update", ctx, mock.Anything).Return(nil)

	err = svc.DeactivateAccount(ctx, 1, "mypassword")
	require.NoError(t, err)
}

func TestUserService_DeactivateAccount_UserNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	repo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	err := svc.DeactivateAccount(ctx, 1, "mypassword")
	require.Error(t, err)
}

func TestUserService_DeactivateAccount_WrongPassword(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	hash, err := auth.HashPassword("mypassword")
	require.NoError(t, err)

	user := &models.User{Username: "testuser", Password: hash, IsActive: true}
	user.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(user, nil)

	err = svc.DeactivateAccount(ctx, 1, "wrongpassword")
	require.Error(t, err)
}

func TestNewUserService(t *testing.T) {
	repo := &mockUserRepo{}
	svc := NewUserService(repo)
	require.NotNil(t, svc)
}

func TestUserService_GetByEmail_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	repo.On("GetByEmail", ctx, "missing@test.com").Return(nil, errors.New("not found"))

	_, err := svc.GetByEmail(ctx, "missing@test.com")
	require.Error(t, err)
}

func TestUserService_GetByUsername_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	repo.On("GetByUsername", ctx, "missing").Return(nil, errors.New("not found"))

	_, err := svc.GetByUsername(ctx, "missing")
	require.Error(t, err)
}

func TestUserService_UpdateUser_UsernameConflict(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	existing := &models.User{Username: "oldname", Email: "test@test.com"}
	existing.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)

	conflict := &models.User{Username: "takenname"}
	repo.On("GetByUsername", ctx, "takenname").Return(conflict, nil)

	user := &models.User{Username: "takenname", Email: "test@test.com"}
	user.ID = 1
	err := svc.UpdateUser(ctx, user)
	require.Error(t, err)
}

func TestUserService_UpdateUser_EmailConflict(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	existing := &models.User{Username: "testuser", Email: "old@test.com"}
	existing.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)

	conflict := &models.User{Email: "taken@test.com"}
	repo.On("GetByEmail", ctx, "taken@test.com").Return(conflict, nil)

	user := &models.User{Username: "testuser", Email: "taken@test.com"}
	user.ID = 1
	err := svc.UpdateUser(ctx, user)
	require.Error(t, err)
}

func TestUserService_UpdateUser_UpdateError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	existing := &models.User{Username: "testuser", Email: "test@test.com"}
	existing.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	repo.On("Update", ctx, mock.Anything).Return(errors.New("db error"))

	user := &models.User{Username: "testuser", Email: "test@test.com"}
	user.ID = 1
	err := svc.UpdateUser(ctx, user)
	require.Error(t, err)
}

func TestUserService_CreateUser_CreateError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	repo.On("GetByUsername", ctx, "newuser").Return(nil, errors.New("not found"))
	repo.On("GetByEmail", ctx, "new@test.com").Return(nil, errors.New("not found"))
	repo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

	user := &models.User{Username: "newuser", Email: "new@test.com", Password: "password123"}
	err := svc.CreateUser(ctx, user)
	require.Error(t, err)
}

func TestUserService_CreateUser_DuplicateOnInsert(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepo{}
	svc := newUserSvc(repo)

	repo.On("GetByUsername", ctx, "newuser").Return(nil, errors.New("not found"))
	repo.On("GetByEmail", ctx, "new@test.com").Return(nil, errors.New("not found"))
	repo.On("Create", ctx, mock.Anything).Return(errors.New("unique constraint violation"))

	user := &models.User{Username: "newuser", Email: "new@test.com", Password: "password123"}
	err := svc.CreateUser(ctx, user)
	require.Error(t, err)
}

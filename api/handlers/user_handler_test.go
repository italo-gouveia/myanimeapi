package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"myanimeapi/api/mocks"
	"myanimeapi/api/models"
	apperrors "myanimeapi/internal/errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// sampleUser returns a minimal User for handler tests.
func sampleUser(id uint) *models.User {
	return &models.User{
		BaseModel: models.BaseModel{ID: id},
		Username:  "testuser",
		Email:     "test@example.com",
		IsActive:  true,
	}
}

// ─── Suite ────────────────────────────────────────────────────────────────────

type UserHandlerSuite struct {
	suite.Suite
	userSvc     *mocks.MockUserServiceInterface
	genreSvc    *mocks.MockGenreServiceInterface
	passwordSvc *mocks.MockPasswordResetServiceInterface
	handler     *UserHandler
}

func (s *UserHandlerSuite) SetupTest() {
	s.userSvc = new(mocks.MockUserServiceInterface)
	s.genreSvc = new(mocks.MockGenreServiceInterface)
	s.passwordSvc = new(mocks.MockPasswordResetServiceInterface)
	s.handler = NewUserHandler(s.userSvc, s.genreSvc, s.passwordSvc)
}

func (s *UserHandlerSuite) TearDownTest() {
	s.userSvc.AssertExpectations(s.T())
	s.genreSvc.AssertExpectations(s.T())
	s.passwordSvc.AssertExpectations(s.T())
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerSuite))
}

// ─── Register ─────────────────────────────────────────────────────────────────

func (s *UserHandlerSuite) TestRegister_Success() {
	s.userSvc.EXPECT().CreateUser(mock.Anything, mock.AnythingOfType("*models.User")).
		Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/users/register",
		jsonBuf(map[string]string{"username": "testuser", "email": "test@example.com", "password": "pass123"}))
	rr := httptest.NewRecorder()
	s.handler.Register(rr, req)

	s.Equal(http.StatusCreated, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal("success", body["status"])
}

func (s *UserHandlerSuite) TestRegister_InvalidBody() {
	req := httptest.NewRequest(http.MethodPost, "/users/register", jsonBuf(`{"bad json`))
	rr := httptest.NewRecorder()
	s.handler.Register(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *UserHandlerSuite) TestRegister_Conflict() {
	s.userSvc.EXPECT().CreateUser(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrConflict, "Username already exists", "", http.StatusConflict, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/users/register",
		jsonBuf(map[string]string{"username": "existing", "email": "existing@example.com", "password": "pass123"}))
	rr := httptest.NewRecorder()
	s.handler.Register(rr, req)

	assertErrorCode(s.T(), rr, http.StatusConflict, apperrors.ErrConflict)
}

func (s *UserHandlerSuite) TestRegister_InternalError() {
	s.userSvc.EXPECT().CreateUser(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/users/register",
		jsonBuf(map[string]string{"username": "testuser", "email": "test@example.com", "password": "pass123"}))
	rr := httptest.NewRecorder()
	s.handler.Register(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ─── Login ────────────────────────────────────────────────────────────────────

func (s *UserHandlerSuite) TestLogin_Success() {
	s.userSvc.EXPECT().ValidateUser(mock.Anything, "testuser", "pass123").
		Return(sampleUser(1), nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/users/login",
		jsonBuf(map[string]string{"username": "testuser", "password": "pass123"}))
	rr := httptest.NewRecorder()
	s.handler.Login(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal("success", body["status"])
}

func (s *UserHandlerSuite) TestLogin_InvalidBody() {
	req := httptest.NewRequest(http.MethodPost, "/users/login", jsonBuf(`{"bad json`))
	rr := httptest.NewRecorder()
	s.handler.Login(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *UserHandlerSuite) TestLogin_InvalidCredentials() {
	s.userSvc.EXPECT().ValidateUser(mock.Anything, "testuser", "wrong").
		Return(nil, apperrors.NewError(apperrors.ErrUnauthorized, "Invalid credentials", "", http.StatusUnauthorized, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/users/login",
		jsonBuf(map[string]string{"username": "testuser", "password": "wrong"}))
	rr := httptest.NewRecorder()
	s.handler.Login(rr, req)

	assertErrorCode(s.T(), rr, http.StatusUnauthorized, apperrors.ErrUnauthorized)
}

// ─── GetProfile ───────────────────────────────────────────────────────────────

func (s *UserHandlerSuite) TestGetProfile_Success() {
	s.userSvc.EXPECT().GetByID(mock.Anything, uint(1)).Return(sampleUser(1), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetProfile(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal("success", body["status"])
}

func (s *UserHandlerSuite) TestGetProfile_NotFound() {
	s.userSvc.EXPECT().GetByID(mock.Anything, uint(1)).
		Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "User not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetProfile(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

func (s *UserHandlerSuite) TestGetProfile_InternalError() {
	s.userSvc.EXPECT().GetByID(mock.Anything, uint(1)).
		Return(nil, apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.GetProfile(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *UserHandlerSuite) TestGetProfile_Unauthenticated() {
	req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
	req = req.WithContext(noAuthCtx())
	rr := httptest.NewRecorder()
	s.handler.GetProfile(rr, req)

	assertErrorCode(s.T(), rr, http.StatusUnauthorized, apperrors.ErrUnauthorized)
}

// ─── UpdateProfile ────────────────────────────────────────────────────────────

func (s *UserHandlerSuite) TestUpdateProfile_Success() {
	s.userSvc.EXPECT().GetByID(mock.Anything, uint(1)).Return(sampleUser(1), nil).Once()
	s.userSvc.EXPECT().UpdateUser(mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPut, "/users/profile",
		jsonBuf(map[string]string{"username": "newname", "email": "new@example.com"}))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.UpdateProfile(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal("success", body["status"])
}

func (s *UserHandlerSuite) TestUpdateProfile_InvalidBody() {
	req := httptest.NewRequest(http.MethodPut, "/users/profile", jsonBuf(`{"bad json`))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.UpdateProfile(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *UserHandlerSuite) TestUpdateProfile_UserNotFound() {
	s.userSvc.EXPECT().GetByID(mock.Anything, uint(1)).
		Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "User not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPut, "/users/profile",
		jsonBuf(map[string]string{"username": "newname"}))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.UpdateProfile(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

func (s *UserHandlerSuite) TestUpdateProfile_UpdateFails() {
	s.userSvc.EXPECT().GetByID(mock.Anything, uint(1)).Return(sampleUser(1), nil).Once()
	s.userSvc.EXPECT().UpdateUser(mock.Anything, mock.AnythingOfType("*models.User")).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPut, "/users/profile",
		jsonBuf(map[string]string{"username": "newname"}))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.UpdateProfile(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ─── ChangePassword ───────────────────────────────────────────────────────────

func (s *UserHandlerSuite) TestChangePassword_Success() {
	s.userSvc.EXPECT().GetByID(mock.Anything, uint(1)).Return(sampleUser(1), nil).Once()
	s.userSvc.EXPECT().ChangePassword(mock.Anything, uint(1), "oldpass", "newpass").Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/users/change-password",
		jsonBuf(map[string]string{"current_password": "oldpass", "new_password": "newpass"}))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.ChangePassword(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal("success", body["status"])
}

func (s *UserHandlerSuite) TestChangePassword_InvalidBody() {
	req := httptest.NewRequest(http.MethodPost, "/users/change-password", jsonBuf(`{bad json`))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.ChangePassword(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *UserHandlerSuite) TestChangePassword_WrongPassword() {
	s.userSvc.EXPECT().GetByID(mock.Anything, uint(1)).Return(sampleUser(1), nil).Once()
	s.userSvc.EXPECT().ChangePassword(mock.Anything, uint(1), "wrong", "newpass").
		Return(apperrors.NewError(apperrors.ErrUnauthorized, "Invalid current password", "", http.StatusUnauthorized, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/users/change-password",
		jsonBuf(map[string]string{"current_password": "wrong", "new_password": "newpass"}))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.ChangePassword(rr, req)

	assertErrorCode(s.T(), rr, http.StatusUnauthorized, apperrors.ErrUnauthorized)
}

// ─── DeactivateAccount ────────────────────────────────────────────────────────

func (s *UserHandlerSuite) TestDeactivateAccount_Success() {
	s.userSvc.EXPECT().GetByID(mock.Anything, uint(1)).Return(sampleUser(1), nil).Once()
	s.userSvc.EXPECT().DeactivateAccount(mock.Anything, uint(1), "mypassword").Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/users/deactivate",
		jsonBuf(map[string]string{"password": "mypassword"}))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.DeactivateAccount(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal("success", body["status"])
}

func (s *UserHandlerSuite) TestDeactivateAccount_InvalidBody() {
	req := httptest.NewRequest(http.MethodPost, "/users/deactivate", jsonBuf(`{bad json`))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.DeactivateAccount(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *UserHandlerSuite) TestDeactivateAccount_WrongPassword() {
	s.userSvc.EXPECT().GetByID(mock.Anything, uint(1)).Return(sampleUser(1), nil).Once()
	s.userSvc.EXPECT().DeactivateAccount(mock.Anything, uint(1), "wrong").
		Return(apperrors.NewError(apperrors.ErrUnauthorized, "Invalid password", "", http.StatusUnauthorized, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/users/deactivate",
		jsonBuf(map[string]string{"password": "wrong"}))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.DeactivateAccount(rr, req)

	assertErrorCode(s.T(), rr, http.StatusUnauthorized, apperrors.ErrUnauthorized)
}

// ─── RequestPasswordReset ─────────────────────────────────────────────────────

func (s *UserHandlerSuite) TestRequestPasswordReset_Success() {
	s.passwordSvc.EXPECT().RequestPasswordReset(mock.Anything, "test@example.com").Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/users/reset-password/request",
		jsonBuf(map[string]string{"email": "test@example.com"}))
	rr := httptest.NewRecorder()
	s.handler.RequestPasswordReset(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal("success", body["status"])
}

func (s *UserHandlerSuite) TestRequestPasswordReset_InvalidBody() {
	req := httptest.NewRequest(http.MethodPost, "/users/reset-password/request", jsonBuf(`{bad json`))
	rr := httptest.NewRecorder()
	s.handler.RequestPasswordReset(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *UserHandlerSuite) TestRequestPasswordReset_InternalError() {
	s.passwordSvc.EXPECT().RequestPasswordReset(mock.Anything, "test@example.com").
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/users/reset-password/request",
		jsonBuf(map[string]string{"email": "test@example.com"}))
	rr := httptest.NewRecorder()
	s.handler.RequestPasswordReset(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// ─── ResetPassword ────────────────────────────────────────────────────────────

func (s *UserHandlerSuite) TestResetPassword_Success() {
	s.passwordSvc.EXPECT().ResetPassword(mock.Anything, "valid-token", "newpass123").Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/users/reset-password",
		jsonBuf(map[string]string{"token": "valid-token", "new_password": "newpass123"}))
	rr := httptest.NewRecorder()
	s.handler.ResetPassword(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal("success", body["status"])
}

func (s *UserHandlerSuite) TestResetPassword_InvalidBody() {
	req := httptest.NewRequest(http.MethodPost, "/users/reset-password", jsonBuf(`{bad json`))
	rr := httptest.NewRecorder()
	s.handler.ResetPassword(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *UserHandlerSuite) TestResetPassword_InvalidToken() {
	s.passwordSvc.EXPECT().ResetPassword(mock.Anything, "bad-token", "newpass123").
		Return(apperrors.NewError(apperrors.ErrUnauthorized, "Invalid or expired token", "", http.StatusUnauthorized, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodPost, "/users/reset-password",
		jsonBuf(map[string]string{"token": "bad-token", "new_password": "newpass123"}))
	rr := httptest.NewRecorder()
	s.handler.ResetPassword(rr, req)

	assertErrorCode(s.T(), rr, http.StatusUnauthorized, apperrors.ErrUnauthorized)
}

// ─── DeleteAccount ────────────────────────────────────────────────────────────

func (s *UserHandlerSuite) TestDeleteAccount_Success() {
	s.userSvc.EXPECT().GetByID(mock.Anything, uint(1)).Return(sampleUser(1), nil).Once()
	s.userSvc.EXPECT().DeleteUser(mock.Anything, uint(1)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/users/delete", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.DeleteAccount(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	body := bodyJSON(s.T(), rr)
	s.Equal("success", body["status"])
}

func (s *UserHandlerSuite) TestDeleteAccount_UserNotFound() {
	s.userSvc.EXPECT().GetByID(mock.Anything, uint(1)).
		Return(nil, apperrors.NewError(apperrors.ErrResourceNotFound, "User not found", "", http.StatusNotFound, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/users/delete", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.DeleteAccount(rr, req)

	assertErrorCode(s.T(), rr, http.StatusNotFound, apperrors.ErrResourceNotFound)
}

func (s *UserHandlerSuite) TestDeleteAccount_InternalError() {
	s.userSvc.EXPECT().GetByID(mock.Anything, uint(1)).Return(sampleUser(1), nil).Once()
	s.userSvc.EXPECT().DeleteUser(mock.Anything, uint(1)).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	req := httptest.NewRequest(http.MethodDelete, "/users/delete", nil)
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.DeleteAccount(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

func (s *UserHandlerSuite) TestDeleteAccount_Unauthenticated() {
	req := httptest.NewRequest(http.MethodDelete, "/users/delete", nil)
	req = req.WithContext(noAuthCtx())
	rr := httptest.NewRecorder()
	s.handler.DeleteAccount(rr, req)

	assertErrorCode(s.T(), rr, http.StatusUnauthorized, apperrors.ErrUnauthorized)
}

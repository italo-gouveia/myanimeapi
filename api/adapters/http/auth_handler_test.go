package httphandler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"myanimeapi/api/mocks"
	"myanimeapi/api/models"
	apperrors "myanimeapi/internal/errors"
	"myanimeapi/internal/logger"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerSuite struct {
	suite.Suite
	svc     *mocks.MockAuthServiceInterface
	handler *AuthHandler
}

func (s *AuthHandlerSuite) SetupTest() {
	s.svc = new(mocks.MockAuthServiceInterface)
	s.handler = NewAuthHandler(s.svc, logger.New())
}

func (s *AuthHandlerSuite) TearDownTest() {
	s.svc.AssertExpectations(s.T())
}

func TestAuthHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerSuite))
}

// Register tests

func (s *AuthHandlerSuite) TestRegister_Success() {
	s.svc.EXPECT().RegisterUser(mock.Anything, mock.AnythingOfType("*models.User")).
		Run(func(ctx context.Context, user *models.User) {
			user.ID = 1
			user.Username = "testuser"
		}).
		Return(nil).Once()

	body := `{"username":"testuser","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.RegisterUserHandler(rr, req)

	s.Equal(http.StatusCreated, rr.Code)
	b := bodyJSON(s.T(), rr)
	s.Equal("success", b["status"])
	s.Equal("User registered successfully", b["message"])
}

func (s *AuthHandlerSuite) TestRegister_InvalidBody() {
	body := `{"bad json`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.RegisterUserHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *AuthHandlerSuite) TestRegister_Conflict() {
	s.svc.EXPECT().RegisterUser(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrConflict, "Username already exists", "", http.StatusConflict, nil, nil)).Once()

	body := `{"username":"existing","email":"existing@example.com","password":"pass123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.RegisterUserHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusConflict, apperrors.ErrConflict)
}

func (s *AuthHandlerSuite) TestRegister_InternalError() {
	s.svc.EXPECT().RegisterUser(mock.Anything, mock.Anything).
		Return(apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	body := `{"username":"testuser","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.RegisterUserHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

// Authenticate tests

func (s *AuthHandlerSuite) TestAuthenticate_Success() {
	s.svc.EXPECT().AuthenticateUser(mock.Anything, mock.AnythingOfType("*models.UserCredentials")).
		Return("jwt-token-here", nil).Once()

	body := `{"username":"testuser","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/authenticate", bytes.NewBufferString(body))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.AuthenticateUserHandler(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	b := bodyJSON(s.T(), rr)
	s.Equal("success", b["status"])
	data, _ := b["data"].(map[string]interface{})
	s.Equal("jwt-token-here", data["token"])
}

func (s *AuthHandlerSuite) TestAuthenticate_InvalidBody() {
	body := `{"bad json`
	req := httptest.NewRequest(http.MethodPost, "/auth/authenticate", bytes.NewBufferString(body))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.AuthenticateUserHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusBadRequest, apperrors.ErrInvalidInput)
}

func (s *AuthHandlerSuite) TestAuthenticate_Unauthorized() {
	s.svc.EXPECT().AuthenticateUser(mock.Anything, mock.Anything).
		Return("", apperrors.NewError(apperrors.ErrUnauthorized, "Invalid credentials", "", http.StatusUnauthorized, nil, nil)).Once()

	body := `{"username":"testuser","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/authenticate", bytes.NewBufferString(body))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.AuthenticateUserHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusUnauthorized, apperrors.ErrUnauthorized)
}

func (s *AuthHandlerSuite) TestAuthenticate_InternalError() {
	s.svc.EXPECT().AuthenticateUser(mock.Anything, mock.Anything).
		Return("", apperrors.NewError(apperrors.ErrInternalServer, "DB error", "", http.StatusInternalServerError, nil, nil)).Once()

	body := `{"username":"testuser","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/authenticate", bytes.NewBufferString(body))
	req = req.WithContext(testCtx())
	rr := httptest.NewRecorder()
	s.handler.AuthenticateUserHandler(rr, req)

	assertErrorCode(s.T(), rr, http.StatusInternalServerError, apperrors.ErrInternalServer)
}

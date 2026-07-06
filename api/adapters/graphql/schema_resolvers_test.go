package graphql

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"myanimeapi/api/mocks"
	"myanimeapi/api/models"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedID  uint
		expectError bool
	}{
		{
			name:        "valid positive integer string",
			input:       "42",
			expectedID:  42,
			expectError: false,
		},
		{
			name:        "zero returns error",
			input:       "0",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "negative integer returns error",
			input:       "-5",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "non-numeric string returns error",
			input:       "abc",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "empty string returns error",
			input:       "",
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "very large number within int range",
			input:       "2147483647",
			expectedID:  2147483647,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := parseID(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, uint(0), id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}

func TestDeleteAnimeResolver(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		setupMock   func(*mocks.MockAnimeServiceInterface)
		expectedOK  bool
		expectError bool
	}{
		{
			name:        "invalid ID returns error without calling service",
			id:          "abc",
			setupMock:   func(m *mocks.MockAnimeServiceInterface) {},
			expectedOK:  false,
			expectError: true,
		},
		{
			name: "valid ID service succeeds returns true",
			id:   "1",
			setupMock: func(m *mocks.MockAnimeServiceInterface) {
				m.On("DeleteAnime", mock.Anything, uint(1)).Return(nil)
			},
			expectedOK:  true,
			expectError: false,
		},
		{
			name: "valid ID service error returns false",
			id:   "2",
			setupMock: func(m *mocks.MockAnimeServiceInterface) {
				m.On("DeleteAnime", mock.Anything, uint(2)).Return(errors.New("not found"))
			},
			expectedOK:  false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAnime := &mocks.MockAnimeServiceInterface{}
			tt.setupMock(mockAnime)
			r := &mutationResolver{&Resolver{AnimeService: mockAnime}}

			ok, err := r.DeleteAnime(context.Background(), tt.id)

			assert.Equal(t, tt.expectedOK, ok)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockAnime.AssertExpectations(t)
		})
	}
}

func TestAnimeQueryResolver(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		setupMock   func(*mocks.MockAnimeServiceInterface)
		expectNil   bool
		expectError bool
	}{
		{
			name:        "invalid ID returns error without calling service",
			id:          "bad",
			setupMock:   func(m *mocks.MockAnimeServiceInterface) {},
			expectNil:   true,
			expectError: true,
		},
		{
			name: "valid ID service returns anime",
			id:   "5",
			setupMock: func(m *mocks.MockAnimeServiceInterface) {
				m.On("GetAnimeByID", mock.Anything, uint(5)).Return(&models.Anime{ID: 5, Title: "Naruto"}, nil)
			},
			expectNil:   false,
			expectError: false,
		},
		{
			name: "valid ID service error returns nil",
			id:   "7",
			setupMock: func(m *mocks.MockAnimeServiceInterface) {
				m.On("GetAnimeByID", mock.Anything, uint(7)).Return((*models.Anime)(nil), errors.New("db error"))
			},
			expectNil:   true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAnime := &mocks.MockAnimeServiceInterface{}
			tt.setupMock(mockAnime)
			r := &queryResolver{&Resolver{AnimeService: mockAnime}}

			result, err := r.Anime(context.Background(), tt.id)

			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockAnime.AssertExpectations(t)
		})
	}
}

func TestGenreQueryResolver(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		setupMock   func(*mocks.MockGenreServiceInterface)
		expectNil   bool
		expectError bool
	}{
		{
			name:        "invalid ID returns error without calling service",
			id:          "bad",
			setupMock:   func(m *mocks.MockGenreServiceInterface) {},
			expectNil:   true,
			expectError: true,
		},
		{
			name: "valid ID service returns genre",
			id:   "3",
			setupMock: func(m *mocks.MockGenreServiceInterface) {
				m.On("GetGenreByID", mock.Anything, uint(3)).Return(&models.Genre{ID: 3, Name: "Action"}, nil)
			},
			expectNil:   false,
			expectError: false,
		},
		{
			name: "valid ID service error returns nil",
			id:   "9",
			setupMock: func(m *mocks.MockGenreServiceInterface) {
				m.On("GetGenreByID", mock.Anything, uint(9)).Return((*models.Genre)(nil), errors.New("not found"))
			},
			expectNil:   true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGenre := &mocks.MockGenreServiceInterface{}
			tt.setupMock(mockGenre)
			r := &queryResolver{&Resolver{GenreService: mockGenre}}

			result, err := r.Genre(context.Background(), tt.id)

			if tt.expectNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockGenre.AssertExpectations(t)
		})
	}
}

package graphql

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"myanimeapi/api/adapters/graphql/model"
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

func strPtr(s string) *string  { return &s }
func intPtr(i int) *int        { return &i }
func f64Ptr(f float64) *float64 { return &f }

func TestGqlFilterToModel_NilFilter(t *testing.T) {
	result := gqlFilterToModel(nil)
	assert.Equal(t, models.AnimeFilter{}, result)
}

func TestGqlFilterToModel_AllFields(t *testing.T) {
	input := &model.AnimeFilterInput{
		Status:      strPtr("Airing"),
		Genre:       strPtr("Action"),
		Genres:      []string{"Action", "Adventure"},
		Tag:         strPtr("magic"),
		Tags:        []string{"magic", "school"},
		RatingMin:   f64Ptr(7.5),
		RatingMax:   f64Ptr(9.5),
		EpisodesMin: intPtr(12),
		EpisodesMax: intPtr(26),
		YearFrom:    intPtr(2010),
		YearTo:      intPtr(2023),
		SortBy:      strPtr("rating"),
		SortOrder:   strPtr("desc"),
	}

	result := gqlFilterToModel(input)

	assert.Equal(t, "Airing", result.Status)
	assert.Equal(t, "Action", result.Genre)
	assert.Equal(t, []string{"Action", "Adventure"}, result.Genres)
	assert.Equal(t, "magic", result.Tag)
	assert.Equal(t, []string{"magic", "school"}, result.Tags)
	assert.Equal(t, 7.5, result.RatingMin)
	assert.Equal(t, 9.5, result.RatingMax)
	assert.Equal(t, 12, result.EpisodesMin)
	assert.Equal(t, 26, result.EpisodesMax)
	assert.Equal(t, 2010, result.YearFrom)
	assert.Equal(t, 2023, result.YearTo)
	assert.Equal(t, "rating", result.SortBy)
	assert.Equal(t, "desc", result.SortOrder)
}

func TestGqlFilterToModel_PartialFields(t *testing.T) {
	input := &model.AnimeFilterInput{
		Status:    strPtr("Completed"),
		RatingMin: f64Ptr(8.0),
	}

	result := gqlFilterToModel(input)

	assert.Equal(t, "Completed", result.Status)
	assert.Equal(t, 8.0, result.RatingMin)
	assert.Equal(t, "", result.Genre)
	assert.Equal(t, 0.0, result.RatingMax)
	assert.Nil(t, result.Genres)
	assert.Nil(t, result.Tags)
}

func TestGqlFilterToModel_EmptySlices(t *testing.T) {
	input := &model.AnimeFilterInput{
		Genres: []string{},
		Tags:   []string{},
	}

	result := gqlFilterToModel(input)

	assert.Nil(t, result.Genres)
	assert.Nil(t, result.Tags)
}

func TestGqlFilterToModel_GenresAndTags(t *testing.T) {
	input := &model.AnimeFilterInput{
		Genres: []string{"Comedy", "Romance"},
		Tags:   []string{"slice-of-life"},
	}

	result := gqlFilterToModel(input)

	assert.Equal(t, []string{"Comedy", "Romance"}, result.Genres)
	assert.Equal(t, []string{"slice-of-life"}, result.Tags)
}

func TestAnimesResolver_Success(t *testing.T) {
	mockAnime := &mocks.MockAnimeServiceInterface{}
	page, limit := 1, 10
	filter := &model.AnimeFilterInput{
		Status: strPtr("Airing"),
	}

	mockAnime.EXPECT().
		GetAllAnimes(mock.Anything, page, limit, mock.MatchedBy(func(f models.AnimeFilter) bool {
			return f.Status == "Airing"
		})).
		Return([]*models.Anime{}, int64(0), nil)

	r := &queryResolver{&Resolver{AnimeService: mockAnime}}
	result, err := r.Animes(context.Background(), &page, &limit, filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Total)
	assert.Empty(t, result.Data)
	mockAnime.AssertExpectations(t)
}

func TestAnimesResolver_NilFilter(t *testing.T) {
	mockAnime := &mocks.MockAnimeServiceInterface{}
	page, limit := 1, 20

	mockAnime.EXPECT().
		GetAllAnimes(mock.Anything, page, limit, mock.MatchedBy(func(f models.AnimeFilter) bool {
			return f.Status == "" && f.Genre == ""
		})).
		Return([]*models.Anime{}, int64(0), nil)

	r := &queryResolver{&Resolver{AnimeService: mockAnime}}
	result, err := r.Animes(context.Background(), &page, &limit, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockAnime.AssertExpectations(t)
}

func TestAnimesResolver_ServiceError(t *testing.T) {
	mockAnime := &mocks.MockAnimeServiceInterface{}
	page, limit := 1, 10

	mockAnime.EXPECT().
		GetAllAnimes(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, int64(0), errors.New("db error"))

	r := &queryResolver{&Resolver{AnimeService: mockAnime}}
	result, err := r.Animes(context.Background(), &page, &limit, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockAnime.AssertExpectations(t)
}

func TestAnimesResolver_WithResults(t *testing.T) {
	mockAnime := &mocks.MockAnimeServiceInterface{}
	page, limit := 1, 10

	animes := []*models.Anime{
		{ID: 1, Title: "Naruto"},
		{ID: 2, Title: "Bleach"},
	}

	mockAnime.EXPECT().
		GetAllAnimes(mock.Anything, page, limit, mock.Anything).
		Return(animes, int64(2), nil)

	r := &queryResolver{&Resolver{AnimeService: mockAnime}}
	result, err := r.Animes(context.Background(), &page, &limit, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Data, 2)
	mockAnime.AssertExpectations(t)
}

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

// ---- helpers ----------------------------------------------------------------

func TestToGQLTag(t *testing.T) {
	tag := &models.Tag{}
	tag.ID = 1
	tag.Name = "Action"
	gql := toGQLTag(tag)
	assert.NotNil(t, gql)
	assert.Equal(t, "1", gql.ID)
	assert.Equal(t, "Action", gql.Name)
}

func TestToGQLTag_Nil(t *testing.T) {
	assert.Nil(t, toGQLTag(nil))
}

func TestToGQLCharacter(t *testing.T) {
	c := &models.Character{Name: "Naruto", Description: "Hero", VoiceActor: "VA"}
	c.ID = 1
	img := "http://img.com/naruto.jpg"
	c.ImageURL = &img

	gql := toGQLCharacter(c)
	assert.NotNil(t, gql)
	assert.Equal(t, "1", gql.ID)
	assert.Equal(t, "Naruto", gql.Name)
	assert.NotNil(t, gql.Description)
	assert.Equal(t, "Hero", *gql.Description)
	assert.NotNil(t, gql.VoiceActor)
	assert.Equal(t, "VA", *gql.VoiceActor)
	assert.Equal(t, &img, gql.ImageURL)
}

func TestToGQLCharacter_Nil(t *testing.T) {
	assert.Nil(t, toGQLCharacter(nil))
}

func TestToGQLCharacter_NoOptionals(t *testing.T) {
	c := &models.Character{Name: "Sasuke"}
	c.ID = 2
	gql := toGQLCharacter(c)
	assert.NotNil(t, gql)
	assert.Nil(t, gql.Description)
	assert.Nil(t, gql.VoiceActor)
}

// ---- auth resolvers ---------------------------------------------------------

func TestRegisterResolver_Success(t *testing.T) {
	mockAuth := &mocks.MockAuthServiceInterface{}
	mockAuth.On("RegisterUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	r := &mutationResolver{&Resolver{AuthService: mockAuth}}
	result, err := r.Register(context.Background(), "user", "user@example.com", "pass")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockAuth.AssertExpectations(t)
}

func TestRegisterResolver_Error(t *testing.T) {
	mockAuth := &mocks.MockAuthServiceInterface{}
	mockAuth.On("RegisterUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("conflict"))

	r := &mutationResolver{&Resolver{AuthService: mockAuth}}
	result, err := r.Register(context.Background(), "user", "user@example.com", "pass")
	assert.Error(t, err)
	assert.Nil(t, result)
	mockAuth.AssertExpectations(t)
}

func TestLoginResolver_Success(t *testing.T) {
	mockAuth := &mocks.MockAuthServiceInterface{}
	mockAuth.On("AuthenticateUser", mock.Anything, mock.AnythingOfType("*models.UserCredentials")).Return("jwt-token", nil)

	r := &mutationResolver{&Resolver{AuthService: mockAuth}}
	result, err := r.Login(context.Background(), "user", "pass")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "jwt-token", result.Token)
	mockAuth.AssertExpectations(t)
}

func TestLoginResolver_Error(t *testing.T) {
	mockAuth := &mocks.MockAuthServiceInterface{}
	mockAuth.On("AuthenticateUser", mock.Anything, mock.AnythingOfType("*models.UserCredentials")).Return("", errors.New("invalid"))

	r := &mutationResolver{&Resolver{AuthService: mockAuth}}
	result, err := r.Login(context.Background(), "user", "wrong")
	assert.Error(t, err)
	assert.Nil(t, result)
	mockAuth.AssertExpectations(t)
}

// ---- CreateAnime resolver ---------------------------------------------------

func TestCreateAnimeResolver_Success(t *testing.T) {
	mockAnime := &mocks.MockAnimeServiceInterface{}
	desc := "Great anime"
	ep := 24
	status := "Airing"

	anime := &models.Anime{Title: "Naruto", Description: desc, Episodes: ep, Status: status}
	anime.ID = 1

	mockAnime.On("CreateAnime", mock.Anything, mock.AnythingOfType("*models.Anime")).Return(nil)
	mockAnime.On("GetAnimeByID", mock.Anything, uint(0)).Return(anime, nil)

	r := &mutationResolver{&Resolver{AnimeService: mockAnime}}
	result, err := r.CreateAnime(context.Background(), "Naruto", &desc, &ep, &status)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockAnime.AssertExpectations(t)
}

func TestCreateAnimeResolver_CreateError(t *testing.T) {
	mockAnime := &mocks.MockAnimeServiceInterface{}
	mockAnime.On("CreateAnime", mock.Anything, mock.AnythingOfType("*models.Anime")).Return(errors.New("db error"))

	r := &mutationResolver{&Resolver{AnimeService: mockAnime}}
	result, err := r.CreateAnime(context.Background(), "Naruto", nil, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockAnime.AssertExpectations(t)
}

// ---- Character resolvers ----------------------------------------------------

func TestCreateCharacterResolver_Success(t *testing.T) {
	mockChar := &mocks.MockCharacterServiceInterface{}
	mockChar.On("CreateCharacter", mock.Anything, mock.AnythingOfType("*models.Character")).Return(nil)

	r := &mutationResolver{&Resolver{CharacterService: mockChar}}
	result, err := r.CreateCharacter(context.Background(), "Naruto", nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockChar.AssertExpectations(t)
}

func TestCreateCharacterResolver_Error(t *testing.T) {
	mockChar := &mocks.MockCharacterServiceInterface{}
	mockChar.On("CreateCharacter", mock.Anything, mock.AnythingOfType("*models.Character")).Return(errors.New("db error"))

	r := &mutationResolver{&Resolver{CharacterService: mockChar}}
	result, err := r.CreateCharacter(context.Background(), "Naruto", nil, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockChar.AssertExpectations(t)
}

func TestUpdateCharacterResolver_InvalidID(t *testing.T) {
	mockChar := &mocks.MockCharacterServiceInterface{}

	r := &mutationResolver{&Resolver{CharacterService: mockChar}}
	result, err := r.UpdateCharacter(context.Background(), "bad", nil, nil, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateCharacterResolver_Success(t *testing.T) {
	mockChar := &mocks.MockCharacterServiceInterface{}
	c := &models.Character{Name: "Naruto"}
	c.ID = 1
	name := "Naruto Uzumaki"

	mockChar.On("GetCharacterByID", mock.Anything, uint(1)).Return(c, nil)
	mockChar.On("UpdateCharacter", mock.Anything, mock.AnythingOfType("*models.Character")).Return(nil)

	r := &mutationResolver{&Resolver{CharacterService: mockChar}}
	result, err := r.UpdateCharacter(context.Background(), "1", &name, nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockChar.AssertExpectations(t)
}

func TestUpdateCharacterResolver_GetError(t *testing.T) {
	mockChar := &mocks.MockCharacterServiceInterface{}
	mockChar.On("GetCharacterByID", mock.Anything, uint(1)).Return((*models.Character)(nil), errors.New("not found"))

	r := &mutationResolver{&Resolver{CharacterService: mockChar}}
	result, err := r.UpdateCharacter(context.Background(), "1", nil, nil, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockChar.AssertExpectations(t)
}

func TestDeleteCharacterResolver_InvalidID(t *testing.T) {
	r := &mutationResolver{&Resolver{CharacterService: &mocks.MockCharacterServiceInterface{}}}
	ok, err := r.DeleteCharacter(context.Background(), "bad")
	assert.Error(t, err)
	assert.False(t, ok)
}

func TestDeleteCharacterResolver_Success(t *testing.T) {
	mockChar := &mocks.MockCharacterServiceInterface{}
	mockChar.On("DeleteCharacter", mock.Anything, uint(1)).Return(nil)

	r := &mutationResolver{&Resolver{CharacterService: mockChar}}
	ok, err := r.DeleteCharacter(context.Background(), "1")
	assert.NoError(t, err)
	assert.True(t, ok)
	mockChar.AssertExpectations(t)
}

func TestDeleteCharacterResolver_Error(t *testing.T) {
	mockChar := &mocks.MockCharacterServiceInterface{}
	mockChar.On("DeleteCharacter", mock.Anything, uint(1)).Return(errors.New("db error"))

	r := &mutationResolver{&Resolver{CharacterService: mockChar}}
	ok, err := r.DeleteCharacter(context.Background(), "1")
	assert.Error(t, err)
	assert.False(t, ok)
	mockChar.AssertExpectations(t)
}

// ---- AnimesByTitle, Genres, Characters resolvers ----------------------------

func TestAnimesByTitleResolver_Success(t *testing.T) {
	mockAnime := &mocks.MockAnimeServiceInterface{}
	animes := []*models.Anime{{ID: 1, Title: "Naruto"}}

	mockAnime.On("GetAllAnimes", mock.Anything, 1, 10, mock.MatchedBy(func(f models.AnimeFilter) bool {
		return f.Title == "Naruto"
	})).Return(animes, int64(1), nil)

	r := &queryResolver{&Resolver{AnimeService: mockAnime}}
	result, err := r.AnimesByTitle(context.Background(), "Naruto", nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
	mockAnime.AssertExpectations(t)
}

func TestAnimesByTitleResolver_Error(t *testing.T) {
	mockAnime := &mocks.MockAnimeServiceInterface{}
	mockAnime.On("GetAllAnimes", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	r := &queryResolver{&Resolver{AnimeService: mockAnime}}
	result, err := r.AnimesByTitle(context.Background(), "X", nil, nil, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockAnime.AssertExpectations(t)
}

func TestGenresResolver_Success(t *testing.T) {
	mockGenre := &mocks.MockGenreServiceInterface{}
	genres := []models.Genre{{Name: "Action"}}

	mockGenre.On("GetAllGenres", mock.Anything, 1, 10).Return(genres, int64(1), nil)

	r := &queryResolver{&Resolver{GenreService: mockGenre}}
	result, err := r.Genres(context.Background(), nil, nil)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockGenre.AssertExpectations(t)
}

func TestGenresResolver_Error(t *testing.T) {
	mockGenre := &mocks.MockGenreServiceInterface{}
	mockGenre.On("GetAllGenres", mock.Anything, 1, 10).Return(nil, int64(0), errors.New("db error"))

	r := &queryResolver{&Resolver{GenreService: mockGenre}}
	result, err := r.Genres(context.Background(), nil, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockGenre.AssertExpectations(t)
}

func TestCharacterResolver_InvalidID(t *testing.T) {
	r := &queryResolver{&Resolver{CharacterService: &mocks.MockCharacterServiceInterface{}}}
	result, err := r.Character(context.Background(), "bad")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCharacterResolver_Success(t *testing.T) {
	mockChar := &mocks.MockCharacterServiceInterface{}
	c := &models.Character{Name: "Naruto"}
	c.ID = 1

	mockChar.On("GetCharacterByID", mock.Anything, uint(1)).Return(c, nil)

	r := &queryResolver{&Resolver{CharacterService: mockChar}}
	result, err := r.Character(context.Background(), "1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockChar.AssertExpectations(t)
}

func TestCharactersResolver_Success(t *testing.T) {
	mockChar := &mocks.MockCharacterServiceInterface{}
	chars := []*models.Character{{Name: "Naruto"}, {Name: "Sasuke"}}

	mockChar.On("GetAllCharacters", mock.Anything, 1, 10).Return(chars, int64(2), nil)

	r := &queryResolver{&Resolver{CharacterService: mockChar}}
	result, err := r.Characters(context.Background(), nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Total)
	mockChar.AssertExpectations(t)
}

func TestCharactersResolver_Error(t *testing.T) {
	mockChar := &mocks.MockCharacterServiceInterface{}
	mockChar.On("GetAllCharacters", mock.Anything, 1, 10).Return(nil, int64(0), errors.New("db error"))

	r := &queryResolver{&Resolver{CharacterService: mockChar}}
	result, err := r.Characters(context.Background(), nil, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockChar.AssertExpectations(t)
}

func TestCharactersByNameResolver_Success(t *testing.T) {
	mockChar := &mocks.MockCharacterServiceInterface{}
	chars := []*models.Character{{Name: "Naruto"}}

	mockChar.On("GetCharactersByName", mock.Anything, "Naruto", 1, 10).Return(chars, int64(1), nil)

	r := &queryResolver{&Resolver{CharacterService: mockChar}}
	result, err := r.CharactersByName(context.Background(), "Naruto", nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
	mockChar.AssertExpectations(t)
}

func TestCharactersByNameResolver_Error(t *testing.T) {
	mockChar := &mocks.MockCharacterServiceInterface{}
	mockChar.On("GetCharactersByName", mock.Anything, "X", 1, 10).Return(nil, int64(0), errors.New("db error"))

	r := &queryResolver{&Resolver{CharacterService: mockChar}}
	result, err := r.CharactersByName(context.Background(), "X", nil, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockChar.AssertExpectations(t)
}

func TestCharactersByAnimeResolver_InvalidID(t *testing.T) {
	r := &queryResolver{&Resolver{CharacterService: &mocks.MockCharacterServiceInterface{}}}
	result, err := r.CharactersByAnime(context.Background(), "bad", nil, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCharactersByAnimeResolver_Success(t *testing.T) {
	mockChar := &mocks.MockCharacterServiceInterface{}
	chars := []*models.Character{{Name: "Naruto"}}

	mockChar.On("GetCharactersByAnime", mock.Anything, uint(1), 1, 10).Return(chars, int64(1), nil)

	r := &queryResolver{&Resolver{CharacterService: mockChar}}
	result, err := r.CharactersByAnime(context.Background(), "1", nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
	mockChar.AssertExpectations(t)
}

func TestMutationResolver(t *testing.T) {
	r := &Resolver{}
	assert.NotNil(t, r.Mutation())
}

func TestQueryResolver(t *testing.T) {
	r := &Resolver{}
	assert.NotNil(t, r.Query())
}

package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myanimeapi/api/models"
)

func newCharSvc() (*CharacterService, *mockCharacterRepo) {
	repo := &mockCharacterRepo{}
	return NewCharacterService(repo), repo
}

func TestCharacterService_NewCharacterService(t *testing.T) {
	repo := &mockCharacterRepo{}
	svc := NewCharacterService(repo)
	assert.NotNil(t, svc)
}

func TestCharacterService_GetCharacterByID_Success(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()
	want := &models.Character{Name: "Naruto"}
	want.ID = 1

	repo.On("GetByID", ctx, uint(1)).Return(want, nil)

	got, err := svc.GetCharacterByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, want, got)
	repo.AssertExpectations(t)
}

func TestCharacterService_GetCharacterByID_Error(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()

	repo.On("GetByID", ctx, uint(99)).Return(nil, errors.New("not found"))

	got, err := svc.GetCharacterByID(ctx, 99)
	assert.Nil(t, got)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_GetAllCharacters_Success(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()

	c1 := &models.Character{Name: "Naruto"}
	c1.ID = 1
	c2 := &models.Character{Name: "Sasuke"}
	c2.ID = 2

	repo.On("GetAll", ctx, 1, 10).Return([]interface{}{c1, c2}, int64(2), nil)

	chars, total, err := svc.GetAllCharacters(ctx, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, chars, 2)
	repo.AssertExpectations(t)
}

func TestCharacterService_GetAllCharacters_Error(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()

	repo.On("GetAll", ctx, 1, 10).Return(nil, int64(0), errors.New("db error"))

	chars, total, err := svc.GetAllCharacters(ctx, 1, 10)
	assert.Nil(t, chars)
	assert.Equal(t, int64(0), total)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_GetCharactersByAnime_Success(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()

	c := models.Character{Name: "Rock Lee"}
	c.ID = 3

	repo.On("GetByAnimeID", ctx, uint(1), 1, 10).Return([]models.Character{c}, int64(1), nil)

	chars, total, err := svc.GetCharactersByAnime(ctx, 1, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, chars, 1)
	repo.AssertExpectations(t)
}

func TestCharacterService_GetCharactersByAnime_Error(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()

	repo.On("GetByAnimeID", ctx, uint(1), 1, 10).Return(nil, int64(0), errors.New("db error"))

	chars, total, err := svc.GetCharactersByAnime(ctx, 1, 1, 10)
	assert.Nil(t, chars)
	assert.Equal(t, int64(0), total)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_GetCharactersByName_Success(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()

	c := models.Character{Name: "Itachi"}
	c.ID = 4

	repo.On("GetByName", ctx, "Itachi", 1, 10).Return([]models.Character{c}, int64(1), nil)

	chars, total, err := svc.GetCharactersByName(ctx, "Itachi", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, chars, 1)
	repo.AssertExpectations(t)
}

func TestCharacterService_GetCharactersByName_Error(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()

	repo.On("GetByName", ctx, "X", 1, 10).Return(nil, int64(0), errors.New("db error"))

	chars, total, err := svc.GetCharactersByName(ctx, "X", 1, 10)
	assert.Nil(t, chars)
	assert.Equal(t, int64(0), total)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_CreateCharacter_Success(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()
	c := &models.Character{Name: "Kakashi"}

	repo.On("Create", ctx, c).Return(nil)

	err := svc.CreateCharacter(ctx, c)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_CreateCharacter_Error(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()
	c := &models.Character{Name: "Kakashi"}

	repo.On("Create", ctx, c).Return(errors.New("db error"))

	err := svc.CreateCharacter(ctx, c)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_UpdateCharacter_Success(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()
	c := &models.Character{Name: "Jiraiya"}
	c.ID = 5

	repo.On("Update", ctx, c).Return(nil)

	err := svc.UpdateCharacter(ctx, c)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_UpdateCharacter_Error(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()
	c := &models.Character{Name: "Jiraiya"}
	c.ID = 5

	repo.On("Update", ctx, c).Return(errors.New("db error"))

	err := svc.UpdateCharacter(ctx, c)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_DeleteCharacter_Success(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()
	c := &models.Character{Name: "Gaara"}
	c.ID = 6

	repo.On("GetByID", ctx, uint(6)).Return(c, nil)
	repo.On("Delete", ctx, uint(6)).Return(nil)

	err := svc.DeleteCharacter(ctx, 6)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_DeleteCharacter_NotFound(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()

	repo.On("GetByID", ctx, uint(99)).Return(nil, errors.New("not found"))

	err := svc.DeleteCharacter(ctx, 99)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_DeleteCharacter_DeleteError(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()
	c := &models.Character{Name: "Gaara"}
	c.ID = 6

	repo.On("GetByID", ctx, uint(6)).Return(c, nil)
	repo.On("Delete", ctx, uint(6)).Return(errors.New("db error"))

	err := svc.DeleteCharacter(ctx, 6)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_AddCharactersToAnime_Success(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()
	ids := []uint{1, 2, 3}

	repo.On("AddToAnime", ctx, uint(10), ids).Return(nil)

	err := svc.AddCharactersToAnime(ctx, 10, ids)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_AddCharactersToAnime_Error(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()
	ids := []uint{1, 2}

	repo.On("AddToAnime", ctx, uint(10), ids).Return(errors.New("db error"))

	err := svc.AddCharactersToAnime(ctx, 10, ids)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_RemoveCharacterFromAnime_Success(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()

	repo.On("RemoveFromAnime", ctx, uint(10), uint(1)).Return(nil)

	err := svc.RemoveCharacterFromAnime(ctx, 10, 1)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCharacterService_RemoveCharacterFromAnime_Error(t *testing.T) {
	svc, repo := newCharSvc()
	ctx := context.Background()

	repo.On("RemoveFromAnime", ctx, uint(10), uint(1)).Return(errors.New("db error"))

	err := svc.RemoveCharacterFromAnime(ctx, 10, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

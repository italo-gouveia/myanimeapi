package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"myanimeapi/api/adapters/cache"
	"myanimeapi/api/models"
)

func newTagSvc(repo *mockTagRepo) *TagService {
	return &TagService{
		tagRepo: repo,
		log:     newTestLogger(),
		cache:   cache.NewNoopCache(),
	}
}

func TestTagService_GetTagByID_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockTagRepo{}
	svc := newTagSvc(repo)

	tag := &models.Tag{Name: "Action"}
	tag.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(tag, nil)

	result, err := svc.GetTagByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "Action", result.Name)
}

func TestTagService_GetTagByID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockTagRepo{}
	svc := newTagSvc(repo)

	repo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	_, err := svc.GetTagByID(ctx, 1)
	require.Error(t, err)
}

func TestTagService_GetAllTags_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockTagRepo{}
	svc := newTagSvc(repo)

	t1 := &models.Tag{Name: "Action"}
	t1.ID = 1
	t2 := &models.Tag{Name: "Comedy"}
	t2.ID = 2
	repo.On("GetAll", ctx, 1, 10).Return([]interface{}{t1, t2}, int64(2), nil)

	result, total, err := svc.GetAllTags(ctx, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

func TestTagService_GetAllTags_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockTagRepo{}
	svc := newTagSvc(repo)

	repo.On("GetAll", ctx, 1, 10).Return(nil, int64(0), errors.New("db error"))

	_, _, err := svc.GetAllTags(ctx, 1, 10)
	require.Error(t, err)
}

func TestTagService_CreateTag_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockTagRepo{}
	svc := newTagSvc(repo)

	repo.On("GetAll", ctx, 1, defaultPaginationLimit).Return([]interface{}{}, int64(0), nil)
	repo.On("Create", ctx, mock.Anything).Return(nil)

	tag := &models.Tag{Name: "NewTag"}
	err := svc.CreateTag(ctx, tag)
	require.NoError(t, err)
}

func TestTagService_CreateTag_NameConflict(t *testing.T) {
	ctx := context.Background()
	repo := &mockTagRepo{}
	svc := newTagSvc(repo)

	existing := &models.Tag{Name: "Action"}
	existing.ID = 1
	repo.On("GetAll", ctx, 1, defaultPaginationLimit).Return([]interface{}{existing}, int64(1), nil)

	tag := &models.Tag{Name: "Action"}
	err := svc.CreateTag(ctx, tag)
	require.Error(t, err)
}

func TestTagService_CreateTag_CreateError(t *testing.T) {
	ctx := context.Background()
	repo := &mockTagRepo{}
	svc := newTagSvc(repo)

	repo.On("GetAll", ctx, 1, defaultPaginationLimit).Return([]interface{}{}, int64(0), nil)
	repo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

	tag := &models.Tag{Name: "NewTag"}
	err := svc.CreateTag(ctx, tag)
	require.Error(t, err)
}

func TestTagService_UpdateTag_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockTagRepo{}
	svc := newTagSvc(repo)

	existing := &models.Tag{Name: "Action"}
	existing.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	repo.On("Update", ctx, mock.Anything).Return(nil)

	tag := &models.Tag{Name: "Action"}
	tag.ID = 1
	err := svc.UpdateTag(ctx, tag)
	require.NoError(t, err)
}

func TestTagService_UpdateTag_NameConflict(t *testing.T) {
	ctx := context.Background()
	repo := &mockTagRepo{}
	svc := newTagSvc(repo)

	existing := &models.Tag{Name: "Action"}
	existing.ID = 1
	other := &models.Tag{Name: "Comedy"}
	other.ID = 2
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	repo.On("GetAll", ctx, 1, defaultPaginationLimit).Return([]interface{}{existing, other}, int64(2), nil)

	tag := &models.Tag{Name: "Comedy"}
	tag.ID = 1
	err := svc.UpdateTag(ctx, tag)
	require.Error(t, err)
}

func TestTagService_DeleteTag_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockTagRepo{}
	svc := newTagSvc(repo)

	existing := &models.Tag{Name: "Action"}
	existing.ID = 1
	repo.On("GetByID", ctx, uint(1)).Return(existing, nil)
	repo.On("Delete", ctx, uint(1)).Return(nil)

	err := svc.DeleteTag(ctx, 1)
	require.NoError(t, err)
}

func TestTagService_DeleteTag_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockTagRepo{}
	svc := newTagSvc(repo)

	repo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("not found"))

	err := svc.DeleteTag(ctx, 1)
	require.Error(t, err)
}

func TestTagService_GetTagsByIDs_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockTagRepo{}
	svc := newTagSvc(repo)

	tags := []models.Tag{{Name: "Action"}, {Name: "Comedy"}}
	repo.On("GetByIDs", ctx, []uint{1, 2}).Return(tags, nil)

	result, err := svc.GetTagsByIDs(ctx, []uint{1, 2})
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestTagService_GetTagsByIDs_Error(t *testing.T) {
	ctx := context.Background()
	repo := &mockTagRepo{}
	svc := newTagSvc(repo)

	repo.On("GetByIDs", ctx, []uint{1}).Return(nil, errors.New("db error"))

	_, err := svc.GetTagsByIDs(ctx, []uint{1})
	require.Error(t, err)
}

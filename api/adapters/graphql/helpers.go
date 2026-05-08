package graphql

// helpers.go contains mapping utilities shared by all resolver files.
// Kept separate so gqlgen regeneration of schema.resolvers.go does not
// accidentally move them into its "unknown code" comment block.

import (
	"strconv"

	"myanimeapi/api/adapters/graphql/model"
	"myanimeapi/api/models"
)

func toGQLAnime(a *models.Anime) *model.Anime {
	if a == nil {
		return nil
	}
	id := strconv.FormatUint(uint64(a.ID), 10)
	gql := &model.Anime{
		ID:    id,
		Title: a.Title,
	}
	if a.Description != "" {
		gql.Description = &a.Description
	}
	if a.Rating != 0 {
		r := a.Rating
		gql.Rating = &r
	}
	if a.Episodes != 0 {
		e := a.Episodes
		gql.Episodes = &e
	}
	if a.Status != "" {
		s := a.Status
		gql.Status = &s
	}
	if !a.StartDate.IsZero() {
		gql.StartDate = &a.StartDate
	}
	if !a.EndDate.IsZero() {
		gql.EndDate = &a.EndDate
	}
	for i := range a.Genres {
		gql.Genres = append(gql.Genres, toGQLGenre(&a.Genres[i]))
	}
	for i := range a.Tags {
		gql.Tags = append(gql.Tags, toGQLTag(&a.Tags[i]))
	}
	return gql
}

func toGQLGenre(g *models.Genre) *model.Genre {
	if g == nil {
		return nil
	}
	return &model.Genre{
		ID:   strconv.FormatUint(uint64(g.ID), 10),
		Name: g.Name,
	}
}

func toGQLTag(t *models.Tag) *model.Tag {
	if t == nil {
		return nil
	}
	return &model.Tag{
		ID:   strconv.FormatUint(uint64(t.ID), 10),
		Name: t.Name,
	}
}

func paginationArgs(page, limit *int) (int, int) {
	p, l := 1, 10
	if page != nil && *page > 0 {
		p = *page
	}
	if limit != nil && *limit > 0 {
		l = *limit
	}
	return p, l
}

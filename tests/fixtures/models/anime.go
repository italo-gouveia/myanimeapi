package models

import (
	"time"

	"myanimeapi/api/models"
)

type AnimeFixture struct {
	ID          uint
	Title       string
	Description string
	Rating      float64
	Episodes    int
	Status      string
	StartDate   time.Time
	EndDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewAnimeFixture(opts ...AnimeFixtureOption) *models.Anime {
	fixture := &AnimeFixture{
		Title:       "Test Anime",
		Description: "A test anime for testing purposes",
		Rating:      8.5,
		Episodes:    12,
		Status:      "Ongoing",
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Time{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, opt := range opts {
		opt(fixture)
	}

	return &models.Anime{
		Title:       fixture.Title,
		Description: fixture.Description,
		Rating:      fixture.Rating,
		Episodes:    fixture.Episodes,
		Status:      fixture.Status,
		StartDate:   fixture.StartDate,
		EndDate:     fixture.EndDate,
		CreatedAt:   fixture.CreatedAt,
		UpdatedAt:   fixture.UpdatedAt,
	}
}

type AnimeFixtureOption func(*AnimeFixture)

func WithTitle(title string) AnimeFixtureOption {
	return func(f *AnimeFixture) {
		f.Title = title
	}
}

func WithDescription(description string) AnimeFixtureOption {
	return func(f *AnimeFixture) {
		f.Description = description
	}
}

func WithStatus(status string) AnimeFixtureOption {
	return func(f *AnimeFixture) {
		f.Status = status
	}
}

func WithEpisodes(episodes int) AnimeFixtureOption {
	return func(f *AnimeFixture) {
		f.Episodes = episodes
	}
}

func WithRating(rating float64) AnimeFixtureOption {
	return func(f *AnimeFixture) {
		f.Rating = rating
	}
}

func WithStartDate(startDate time.Time) AnimeFixtureOption {
	return func(f *AnimeFixture) {
		f.StartDate = startDate
	}
}

func WithEndDate(endDate time.Time) AnimeFixtureOption {
	return func(f *AnimeFixture) {
		f.EndDate = endDate
	}
}

// Predefined anime fixtures
func OngoingAnime() *models.Anime {
	return NewAnimeFixture()
}

func CompletedAnime() *models.Anime {
	return NewAnimeFixture(
		WithStatus("Completed"),
		WithEndDate(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)),
	)
}

func HighRatedAnime() *models.Anime {
	return NewAnimeFixture(WithRating(9.5))
}

func TestAnime() *models.Anime {
	return NewAnimeFixture(
		WithTitle("Test Anime"),
		WithDescription("A test anime for testing"),
	)
}

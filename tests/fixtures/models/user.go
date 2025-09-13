package models

import (
	"time"

	"myanimeapi/api/models"
)

type UserFixture struct {
	ID          uint
	Username    string
	Email       string
	Password    string
	IsActive    bool
	ProfilePic  string
	Bio         string
	SocialLinks models.JSON
	IsAdmin     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewUserFixture(opts ...UserFixtureOption) *models.User {
	fixture := &UserFixture{
		Username:    "testuser",
		Email:       "test@example.com",
		Password:    "password123",
		IsActive:    true,
		ProfilePic:  "",
		Bio:         "Test user bio",
		SocialLinks: models.JSON{},
		IsAdmin:     false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, opt := range opts {
		opt(fixture)
	}

	return &models.User{
		Username:    fixture.Username,
		Email:       fixture.Email,
		Password:    fixture.Password,
		IsActive:    fixture.IsActive,
		ProfilePic:  fixture.ProfilePic,
		Bio:         fixture.Bio,
		SocialLinks: fixture.SocialLinks,
		IsAdmin:     fixture.IsAdmin,
		BaseModel: models.BaseModel{
			CreatedAt: fixture.CreatedAt,
			UpdatedAt: fixture.UpdatedAt,
		},
	}
}

type UserFixtureOption func(*UserFixture)

func WithUsername(username string) UserFixtureOption {
	return func(f *UserFixture) {
		f.Username = username
	}
}

func WithEmail(email string) UserFixtureOption {
	return func(f *UserFixture) {
		f.Email = email
	}
}

func WithIsActive(isActive bool) UserFixtureOption {
	return func(f *UserFixture) {
		f.IsActive = isActive
	}
}

func WithIsAdmin(isAdmin bool) UserFixtureOption {
	return func(f *UserFixture) {
		f.IsAdmin = isAdmin
	}
}

func WithBio(bio string) UserFixtureOption {
	return func(f *UserFixture) {
		f.Bio = bio
	}
}

// Predefined user fixtures
func ActiveUser() *models.User {
	return NewUserFixture()
}

func InactiveUser() *models.User {
	return NewUserFixture(WithIsActive(false))
}

func AdminUser() *models.User {
	return NewUserFixture(
		WithUsername("admin"),
		WithEmail("admin@example.com"),
		WithIsAdmin(true),
	)
}

func TestUser() *models.User {
	return NewUserFixture(
		WithUsername("testuser"),
		WithEmail("test@example.com"),
	)
}

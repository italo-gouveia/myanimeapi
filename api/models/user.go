package models

// User represents a user in the system.
// It includes fields for user details, authentication, and relationships with reviews.
//
// Example:
//
//	{
//	  "id": 1,
//	  "created_at": "2025-02-20T19:27:00Z",
//	  "updated_at": "2025-02-20T19:27:00Z",
//	  "username": "john_doe",
//	  "email": "john@example.com",
//	  "is_admin": false,
//	  "is_active": true,
//	  "profile_pic": "https://example.com/profile.jpg",
//	  "bio": "Anime enthusiast",
//	  "social_links": {
//	    "twitter": "@johndoe",
//	    "instagram": "@johndoe"
//	  }
//	}
type User struct {
	BaseModel
	Username    string   `json:"username" gorm:"unique;not null" validate:"required,min=3,max=50" example:"john_doe"`                         // Unique username for the user
	Email       string   `json:"email" gorm:"unique;not null" validate:"required,email,min=5,max=100" example:"john@example.com"`             // Unique Email address of the user
	Password    string   `json:"-" gorm:"not null" validate:"required,min=5,max=100" example:"password123"`                                   // Password of the user (never serialized)
	IsActive    bool     `json:"is_active" gorm:"default:true" example:"true"`                                                                // Indicates if the user account is active
	ProfilePic  string   `json:"profile_pic" gorm:"default:null" example:"https://example.com/profile.jpg"`                                   // URL to the user's profile picture
	Bio         string   `json:"bio" gorm:"type:text;default:null" example:"Anime enthusiast"`                                                // User's biography
	SocialLinks JSON     `json:"social_links" gorm:"type:jsonb;default:null" example:"{\"twitter\":\"@johndoe\",\"instagram\":\"@johndoe\"}"` // User's social media links
	IsAdmin     bool     `json:"is_admin" gorm:"default:false" example:"false"`                                                               // Indicates if the user has admin privileges
	Reviews     []Review `json:"reviews,omitempty" gorm:"foreignKey:UserID"`                                                                  // User's reviews
	Favorites   []Anime  `json:"favorites,omitempty" gorm:"many2many:user_favorites;"`                                                        // User's favorite anime
	Genres      []Genre  `json:"genres,omitempty" gorm:"many2many:user_genres;"`                                                              // User's preferred genres
}

// UserUpdateRequest represents the request payload for updating a user's profile.
// It includes optional fields that can be updated, with validation rules for each field.
//
// Example:
//
//	{
//	  "username": "johndoe",
//	  "email": "john@example.com",
//	  "profile_pic": "https://example.com/profile.jpg",
//	  "bio": "Anime enthusiast",
//	  "social_links": {
//	    "twitter": "@johndoe",
//	    "instagram": "@johndoe"
//	  },
//	  "genre_ids": [1, 2, 3]
//	}
type UserUpdateRequest struct {
	Username    string `json:"username" validate:"omitempty,min=3,max=50" example:"johndoe"`                                      // New username (optional, 3-50 chars)
	Email       string `json:"email" validate:"omitempty,email" example:"john@example.com"`                                       // New email address (optional, valid email)
	ProfilePic  string `json:"profile_pic" validate:"omitempty,url" example:"https://example.com/profile.jpg"`                    // New profile picture URL (optional, valid URL)
	Bio         string `json:"bio" validate:"omitempty,max=500" example:"Anime enthusiast"`                                       // New biography (optional, max 500 chars)
	SocialLinks JSON   `json:"social_links" validate:"omitempty" example:"{\"twitter\":\"@johndoe\",\"instagram\":\"@johndoe\"}"` // New social media links (optional)
	GenreIDs    []uint `json:"genre_ids" validate:"omitempty,dive,min=1" example:"[1,2,3]"`                                       // New preferred genre IDs (optional)
}

// ChangePasswordRequest represents the request payload for changing a user's password.
// It requires both the current password for verification and the new password.
//
// Example:
//
//	{
//	  "current_password": "oldpassword123",
//	  "new_password": "newpassword123"
//	}
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=6" example:"oldpassword123"` // Current password for verification
	NewPassword     string `json:"new_password" validate:"required,min=6" example:"newpassword123"`     // New password to set
}

// DeactivateAccountRequest represents the request payload for deactivating a user account.
// It requires the user's password for security verification.
//
// Example:
//
//	{
//	  "password": "password123"
//	}
type DeactivateAccountRequest struct {
	Password string `json:"password" validate:"required" example:"password123"` // User's password for verification
}

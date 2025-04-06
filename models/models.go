package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	Fullname       string    `gorm:"not null"`
	Email          string    `gorm:"unique;not null"`
	Password       string    `gorm:"not null"`
	Phone          string    `gorm:"not null"`
	FirstTimeLogin bool      `gorm:"default:false"`
	OtpCode        string    `gorm:"not null"`
	RoleID         string    `gorm:"not null"`
	Role           Role      `gorm:"foreignKey:RoleID"`
	Status         string    `gorm:"not null"`
	EmailVerified  bool      `gorm:"default:null"`
	gorm.Model
}

type Api struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	ProjectID      uuid.UUID `gorm:"type:uuid;not null"`
	ApiName        string    `gorm:"not null"`
	ApiMethod      string    `gorm:"not null"`
	ApiDescription string    `gorm:"not null"`
	Endpoint       string    `gorm:"not null"`
	RequiresAuth   bool      `gorm:"default:false"`
	gorm.Model
}

type ApiKeys struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID         uuid.UUID `gorm:"type:uuid;not null"`
	User           User      `gorm:"foreignKey:UserID"`
	ClientID       string    `gorm:"not null"`
	ClientSecret   string    `gorm:"not null"`
	IsActive       bool      `gorm:"default:true"`
	PIN            string    `gorm:"default:null"`
	OAuthSignature string    `gorm:"default:null"`
	FloatBalance   string    `gorm:"default:null"`
	gorm.Model
}

type Notifications struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID  uuid.UUID `gorm:"type:uuid;not null"`
	User    User      `gorm:"foreignKey:UserID"`
	Message string    `gorm:"not null"`
	Status  bool      `gorm:"default:false"`
	gorm.Model
}

type Permission struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key"`
	Name string    `gorm:"not null;unique"`
}

type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primary_key"`
	Name        string       `gorm:"not null;unique"`            // e.g. "admin", "recruiter"
	Permissions []Permission `gorm:"many2many:role_permissions"` // Many-to-Many relationship
}

type Transactions struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Reference string    `gorm:"not null"`
	Amount    string    `gorm:"not null"`
	Status    string    `gorm:"not null"`
	Customer  string    `gorm:"not null"`
	Channel   string    `gorm:"not null"`
	Type      string    `gorm:"not null"`
	Narration string    `gorm:"default:null"`
	Date      time.Time
	gorm.Model
}

type CheckOutUrls struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	OrderID       string    `gorm:"default:null"`
	CustomerName  string    `gorm:"default:null"`
	CustomerEmail string    `gorm:"default:null"`
	SuccessUrl    string    `gorm:"default:null"`
	FailedUrl     string    `gorm:"default:null"`
	CancelUrl     string    `gorm:"default:null"`
	Amount        string    `gorm:"default:null"`
	GeneratedUrl  string    `gorm:"default:null"`
	TReference    string    `gorm:"default:null"`
	gorm.Model
}

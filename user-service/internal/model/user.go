package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
	RoleAdmin   UserRole = "admin"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Username  string    `gorm:"type:varchar(70);not null;uniqueIndex"`
	FirstName string    `gorm:"type:varchar(35);not null;"`
	LastName  string    `gorm:"type:varchar(35);not null;"`
	Role      UserRole  `gorm:"type:varchar(10);not null;default:'student';index"`
	IsActive  bool      `gorm:"type:boolean;not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

//func (u *User) BeforeSave(tx *gorm.DB) course_errors {
//	if u.ID == uuid.Nil {
//		u.ID = uuid.New()
//	}
//	return nil
//}

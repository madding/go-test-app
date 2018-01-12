package model

import (
	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	Login      string `gorm:"not null"`
	Pass       string `gorm:"not null"`
	WorkNumber int32
	gorm.Model
}

// Get - get user form database
func (u *User) Get(login string, pass string) error {
	return DBConn.Where("login = ? AND pass = ?", login, pass).First(u).Error
}

// Save - update user to database
func (u *User) Save() error {
	return DBConn.Save(u).Error
}

// CreateWhenNotExits - create user in database or skip when he exists
func (u *User) CreateWhenNotExits() error {
	if DBConn.NewRecord(u) {
		notFound := DBConn.Where("login = ? AND pass = ?", u.Login, u.Pass).First(&User{}).RecordNotFound()
		if notFound {
			return DBConn.Create(u).Error
		}
		return nil
	}
	return nil
}

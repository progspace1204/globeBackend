package models

type User struct {
	ID       uint   `gorm:"primarykey"`
	Username string `gorm:"size:64"`
	Author   bool   `gorm:"index"`
	Password string `gorm:"size:64"`
}

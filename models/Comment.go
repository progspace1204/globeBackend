package models

type Comment struct {
	ID      int32 `gorm:"primarykey"`
	User_id int32 `gorm:"index"`
	Post_id int32 `gorm:"index"`
	Content string
}

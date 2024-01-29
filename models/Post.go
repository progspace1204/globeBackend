package models

type Post struct {
	ID      int32  `gorm:"primarykey"`
	Title   string `gorm:"size:256"`
	Content string
	User_id int32 `gorm:"index"`
}

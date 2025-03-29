package models

type File struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(255);not null"`
	Path string `gorm:"type:varchar(255);not null"`
}

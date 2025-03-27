package models

type MainUser struct {
	ID     string `json:"id" gorm:"primaryKey"`
	NameTh string `json:"name_th"`
	NameEn string `json:"name_en"`
	Email  string `json:"email"`
}

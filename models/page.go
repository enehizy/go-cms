package models

type Page struct {
	BaseModel
	Title   string `gorm:"not null" json:"title"`
	Slug    string `gorm:"size:255;unique;not null" json:"slug"`
	Content string `gorm:"not null" json:"content"`
}

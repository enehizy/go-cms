package models

type Post struct {
	BaseModel
	Title string `gorm:"not null" json:"title"`

	Content   string `gorm:"not null" json:"content"`
	Published bool   `gorm:"default:false" json:"published"`

	Author string  `gorm:"not null" json:"author"`
	Media  []Media `gorm:"many2many" json:"media"`
}

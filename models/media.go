package models

var (
	ImagePng  = "image/png"
	ImageJpeg = "image/Jpeg"
	VideoMp4  = "video/mp4"
	VideoWebp = "video/webp"
)

type Media struct {
	BaseModel
	Uri      string `gorm:"not null; unique" json:"uri"`
	MimeType string `gorm:"column:mime_type;not null;check:mime_type in ('image/jpg','image/png','video/mp4','video/webp')" json:"mime_type"`
	Post     []Post `gorm:"many2many" json:"post"`
}

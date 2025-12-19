package models

func (Genre) TableName() string {
	return "genre"
}

type Genre struct {
	ID   int64  `gorm:"primaryKey;column:id"`
	Name string `gorm:"column:name;type:varchar(255);not null"`
}

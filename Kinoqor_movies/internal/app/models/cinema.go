package models

type Cinema struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	DetailsID int64          `gorm:"column:details_id;not null" json:"details_id"`
	Details   CinemaDetails  `json:"details"`
}

func (Cinema) TableName() string {
	return "cinema"
}

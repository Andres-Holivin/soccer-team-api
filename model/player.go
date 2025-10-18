package model

type PlayerPosition string

const (
	Striker    PlayerPosition = "Striker"
	Midfielder PlayerPosition = "Midfielder"
	Defender   PlayerPosition = "Defender"
	Goalkeeper PlayerPosition = "Goalkeeper"
)

type Player struct {
	BaseModel
	Name         string `gorm:"type:varchar(100);not null" json:"name" binding:"required"`
	Height       int    `gorm:"not null" json:"height" binding:"required"`
	Weight       int    `gorm:"not null" json:"weight" binding:"required"`
	JerseyNumber int    `gorm:"not null" json:"jersey_number" binding:"required"`

	Position PlayerPosition `gorm:"type:varchar(20);not null" json:"position" binding:"required,oneof=Striker Midfielder Defender Goalkeeper"`

	TeamID uint `gorm:"not null;index" json:"team_id" binding:"required"`
	Team   Team `gorm:"foreignKey:TeamID" json:"team,omitempty" binding:"-"`

	Goals []Goal `gorm:"foreignKey:PlayerID" json:"goals,omitempty" binding:"-"`
}

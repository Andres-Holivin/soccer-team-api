package model

type Goal struct {
	BaseModel
	MatchID    uint   `gorm:"not null;index" json:"match_id" binding:"required"`
	Match      Match  `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	PlayerID   uint   `gorm:"not null;index" json:"player_id" binding:"required"`
	Player     Player `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	GoalMinute int    `gorm:"not null" json:"goal_minute" binding:"required"`
}

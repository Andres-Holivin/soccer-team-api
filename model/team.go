package model

type Team struct {
	BaseModel
	Name        string `gorm:"type:varchar(100);not null" json:"name" binding:"required"`
	LogoURL     string `gorm:"type:varchar(255)" json:"logo_url"`
	FoundedYear int    `gorm:"not null" json:"founded_year" binding:"required"`
	Address     string `gorm:"type:text" json:"address" binding:"required"`
	City        string `gorm:"type:varchar(100);not null" json:"city" binding:"required"`

	Players     []Player `gorm:"foreignKey:TeamID" json:"players,omitempty"`
	HomeMatches []Match  `gorm:"foreignKey:HomeTeamID" json:"home_matches,omitempty"`
	AwayMatches []Match  `gorm:"foreignKey:AwayTeamID" json:"away_matches,omitempty"`
}

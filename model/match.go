package model

import (
	"time"
)

type MatchStatus string

const (
	MatchStatusScheduled MatchStatus = "scheduled"
	MatchStatusCompleted MatchStatus = "completed"
	MatchStatusCancelled MatchStatus = "cancelled"
)

type Match struct {
	BaseModel
	HomeTeamID    uint        `gorm:"not null;index" json:"home_team_id" binding:"required"`
	AwayTeamID    uint        `gorm:"not null;index" json:"away_team_id" binding:"required"`
	MatchDate     time.Time   `gorm:"type:timestamptz;not null" json:"match_date" binding:"required"`
	MatchTime     string      `gorm:"type:varchar(10);not null" json:"match_time" binding:"required"`
	Status        MatchStatus `gorm:"type:varchar(20);default:'scheduled'" json:"status"`
	HomeTeamScore int         `gorm:"default:0" json:"home_team_score"`
	AwayTeamScore int         `gorm:"default:0" json:"away_team_score"`

	HomeTeam Team   `gorm:"foreignKey:HomeTeamID" json:"home_team,omitempty" binding:"-"`
	AwayTeam Team   `gorm:"foreignKey:AwayTeamID" json:"away_team,omitempty" binding:"-"`
	Goals    []Goal `gorm:"foreignKey:MatchID" json:"goals,omitempty" binding:"-"`
}

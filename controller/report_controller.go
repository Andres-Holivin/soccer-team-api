package controller

import (
	"net/http"
	"soccer-team-api/model"
	"soccer-team-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MatchReportResponse struct {
	Match             model.Match    `json:"match"`
	HomeTeamName      string         `json:"home_team_name"`
	AwayTeamName      string         `json:"away_team_name"`
	FinalScore        string         `json:"final_score"`
	MatchResult       string         `json:"match_result"`
	TopScorer         *TopScorerInfo `json:"top_scorer"`
	HomeTeamTotalWins int64          `json:"home_team_total_wins"`
	AwayTeamTotalWins int64          `json:"away_team_total_wins"`
	GoalDetails       []GoalDetail   `json:"goal_details"`
}

type TopScorerInfo struct {
	PlayerID   uint   `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	GoalCount  int    `json:"goal_count"`
}

type GoalDetail struct {
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	GoalMinute int    `json:"goal_minute"`
}

const (
	QueryStatusEquals = "status = ?"
	QueryMatchDateLTE = "match_date <= ?"
	QueryTeamWins     = "(home_team_id = ? AND home_team_score > away_team_score) OR (away_team_id = ? AND away_team_score > home_team_score)"
)

func GetMatchReport(c *gin.Context) {
	id := c.Param("id")

	var match model.Match
	if err := utils.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Goals.Player.Team").First(&match, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Match not found", err.Error())
		return
	}

	if match.Status != model.MatchStatusCompleted {
		utils.ErrorResponse(c, http.StatusBadRequest, "Match not completed", nil)
		return
	}

	var matchResult string
	if match.HomeTeamScore > match.AwayTeamScore {
		matchResult = "Home Team Wins"
	} else if match.HomeTeamScore < match.AwayTeamScore {
		matchResult = "Away Team Wins"
	} else {
		matchResult = "Draw"
	}

	finalScore := strconv.Itoa(match.HomeTeamScore) + " - " + strconv.Itoa(match.AwayTeamScore)

	type PlayerGoalCount struct {
		PlayerID   uint
		PlayerName string
		TeamName   string
		GoalCount  int
	}

	var topScorer *TopScorerInfo
	if len(match.Goals) > 0 {
		goalCountMap := make(map[uint]PlayerGoalCount)
		for _, goal := range match.Goals {
			if count, exists := goalCountMap[goal.PlayerID]; exists {
				count.GoalCount++
				goalCountMap[goal.PlayerID] = count
			} else {
				goalCountMap[goal.PlayerID] = PlayerGoalCount{
					PlayerID:   goal.Player.ID,
					PlayerName: goal.Player.Name,
					TeamName:   goal.Player.Team.Name,
					GoalCount:  1,
				}
			}
		}

		maxGoals := 0
		for _, playerCount := range goalCountMap {
			if playerCount.GoalCount > maxGoals {
				maxGoals = playerCount.GoalCount
				topScorer = &TopScorerInfo{
					PlayerID:   playerCount.PlayerID,
					PlayerName: playerCount.PlayerName,
					TeamName:   playerCount.TeamName,
					GoalCount:  playerCount.GoalCount,
				}
			}
		}
	}

	goalDetails := make([]GoalDetail, 0)
	for _, goal := range match.Goals {
		goalDetails = append(goalDetails, GoalDetail{
			PlayerName: goal.Player.Name,
			TeamName:   goal.Player.Team.Name,
			GoalMinute: goal.GoalMinute,
		})
	}

	var homeTeamTotalWins int64
	utils.DB.Model(&model.Match{}).
		Where(QueryTeamWins, match.HomeTeamID, match.HomeTeamID).
		Where(QueryStatusEquals, model.MatchStatusCompleted).
		Where(QueryMatchDateLTE, match.MatchDate).
		Count(&homeTeamTotalWins)

	var awayTeamTotalWins int64
	utils.DB.Model(&model.Match{}).
		Where(QueryTeamWins, match.AwayTeamID, match.AwayTeamID).
		Where(QueryStatusEquals, model.MatchStatusCompleted).
		Where(QueryMatchDateLTE, match.MatchDate).
		Count(&awayTeamTotalWins)

	report := MatchReportResponse{
		Match:             match,
		HomeTeamName:      match.HomeTeam.Name,
		AwayTeamName:      match.AwayTeam.Name,
		FinalScore:        finalScore,
		MatchResult:       matchResult,
		TopScorer:         topScorer,
		HomeTeamTotalWins: homeTeamTotalWins,
		AwayTeamTotalWins: awayTeamTotalWins,
		GoalDetails:       goalDetails,
	}

	utils.SuccessResponse(c, http.StatusOK, "Match report retrieved successfully", report)
}

func GetAllMatchReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var matches []model.Match
	var total int64

	query := utils.DB.Model(&model.Match{}).
		Where(QueryStatusEquals, model.MatchStatusCompleted).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Goals.Player.Team")

	query.Count(&total)

	if err := query.Order("match_date DESC, match_time DESC").Offset(offset).Limit(limit).Find(&matches).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error retrieving match reports", err.Error())
		return
	}

	reports := make([]MatchReportResponse, 0)
	for _, match := range matches {
		var matchResult string
		if match.HomeTeamScore > match.AwayTeamScore {
			matchResult = "Home Team Wins"
		} else if match.HomeTeamScore < match.AwayTeamScore {
			matchResult = "Away Team Wins"
		} else {
			matchResult = "Draw"
		}

		finalScore := strconv.Itoa(match.HomeTeamScore) + " - " + strconv.Itoa(match.AwayTeamScore)

		var topScorer *TopScorerInfo
		if len(match.Goals) > 0 {
			type PlayerGoalCount struct {
				PlayerID   uint
				PlayerName string
				TeamName   string
				GoalCount  int
			}

			goalCountMap := make(map[uint]PlayerGoalCount)
			for _, goal := range match.Goals {
				if count, exists := goalCountMap[goal.PlayerID]; exists {
					count.GoalCount++
					goalCountMap[goal.PlayerID] = count
				} else {
					goalCountMap[goal.PlayerID] = PlayerGoalCount{
						PlayerID:   goal.Player.ID,
						PlayerName: goal.Player.Name,
						TeamName:   goal.Player.Team.Name,
						GoalCount:  1,
					}
				}
			}

			maxGoals := 0
			for _, playerCount := range goalCountMap {
				if playerCount.GoalCount > maxGoals {
					maxGoals = playerCount.GoalCount
					topScorer = &TopScorerInfo{
						PlayerID:   playerCount.PlayerID,
						PlayerName: playerCount.PlayerName,
						TeamName:   playerCount.TeamName,
						GoalCount:  playerCount.GoalCount,
					}
				}
			}
		}

		goalDetails := make([]GoalDetail, 0)
		for _, goal := range match.Goals {
			goalDetails = append(goalDetails, GoalDetail{
				PlayerName: goal.Player.Name,
				TeamName:   goal.Player.Team.Name,
				GoalMinute: goal.GoalMinute,
			})
		}

		// Calculate wins
		var homeTeamTotalWins int64
		utils.DB.Model(&model.Match{}).
			Where(QueryTeamWins, match.HomeTeamID, match.HomeTeamID).
			Where(QueryStatusEquals, model.MatchStatusCompleted).
			Where(QueryMatchDateLTE, match.MatchDate).
			Count(&homeTeamTotalWins)

		var awayTeamTotalWins int64
		utils.DB.Model(&model.Match{}).
			Where(QueryTeamWins, match.AwayTeamID, match.AwayTeamID).
			Where(QueryStatusEquals, model.MatchStatusCompleted).
			Where(QueryMatchDateLTE, match.MatchDate).
			Count(&awayTeamTotalWins)

		report := MatchReportResponse{
			Match:             match,
			HomeTeamName:      match.HomeTeam.Name,
			AwayTeamName:      match.AwayTeam.Name,
			FinalScore:        finalScore,
			MatchResult:       matchResult,
			TopScorer:         topScorer,
			HomeTeamTotalWins: homeTeamTotalWins,
			AwayTeamTotalWins: awayTeamTotalWins,
			GoalDetails:       goalDetails,
		}

		reports = append(reports, report)
	}

	utils.SuccessResponseWithPagination(c, http.StatusOK, "Match reports retrieved successfully", reports, page, limit, total)
}

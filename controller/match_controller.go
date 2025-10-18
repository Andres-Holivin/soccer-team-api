package controller

import (
	"net/http"
	"soccer-team-api/model"
	"soccer-team-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateMatch(c *gin.Context) {
	var match model.Match
	if err := c.ShouldBindJSON(&match); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate teams exist
	var homeTeam, awayTeam model.Team
	if err := utils.DB.First(&homeTeam, match.HomeTeamID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Home team not found", nil)
		return
	}
	if err := utils.DB.First(&awayTeam, match.AwayTeamID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Away team not found", nil)
		return
	}

	// Validate teams are different
	if match.HomeTeamID == match.AwayTeamID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Teams must be different", nil)
		return
	}

	// Set initial status
	match.Status = model.MatchStatusScheduled
	match.HomeTeamScore = 0
	match.AwayTeamScore = 0

	if err := utils.DB.Create(&match).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create match", err.Error())
		return
	}
	utils.DB.Preload("HomeTeam").Preload("AwayTeam").First(&match, match.ID)

	utils.SuccessResponse(c, http.StatusCreated, "Match scheduled successfully", match)
}
func GetMatches(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	status := c.Query("status")
	teamID := c.Query("team_id")

	var matches []model.Match
	var total int64

	query := utils.DB.Model(&model.Match{}).Preload("HomeTeam").Preload("AwayTeam")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if teamID != "" {
		query = query.Where("home_team_id = ? OR away_team_id = ?", teamID, teamID)
	}

	query.Count(&total)

	if err := query.Order("match_date DESC, match_time DESC").Offset(offset).Limit(limit).Find(&matches).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve matches", err.Error())
		return
	}

	utils.SuccessResponseWithPagination(c, http.StatusOK, "Matches retrieved successfully", matches, page, limit, total)
}

func GetMatch(c *gin.Context) {
	id := c.Param("id")

	var match model.Match
	if err := utils.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Goals.Player").First(&match, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Match not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Match retrieved successfully", match)
}

func UpdateMatch(c *gin.Context) {
	id := c.Param("id")

	var match model.Match
	if err := utils.DB.First(&match, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Match not found", err.Error())
		return
	}

	if match.Status != model.MatchStatusScheduled {
		utils.ErrorResponse(c, http.StatusBadRequest, "Match can only be updated if scheduled", nil)
		return
	}

	var updateData model.Match
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if updateData.HomeTeamID == updateData.AwayTeamID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Teams must be different", nil)
		return
	}

	match.HomeTeamID = updateData.HomeTeamID
	match.AwayTeamID = updateData.AwayTeamID
	match.MatchDate = updateData.MatchDate
	match.MatchTime = updateData.MatchTime

	if err := utils.DB.Save(&match).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update match", err.Error())
		return
	}

	utils.DB.Preload("HomeTeam").Preload("AwayTeam").First(&match, match.ID)

	utils.SuccessResponse(c, http.StatusOK, "Match updated successfully", match)
}

func DeleteMatch(c *gin.Context) {
	id := c.Param("id")

	var match model.Match
	if err := utils.DB.First(&match, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Match not found", err.Error())
		return
	}

	if err := utils.DB.Delete(&match).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete match", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Match deleted successfully", nil)
}

type MatchResultRequest struct {
	Goals []GoalInput `json:"goals" binding:"required"`
}

type GoalInput struct {
	PlayerID   uint `json:"player_id" binding:"required"`
	GoalMinute int  `json:"goal_minute" binding:"required,min=1,max=120"`
}

func RecordMatchResult(c *gin.Context) {
	id := c.Param("id")

	var match model.Match
	if err := utils.DB.Preload("HomeTeam").Preload("AwayTeam").First(&match, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Match not found", err.Error())
		return
	}

	// Only allow recording results for scheduled matches
	if match.Status != model.MatchStatusScheduled {
		utils.ErrorResponse(c, http.StatusBadRequest, "Match result already recorded", nil)
		return
	}

	var req MatchResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Start transaction
	tx := utils.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Validate and create goals
	homeScore := 0
	awayScore := 0

	for _, goalInput := range req.Goals {
		// Validate player exists and belongs to one of the teams
		var player model.Player
		if err := tx.First(&player, goalInput.PlayerID).Error; err != nil {
			tx.Rollback()
			utils.ErrorResponse(c, http.StatusNotFound, "Player not found", err.Error())
			return
		}

		if player.TeamID != match.HomeTeamID && player.TeamID != match.AwayTeamID {
			tx.Rollback()
			utils.ErrorResponse(c, http.StatusBadRequest, "Player is not in the match teams", nil)
			return
		}

		// Create goal record
		goal := model.Goal{
			MatchID:    match.ID,
			PlayerID:   goalInput.PlayerID,
			GoalMinute: goalInput.GoalMinute,
		}

		if err := tx.Create(&goal).Error; err != nil {
			tx.Rollback()
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to record goal", err.Error())
			return
		}

		// Update score
		if player.TeamID == match.HomeTeamID {
			homeScore++
		} else {
			awayScore++
		}
	}

	// Update match scores and status
	match.HomeTeamScore = homeScore
	match.AwayTeamScore = awayScore
	match.Status = model.MatchStatusCompleted

	if err := tx.Save(&match).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update match", err.Error())
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction", err.Error())
		return
	}

	// Load complete match data
	utils.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Goals.Player").First(&match, match.ID)

	utils.SuccessResponse(c, http.StatusOK, "Match result recorded successfully", match)
}

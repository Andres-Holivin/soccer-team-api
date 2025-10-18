package controller

import (
	"net/http"
	"soccer-team-api/model"
	"soccer-team-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePlayer(c *gin.Context) {
	var player model.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	var team model.Team
	if err := utils.DB.First(&team, player.TeamID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Team not found", nil)
		return
	}

	var existingPlayer model.Player
	err := utils.DB.Where("team_id = ? AND jersey_number = ?", player.TeamID, player.JerseyNumber).First(&existingPlayer).Error
	if err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Jersey number already taken", nil)
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to check jersey number", err.Error())
		return
	}

	if err := utils.DB.Create(&player).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create player", err.Error())
		return
	}

	utils.DB.Preload("Team").First(&player, player.ID)

	utils.SuccessResponse(c, http.StatusCreated, "Player created successfully", player)
}

func GetPlayers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	teamID := c.Query("team_id")
	position := c.Query("position")

	var players []model.Player
	var total int64

	query := utils.DB.Model(&model.Player{}).Preload("Team")

	if teamID != "" {
		query = query.Where("team_id = ?", teamID)
	}

	if position != "" {
		query = query.Where("position = ?", position)
	}

	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&players).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve players", err.Error())
		return
	}

	utils.SuccessResponseWithPagination(c, http.StatusOK, "Players retrieved successfully", players, page, limit, total)
}

func GetPlayer(c *gin.Context) {
	id := c.Param("id")

	var player model.Player
	if err := utils.DB.Preload("Team").Preload("Goals").First(&player, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Player not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Player retrieved successfully", player)
}

func UpdatePlayer(c *gin.Context) {
	id := c.Param("id")

	var player model.Player
	if err := utils.DB.First(&player, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Player not found", err.Error())
		return
	}

	var updateData model.Player
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if updateData.JerseyNumber != player.JerseyNumber {
		var existingPlayer model.Player
		err := utils.DB.Where("team_id = ? AND jersey_number = ? AND id != ?", player.TeamID, updateData.JerseyNumber, player.ID).First(&existingPlayer).Error
		if err == nil {
			utils.ErrorResponse(c, http.StatusConflict, "Jersey number already taken", nil)
			return
		}
	}

	player.Name = updateData.Name
	player.Height = updateData.Height
	player.Weight = updateData.Weight
	player.Position = updateData.Position
	player.JerseyNumber = updateData.JerseyNumber

	if err := utils.DB.Save(&player).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update player", err.Error())
		return
	}

	utils.DB.Preload("Team").First(&player, player.ID)

	utils.SuccessResponse(c, http.StatusOK, "Player updated successfully", player)
}

func DeletePlayer(c *gin.Context) {
	id := c.Param("id")

	var player model.Player
	if err := utils.DB.First(&player, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Player not found", err.Error())
		return
	}

	if err := utils.DB.Delete(&player).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete player", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Player deleted successfully", nil)
}

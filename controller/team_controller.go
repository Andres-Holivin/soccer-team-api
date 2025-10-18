package controller

import (
	"net/http"
	"soccer-team-api/model"
	"soccer-team-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTeam(c *gin.Context) {
	name := c.PostForm("name")
	foundedYear := c.PostForm("founded_year")
	address := c.PostForm("address")
	city := c.PostForm("city")

	if name == "" || foundedYear == "" || address == "" || city == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "All fields are required")
		return
	}

	year, err := strconv.Atoi(foundedYear)
	if err != nil || year < 1800 || year > 2100 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid founded year", "Founded year must be a valid year between 1800 and 2100")
		return
	}

	var logoURL string
	file, err := c.FormFile("logo")
	if err == nil {
		logoURL, _, err = utils.UploadImage(file, "team-logos")
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload logo", err.Error())
			return
		}
	} else {
		logoURL = c.PostForm("logo_url")
		if logoURL == "" {
			utils.ErrorResponse(c, http.StatusBadRequest, "Logo is required", "Either upload a logo file or provide a logo URL")
			return
		}
	}

	team := model.Team{
		Name:        name,
		LogoURL:     logoURL,
		FoundedYear: year,
		Address:     address,
		City:        city,
	}

	if err := utils.DB.Create(&team).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create team", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Team created successfully", team)
}

func GetTeams(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var teams []model.Team
	var total int64

	query := utils.DB.Model(&model.Team{})

	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&teams).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve teams", err.Error())
		return
	}

	utils.SuccessResponseWithPagination(c, http.StatusOK, "Teams retrieved successfully", teams, page, limit, total)
}

func GetTeam(c *gin.Context) {
	id := c.Param("id")

	var team model.Team
	if err := utils.DB.Preload("Players").First(&team, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Team not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Team retrieved successfully", team)
}

func UpdateTeam(c *gin.Context) {
	id := c.Param("id")

	var team model.Team
	if err := utils.DB.First(&team, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Team not found", err.Error())
		return
	}

	if name := c.PostForm("name"); name != "" {
		team.Name = name
	}
	if foundedYear := c.PostForm("founded_year"); foundedYear != "" {
		year, err := strconv.Atoi(foundedYear)
		if err == nil && year >= 1800 && year <= 2100 {
			team.FoundedYear = year
		}
	}
	if address := c.PostForm("address"); address != "" {
		team.Address = address
	}
	if city := c.PostForm("city"); city != "" {
		team.City = city
	}

	file, err := c.FormFile("logo")
	if err == nil {
		logoURL, _, err := utils.UploadImage(file, "team-logos")
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload logo", err.Error())
			return
		}
		team.LogoURL = logoURL
	} else if logoURL := c.PostForm("logo_url"); logoURL != "" {
		team.LogoURL = logoURL
	}

	if err := utils.DB.Save(&team).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update team", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Team updated successfully", team)
}

func DeleteTeam(c *gin.Context) {
	id := c.Param("id")

	var team model.Team
	if err := utils.DB.First(&team, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Team not found", err.Error())
		return
	}

	if err := utils.DB.Delete(&team).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete team", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Team deleted successfully", nil)
}

package controllers

import (
	"georgslauf/models"
	log "github.com/sirupsen/logrus"
)

var (
	totalLogin int64 = 0
	totalGroup int64 = 0
	totalGroupPoint int64 = 0
	totalRole int64 = 0
	totalStation int64 = 0
	totalStationPoint int64 = 0
	totalTribe int64 = 0
)

func InitTotal() {
	totalLogin = InitLogin()
	totalGroup = InitGroup()
	totalGroupPoint = InitGroupPoint()
	totalRole = InitRole()
	totalStation = InitStation()
	totalStationPoint = InitStationPoint()
	totalTribe = InitTribe()
}

func InitLogin() int64 {
	var model []models.Login
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init login failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func InitGroup() int64 {
	var model []models.Group
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init group failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func InitGroupPoint() int64 {
	var model []models.GroupPoint
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init grouppoint failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func InitRole() int64 {
	var model []models.Role
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init role failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func InitStation() int64 {
	var model []models.Station
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init station failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func InitStationPoint() int64 {
	var model []models.StationPoint
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init stationpoint failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func InitTribe() int64 {
	var model []models.Tribe
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init tribe failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

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
	totalGrouping int64 = 0
	totalContent int64 = 0
	totalGroupTop int64 = 0
	totalStationTop int64 = 0
	totalRun int64 = 0
	totalContentType int64 = 0
)

func InitTotal() {
	totalLogin = initLogin()
	totalGroup = initGroup()
	totalGroupPoint = initGroupPoint()
	totalRole = initRole()
	totalStation = initStation()
	totalStationPoint = initStationPoint()
	totalTribe = initTribe()
	totalGrouping = initGrouping()
	totalContent = initContent()
	totalGroupTop = initGroupTop()
	totalStationTop = initStationPoint()
	totalRun = initRun()
	totalContentType = initContentType()
	log.Info("Database total calculation complete.")
}

func initContentType() int64 {
	var model []models.ContentType
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init run failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}


func initRun() int64 {
	var model []models.Run
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init run failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func initStationTop() int64 {
	var model []models.StationTop
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init stationtop failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func initGroupTop() int64 {
	var model []models.GroupTop
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init grouptop failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func initLogin() int64 {
	var model []models.Login
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init login failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func initGroup() int64 {
	var model []models.Group
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init group failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func initGroupPoint() int64 {
	var model []models.GroupPoint
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init grouppoint failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func initRole() int64 {
	var model []models.Role
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init role failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func initStation() int64 {
	var model []models.Station
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init station failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func initStationPoint() int64 {
	var model []models.StationPoint
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init stationpoint failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func initTribe() int64 {
	var model []models.Tribe
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init tribe failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func initGrouping() int64 {
	var model []models.Grouping
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init grouping failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}

func initContent() int64 {
	var model []models.Content
	result := models.DB.Find(&model)
	if result.Error != nil {
		log.Warn("Init content failed.")
		return 0
	} else {
		return result.RowsAffected
	}
}
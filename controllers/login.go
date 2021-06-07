package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"georgslauf/models"
	"strconv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	jwt "github.com/appleboy/gin-jwt/v2"
)

type AuthenticateLoginInput struct {
	Username	string	`json:"username"	binding:"required"`
	Password	string	`json:"password"	binding:"required"`
}

type CreateLoginInput struct {
	Username 	string		`json:"username"	binding:"required"`
	Password	string		`json:"password"	binding:"required"`
	Reset		bool		`json:"reset"		binding:"required"`
	Active		bool		`json:"active"		binding:"required"`
	Confirmed	bool		`json:"confirmed"	binding:"required"`
	Phone		string		`json:"phone"		binding:"required"`
	Email		string		`json:"email"		binding:"required"`
	Contact		string		`json:"contact"		binding:"required"`
	Avatar		string		`json:"avatar"		binding:"required"`
}

type UpdateLoginInput struct {
	Username	string		`json:"username"`
	Password	string		`json:"password"`
	UpdatePW	bool		`json:"updatepw"`
	Reset		bool		`json:"reset"`
	Active		bool		`json:"active"`
	Confirmed	bool		`json:"confirmed"`
	Phone		string		`json:"phone"`
	Email		string		`json:"email"`
	Contact		string		`json:"contact"`
	Avatar		string		`json:"avatar"`
}

func hashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
	}
	return string(hash)
}

func comparePasswords(hashPw string, plainPw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPw), []byte(plainPw))
	if err != nil {
		log.Warn(err)
		return false
	}
	return true
}

func Login(c *gin.Context) (interface{}, error) {
	// Validate input
	var input AuthenticateLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	// Get model if exist
	var login models.Login
	if err := models.DB.Where("username = ?", input.Username).First(&login).Error; err != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	matching := comparePasswords(login.Password, input.Password)
	if !matching {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed."})
		log.Info("Password wrong for ", input.Username)
		return nil, jwt.ErrFailedAuthentication
	}
	return &login, nil
}

func GetLogins(c *gin.Context) {
	var logins []models.Login
	_start := c.DefaultQuery("_start", "0")
	_end := c.DefaultQuery("_end", "10")
	_sortOrder := c.DefaultQuery("_sort", "id") + " " + c.DefaultQuery("_order", "ASC")
	result := models.DB.Where("id BETWEEN ? +1 AND ?", _start, _end).Order(_sortOrder).Find(&logins)
	if result.Error != nil {
		c.AbortWithStatus(500)
		log.Warn("Get logins failed.")
	} else {
		c.Header("Access-Control-Expose-Headers", "X-Total-Count")
		c.Header("X-Total-Count", strconv.FormatInt(totalLogin, 10))
		c.JSON(http.StatusOK, logins)
	}
}

func GetLogin(c *gin.Context) {
	// Get model if exist
	var login models.Login
	if err := models.DB.Where("id = ?", c.Param("id")).First(&login).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		log.Warn("Get login failed.")
		return
	}
	c.JSON(http.StatusOK, login)
}

func PostLogin(c *gin.Context) {
	// Validate input
	var input CreateLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Post login failed.", err)
		return
	}
	if len(input.Password) < 5 { //TODO check this in config model
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters."})
		return
	}
	// Create login
	login := models.Login{
		Username: input.Username,
		Password: hashAndSalt(input.Password),
		Reset: input.Reset,
		Active: input.Active,
		Confirmed: input.Confirmed,
		Phone: input.Phone,
		Email: input.Email,
		Contact: input.Contact,
		Avatar: input.Avatar,
	}
	models.DB.Create(&login)
	c.JSON(http.StatusOK, login)
	totalLogin+=1
}

func PutLogin(c *gin.Context) {
	// Validate input
	var input models.Login
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Warn("Put login failed.")
		return
	}
	if input.UpdatePW {
		if len(input.Password) < 5{ //TODO check this in config model
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters."})
			return
		}
		input.Password = hashAndSalt(input.Password)
	}
	// Put Login
	models.DB.Save(&input)
	c.JSON(http.StatusOK, input)
}

func PatchLogin(c *gin.Context) {
	// Get model if exist
	var login models.Login
	if err := models.DB.Where("id = ?", c.Param("id")).First(&login).Error; err != nil {
		log.Warn("Patch login failed.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	log.Debug("id: ", login.ID)
	// Validate input
	var input UpdateLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Warn("Past login failed.")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	patch := models.Login{
		Username: input.Username,
		Password: hashAndSalt(input.Password),
		Reset: input.Reset,
		Active: input.Active,
		Confirmed: input.Confirmed,
		Phone: input.Phone,
		Email: input.Email,
		Contact: input.Contact,
		Avatar: input.Avatar,
	}
	log.Debug("ok")
	//input.Password = hashAndSalt(input.Password)
	log.Debug("here")
	models.DB.Model(&login).Updates(patch)
	log.Debug("wat")
	c.JSON(http.StatusOK, login)
}

func DeleteLogin(c *gin.Context) {
	// Get model if exist
	var login models.Login
	if err := models.DB.Where("id = ?", c.Param("id")).First(&login).Error; err != nil {
		log.Warn("Delete login failed.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	models.DB.Delete(&login)
	c.JSON(http.StatusOK, true)
}

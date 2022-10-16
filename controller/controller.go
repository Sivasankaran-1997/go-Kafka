package controller

import (
	"kafka/domain"
	"kafka/middleware"
	"kafka/services"
	"kafka/utils"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var newUser domain.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	

	result, resterr := services.CreateUser(newUser)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func Login(c *gin.Context) {
	var newUser domain.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.LoginUser(newUser)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func GetUserByEmail(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")

	if strings.TrimSpace(tokenString) == "" {
		resterr := utils.BadRequest("Empty Token Param")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := middleware.UsergetValidate(tokenString)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func DeleteUserByEmail(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")

	if strings.TrimSpace(tokenString) == "" {
		resterr := utils.BadRequest("Empty Token Param")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := middleware.UserdeleteValidate(tokenString)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	} else {
		getResult, err := services.DeleteUser(result.(string))

		if err != nil {
			c.JSON(resterr.Status, resterr)
			return
		}
		c.JSON(http.StatusOK, getResult)
	}
}

func UpdateUserByEmail(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	var newUser domain.User
	if strings.TrimSpace(tokenString) == "" {
		resterr := utils.BadRequest("Empty Token Param")
		c.JSON(resterr.Status, resterr)
	}

	if err := c.ShouldBindJSON(&newUser); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := middleware.UserupdateValidate(tokenString)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	} else {
		newUser.Email = result.(string)
		isParital := c.Request.Method == http.MethodPatch
		result, resterr := services.UpdateUser(isParital, newUser)

		if resterr != nil {
			c.JSON(resterr.Status, resterr)
			return
		}
		c.JSON(http.StatusOK, result)
	}

}

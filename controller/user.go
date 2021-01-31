package controller

import (
	"kratos/logic"
	"kratos/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	var p models.ParamSignUp
	err := c.ShouldBindJSON(&p)
	if err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": errs.Translate(trans),
		})
		return
	}
	if err := logic.SignUp(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Sign Up failed",
			// "err": err.Error(),
		})
		return
	}
	zap.L().Info("User Sign Up", zap.String("user", p.Username))
	c.JSON(http.StatusOK, gin.H{
		"msg": "OK",
	})
}

func LoginHandler(c *gin.Context) {
	var p models.ParamLogin
	err := c.ShouldBindJSON(&p)
	if err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": errs.Translate(trans),
		})
		return
	}
	if err = logic.Login(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	zap.L().Info("User Login", zap.String("user", p.Username))
	c.JSON(http.StatusOK, gin.H{
		"msg": "Login OK",
	})
}

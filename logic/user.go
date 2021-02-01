package logic

import (
	"errors"
	"kratos/dao/mysql"
	"kratos/models"
	"kratos/pkg/jwt"
	"kratos/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) error {
	exist, err := mysql.CheckUserExist(p.Username)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("user existed")
	}
	userID := snowflake.GenID()
	user := models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	err = mysql.InsertUser(&user)
	if err != nil {
		return err
	}
	return nil
}

func Login(p *models.ParamLogin) (token string, err error) {
	user := models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err = mysql.Login(&user); err != nil {
		return "", err
	}
	if token, err = jwt.GenToken(user.UserID, user.Username); err != nil {
		return "", err
	}
	return token, nil
}

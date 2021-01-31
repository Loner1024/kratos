package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"kratos/models"
)

func CheckUserExist(username string) (bool, error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return false, err
	}
	return count > 0, nil
}

func InsertUser(user *models.User) (err error) {
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, encryptPassword(user.Password))
	if err != nil {
		return err
	}
	return nil
}

func Login(user *models.User) error {
	exist, err := CheckUserExist(user.Username)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("user not existed")
	}
	sqlStr := `select username,password from user where username = ?`
	var userinfo models.User
	err = db.Get(&userinfo, sqlStr, user.Username)
	if err != nil {
		return err
	}
	if userinfo.Password != encryptPassword(user.Password) {
		return errors.New("password wrong")
	}
	return nil
}

func encryptPassword(data string) string {
	h := md5.New()
	return hex.EncodeToString(h.Sum([]byte(data)))
}

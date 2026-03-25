package db_model

import (
	"log/slog"

	"time"
)

type UserModel struct {
	UserId     uint32
	UpdateTime int64
	Name       string
}

func SaveUserData() *UserModel {
	slog.Debug("SaveUserData IN=================>222")

	time.Sleep(time.Second * 2) // 模拟db操作需要2秒时间
	user := &UserModel{1, time.Now().Unix(), "14March"}

	slog.Debug("SaveUserData OUT==================>222")
	return user
}

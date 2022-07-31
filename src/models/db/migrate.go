package db

import (
	connect "github.com/hrs-o/docker-go/config"
	entity "github.com/hrs-o/docker-go/models/entity"
)

func Open() {
	v2_db := connect.Dbconnect()

	// Usersというテーブルを作ります
	usertable := v2_db.Migrator().HasTable(&entity.User{})
	if !usertable {
		v2_db.Migrator().CreateTable(&entity.User{})
	}
}

package entity

import "time"

// DB上のテーブル、カラムと構造体との関連付けが自動的に行われる

type User struct {
	Id             int
	Name           string
	Email          string
	Password       string
	Remember_token string
	Created_at     time.Time
	Updated_at     time.Time
}

package model

type CreateInfo struct {
	Usernames []string `db:"users"`
	Id        int64    `db:"id"`
}

type DeleteInfo struct {
	Id int64 `db:"id"`
}

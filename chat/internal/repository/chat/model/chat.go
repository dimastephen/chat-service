package model

import (
	"database/sql"
	"time"
)

type NewChat struct {
	Id         int
	Usernames  []string
	Created_at time.Time
	Updatet_at sql.NullTime
}

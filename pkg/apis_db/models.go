package apis_db

import "time"

type Entry struct {
	ID        string `gorm:"type:uuid;primary_key"`
	Name      string
	CreatedAt time.Time
}

package models

import "time"

type BWListItem struct {
	ID        int       `db:"id"`
	Cidr      string    `db:"cidr"`
	CreatedAt time.Time `db:"created_at"`
}

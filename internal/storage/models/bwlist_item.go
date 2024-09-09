package models

type BWListItem struct {
	ID        int    `db:"id"`
	Cidr      string `db:"cidr"`
	CreatedAt string `db:"created_at"`
}

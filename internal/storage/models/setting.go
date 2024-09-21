package models

type Setting struct {
	ID            int `db:"id"`
	LoginCount    int `db:"login_count"`
	PasswordCount int `db:"password_count"`
	IPCount       int `db:"ip_count"`
}

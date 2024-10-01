package models

type CheckLoginAttemptRequest struct {
	Login    string
	Password string
	IP       string
}

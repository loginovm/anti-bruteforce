package http

import "github.com/gorilla/mux"

func (s *Server) AddRouting(r *mux.Router) {
	r.Path("/check-login-attempt").
		HandlerFunc(s.CheckLoginAttempt).
		Methods("PUT")

	r.Path("/reset-ip").
		HandlerFunc(s.ResetIPBucket).
		Methods("PUT")

	r.Path("/reset-login").
		HandlerFunc(s.ResetLoginBucket).
		Methods("PUT")

	r.Path("/blacklist").
		HandlerFunc(s.GetBlackList).
		Methods("GET")

	r.Path("/whitelist").
		HandlerFunc(s.GetWhiteList).
		Methods("GET")

	r.Path("/blacklist").
		HandlerFunc(s.AddIPToBlackList).
		Methods("POST")

	r.Path("/whitelist").
		HandlerFunc(s.AddIPToWhiteList).
		Methods("POST")

	r.Path("/blacklist").
		HandlerFunc(s.DeleteIPFromBlackList).
		Methods("DELETE")

	r.Path("/whitelist").
		HandlerFunc(s.DeleteIPFromWhiteList).
		Methods("DELETE")

	r.Path("/settings").
		HandlerFunc(s.GetSettings).
		Methods("GET")

	r.Path("/settings").
		HandlerFunc(s.UpdateSettings).
		Methods("PUT")
}

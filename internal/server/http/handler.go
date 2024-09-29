package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/loginovm/anti-bruteforce/internal/app"
	"github.com/loginovm/anti-bruteforce/internal/server/http/models"
	store "github.com/loginovm/anti-bruteforce/internal/storage/models"
)

func (s *Server) CheckLoginAttempt(w http.ResponseWriter, r *http.Request) {
	var req models.CheckLoginAttemptRequest
	if !decodeRequestBody(w, r.Body, &req) {
		return
	}
	ctx := r.Context()
	isAttemptValid, err := s.app.CheckLoginAttempt(ctx, req.Login, req.Password, req.IP)
	if err != nil {
		setErrResponse(err, w)
		return
	}
	ok(w, models.CheckAttemptResponse{IsValid: isAttemptValid})
}

func (s *Server) ResetIPBucket(w http.ResponseWriter, r *http.Request) {
	var req models.IPRequest
	if !decodeRequestBody(w, r.Body, &req) {
		return
	}
	if len(req.IP) > 0 {
		s.app.ResetIPBucket(req.IP)
	}
}

func (s *Server) ResetLoginBucket(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if !decodeRequestBody(w, r.Body, &req) {
		return
	}
	if len(req.Login) > 0 {
		s.app.ResetLoginBucket(req.Login)
	}
}

func (s *Server) GetBlackList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	list, err := s.app.GetBList(ctx)
	if err != nil {
		setErrResponse(err, w)
		return
	}
	resp := models.BWListResponse{Items: list}
	ok(w, resp)
}

func (s *Server) GetWhiteList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	list, err := s.app.GetWList(ctx)
	if err != nil {
		setErrResponse(err, w)
		return
	}
	resp := models.BWListResponse{Items: list}
	ok(w, resp)
}

func (s *Server) AddIPToBlackList(w http.ResponseWriter, r *http.Request) {
	var req models.CidrRequest
	if !decodeRequestBody(w, r.Body, &req) {
		return
	}
	ctx := r.Context()
	if err := s.app.AddIPToBlackList(ctx, req.Cidr); err != nil {
		setErrResponse(err, w)
		return
	}
}

func (s *Server) AddIPToWhiteList(w http.ResponseWriter, r *http.Request) {
	var req models.CidrRequest
	if !decodeRequestBody(w, r.Body, &req) {
		return
	}
	ctx := r.Context()
	if err := s.app.AddIPToWhiteList(ctx, req.Cidr); err != nil {
		setErrResponse(err, w)
		return
	}
}

func (s *Server) DeleteIPFromBlackList(w http.ResponseWriter, r *http.Request) {
	var req models.CidrRequest
	if !decodeRequestBody(w, r.Body, &req) {
		return
	}
	ctx := r.Context()
	if err := s.app.DeleteIPFromBlackList(ctx, req.Cidr); err != nil {
		setErrResponse(err, w)
		return
	}
}

func (s *Server) DeleteIPFromWhiteList(w http.ResponseWriter, r *http.Request) {
	var req models.CidrRequest
	if !decodeRequestBody(w, r.Body, &req) {
		return
	}
	ctx := r.Context()
	if err := s.app.DeleteIPFromWhiteList(ctx, req.Cidr); err != nil {
		setErrResponse(err, w)
		return
	}
}

func (s *Server) GetSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	settings, err := s.app.GetSettings(ctx)
	if err != nil {
		setErrResponse(err, w)
		return
	}
	model := models.Setting{
		IPCount:       settings.IPCount,
		LoginCount:    settings.LoginCount,
		PasswordCount: settings.PasswordCount,
	}
	ok(w, model)
}

func (s *Server) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var req models.Setting
	if !decodeRequestBody(w, r.Body, &req) {
		return
	}
	ctx := r.Context()
	dbSettings, err := s.app.GetSettings(ctx)
	if err != nil {
		setErrResponse(err, w)
		return
	}
	settings := store.Setting{
		ID:            dbSettings.ID,
		IPCount:       req.IPCount,
		LoginCount:    req.LoginCount,
		PasswordCount: req.PasswordCount,
	}
	if err = s.app.UpdateSettings(ctx, settings); err != nil {
		setErrResponse(err, w)
		return
	}
}

func decodeRequestBody(w http.ResponseWriter, body io.ReadCloser, model any) bool {
	err := json.NewDecoder(body).Decode(model)
	if err != nil {
		badRequest(w, "Invalid request body")
		return false
	}

	return true
}

func ok(w http.ResponseWriter, response any) {
	JSON(w, response)
}

func setErrResponse(err error, w http.ResponseWriter) {
	var appErr *app.Error
	if errors.As(err, &appErr) {
		badRequest(w, err.Error())
	} else {
		setError(w)
	}
}

func setError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	r := struct {
		Message string `json:"message"`
	}{Message: message}
	JSON(w, r)
}

func JSON(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

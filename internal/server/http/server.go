package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/loginovm/anti-bruteforce/internal/app"
	"github.com/rs/cors"
)

type Server struct {
	addr       string
	srv        *http.Server
	app        *app.App
	swaggerURL string
}

func NewServer(addr string, app *app.App, swaggerURL string) *Server {
	return &Server{
		addr:       addr,
		app:        app,
		swaggerURL: swaggerURL,
	}
}

func (s *Server) Start(ctx context.Context) error {
	r := mux.NewRouter()
	s.AddRouting(r)

	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedOrigins: []string{s.swaggerURL},
	})
	handler := c.Handler(r)

	server := &http.Server{
		Addr:              s.addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	s.srv = server
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	<-ctx.Done()

	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	return s.srv.Shutdown(ctx)
}

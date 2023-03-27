package application

import (
	"context"
	"net/http"

	repository "github.com/Nickolaygoloburdin/droneapp/internal/database"
	"github.com/Nickolaygoloburdin/droneapp/internal/handlers"
	mw "github.com/Nickolaygoloburdin/droneapp/internal/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type app struct {
	ctx  context.Context
	repo *repository.Repository
}

func (a app) Routes(r *httprouter.Router) {
	dbHandler := handlers.NewDBHandler(a.ctx, a.repo)
	r.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})
	r.GET("/auth/signin", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		//a.LoginPage(rw, "")
	})
	r.POST("/auth/update", mw.Logger(dbHandler.Login, "update"))
	r.POST("/auth/signup", mw.Logger(dbHandler.Login, "signup"))
}

func NewApp(ctx context.Context, dbpool *pgxpool.Pool) *app {
	return &app{ctx, repository.NewRepository(dbpool)}
}

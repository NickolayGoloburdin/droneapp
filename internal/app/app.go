package application

import (
	"context"

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
	r.GET("/auth/update", mw.Logger(handlers.GetDataFromToken, "Token"))

	r.POST("/auth/signup", mw.Logger(dbHandler.Signup, "signup"))
	r.POST("/auth/signin", mw.Logger(dbHandler.Login, "login"))
}

func NewApp(ctx context.Context, dbpool *pgxpool.Pool) *app {
	return &app{ctx, repository.NewRepository(dbpool)}
}

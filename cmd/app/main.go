package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	application "github.com/Nickolaygoloburdin/droneapp/internal/app"
	repository "github.com/Nickolaygoloburdin/droneapp/internal/database"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	ctx := context.Background()
	dbpool, err := repository.InitDBConn(ctx)
	if err != nil {
		log.Fatalf("%w failed to init DB connection", err)
	}
	defer dbpool.Close()
	a := application.NewApp(ctx, dbpool)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})
	r := httprouter.New()
	router := c.Handler(r)
	a.Routes(r)
	srv := &http.Server{Addr: "192.168.2.135:8000", Handler: router}
	fmt.Println("It is alive! Try http://192.168.2.135:8000")
	srv.ListenAndServe()
}

package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"mohaafaridd.dev/ecom/services/user"
)

type APIServer struct {
	address string
	db      *sql.DB
}

func NewAPIServer(address string, db *sql.DB) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
	}
}

func (server *APIServer) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(server.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subRouter)

	log.Println("Listening on", server.address)
	return http.ListenAndServe(server.address, router)
}

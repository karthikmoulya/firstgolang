package router

import (
	"github.com/gorilla/mux"
	"github.com/karthikmoulya/mongoapi/controller"
)

func Router() *mux.Router{
	router := mux.NewRouter()

	router.HandleFunc("/api/movies", controller.GetMyAllMovies).Methods("GET")
	router.HandleFunc("/api/movie", controller.CreateMovie).Methods("POST")
	router.HandleFunc("/api/movie/{id}", controller.MarkasWatched).Methods("PUT")
	router.HandleFunc("/api/movie/{id}", controller.DeleteMyOneMovie).Methods("DELETE")
	router.HandleFunc("/api/deleteallmovie", controller.DeleteAllMyMovies).Methods("DELETE")

	return router
}
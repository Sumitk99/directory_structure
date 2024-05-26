package routes

import (
	"directory_structure_api/pkg/controllers"
	"github.com/gorilla/mux"
)

func GenerateRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", controllers.Read_func).Methods("GET")
	router.HandleFunc("/", controllers.Delete_func).Methods("DELETE")
	router.HandleFunc("/", controllers.Create_func).Methods("POST")
	router.HandleFunc("/", controllers.Update_func).Methods("PUT")

	return router
}

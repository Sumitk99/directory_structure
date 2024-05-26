package main

import (
	_ "directory_structure_api/pkg/models"
	"directory_structure_api/pkg/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := routes.GenerateRoutes()

	fmt.Println("Starting server on port 9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}

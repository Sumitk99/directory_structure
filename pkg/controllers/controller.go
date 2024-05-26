package controllers

import (
	"directory_structure_api/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func Read_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var InputPath models.Input
	json.NewDecoder(req.Body).Decode(&InputPath)
	fmt.Println(InputPath)
	json.NewEncoder(res).Encode(models.ListItems(InputPath.Path))
}

func Create_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var InputPath models.Input
	json.NewDecoder(req.Body).Decode(&InputPath)
	models.CreatePath(InputPath.Path)
	json.NewEncoder(res).Encode(InputPath)
}

func Update_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var InputPath models.Input
	json.NewDecoder(req.Body).Decode(&InputPath)
	models.UpdatePath(InputPath.Path)
	json.NewEncoder(res).Encode(models.ListItems("/"))
}

func Delete_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var InputPath models.Input
	json.NewDecoder(req.Body).Decode(&InputPath)
	models.DeletePath(InputPath.Path)
}

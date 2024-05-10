package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var graph = make(map[string][]node)

var root_node = node{Name: "root"}
var current_directory = &root_node

func parse(s string, index int) (string, int) {
	t := []rune(s)
	start := index
	ans := make([]rune, 0)
	for i := index; i < len(s) && t[i] != '/'; i++ {
		index++
	}

	return string(append(ans, t[start:index]...)), index
}

type node struct {
	Name   string `json:"name"`
	parent *node
}

func (n node) insert(s string) {
	start := 1
	fmt.Println("inserting\n")
	for start < len(s) {
		partial_path, end := parse(s, start)
		new_path := node{Name: partial_path, parent: &n}
		graph[n.Name] = append(graph[n.Name], new_path)
		fmt.Printf("%s inserted with parent node %s\n", graph[n.Name][len(graph[n.Name])-1].Name, n.Name)
		n = new_path
		start = end + 1
	}
}

type path struct {
	S string `json:"path"`
}

func read_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(graph[current_directory.Name])
}

func create_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var p path
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	current_directory.insert(p.S)
	json.NewEncoder(res).Encode(graph[current_directory.Name])
}

func compare(s string) bool {
	if len(s) < 5 {
		return false
	}
	t := []rune(s)
	cd := "cd .."
	cds := []rune(cd)
	for i := range t {
		if t[i] != cds[i] {
			return false
		}
	}
	return true
}
func update_func(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Update started\n")
	res.Header().Set("Content-Type", "application/json")
	var p path
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	start := 1
	if compare(p.S) {
		current_directory = current_directory.parent
	} else {
		for start < len(p.S) {
			partial_path, end := parse(p.S, start)
			for j, k := range graph[current_directory.Name] {
				if k.Name == partial_path {
					current_directory = &graph[current_directory.Name][j]
					break
				}
			}
			fmt.Println("supposed to go %s and here it is %s\n", partial_path, current_directory.Name)
			start = end + 1
		}
	}
	json.NewEncoder(res).Encode(graph[current_directory.Name])
}

//func delete_func(res http.ResponseWriter, req *http.Request) {
//	res.Header().Set("Content-Type", "application/json")
//	var message Message
//	err := json.NewDecoder(req.Body).Decode(&message)
//	if err != nil {
//		http.Error(res, err.Error(), http.StatusBadRequest)
//		return
//	}
//	db[message.Key] = message.Value
//	delete(db, message.Key)
//	json.NewEncoder(res).Encode(db)
//}

func main() {
	router := mux.NewRouter()
	graph[root_node.Name] = append(graph[root_node.Name], node{})

	root_node.insert("/dcim/photos")
	for _, j := range graph[root_node.Name] {
		fmt.Println(j.Name)
	}
	router.HandleFunc("/", read_func).Methods("GET")
	//	router.HandleFunc("/delete", delete_func).Methods("DELETE")
	router.HandleFunc("/", create_func).Methods("POST")
	router.HandleFunc("/", update_func).Methods("PUT")

	fmt.Println("Starting server on port 9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var graph = make(map[string][]node)
var stack = make([]*node, 0)
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

func (n node) insert(s string, index int) {
	start := index
	fmt.Println("inserting\n")
	for start < len(s) {
		partial_path, end := parse(s, start)
		new_path := node{Name: partial_path, parent: &n}
		graph[n.Name] = append(graph[n.Name], new_path)
		fmt.Printf("%s inserted with parent node %s\n", graph[n.Name][len(graph[n.Name])-1].Name, graph[n.Name][len(graph[n.Name])-1].parent.Name)
		n = new_path
		start = end + 1
	}
}

func navigate(s string, index int, temp_directory *node) *node {
	in := index
	start := index
	for start < len(s) {
		partial_path, end := parse(s, start)
		for j, k := range graph[temp_directory.Name] {
			if k.Name == partial_path {
				temp_directory = &graph[temp_directory.Name][j]
				if in == 1 {
					stack = append(stack, temp_directory)
				}
				break
			}
		}
		fmt.Printf("supposed to go %s and here it is %s\n", partial_path, current_directory.Name)
		start = end + 1
	}
	return temp_directory

}

type path struct {
	S string `json:"path"`
}

func check_path(s string, root_directory *node) int {

	for i, child := range graph[root_directory.Name] {
		if child.Name == s {
			return i
		}
	}
	return -1
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
	temp_directory := current_directory
	start := 1
	for start < len(p.S) {
		partial_path, end := parse(p.S, start)
		if check_path(partial_path, temp_directory) != -1 {
			temp_directory = navigate(partial_path, 0, temp_directory)
		} else {
			break
		}
		start = end + 1
	}
	temp_directory.insert(p.S, start)
	json.NewEncoder(res).Encode(graph[current_directory.Name])
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
	temp_directory := current_directory
	if p.S == ".." {
		stack = stack[:len(stack)-1]
		fmt.Printf("cur - %s, parent - %s\n", current_directory.Name, current_directory.parent.Name)
		current_directory = stack[len(stack)-1]
	} else {
		current_directory = navigate(p.S, 1, temp_directory)
	}
	json.NewEncoder(res).Encode(graph[current_directory.Name])
}

func delete_func(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var input_path path
	err := json.NewDecoder(req.Body).Decode(&input_path)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	start := 1
	temp_directory := current_directory
	for start < len(input_path.S) {
		partial_path, end := parse(input_path.S, start)
		if end >= len(input_path.S) {
			for i, j := range graph[temp_directory.Name] {
				if j.Name == partial_path {
					graph[temp_directory.Name] = append(graph[temp_directory.Name][:i], graph[temp_directory.Name][i+1:]...)
					break
				}
			}
			break
		} else {
			temp_directory = navigate(partial_path, start, temp_directory)
			start = end + 1
		}
	}
	json.NewEncoder(res).Encode(graph[current_directory.Name])
}

func main() {
	router := mux.NewRouter()
	graph[root_node.Name] = append(graph[root_node.Name], node{})
	stack = append(stack, &root_node)
	root_node.insert("/dcim/photos", 1)
	for _, j := range graph[root_node.Name] {
		fmt.Println(j.Name)
	}
	router.HandleFunc("/", read_func).Methods("GET")
	router.HandleFunc("/", delete_func).Methods("DELETE")
	router.HandleFunc("/", create_func).Methods("POST")
	router.HandleFunc("/", update_func).Methods("PUT")

	fmt.Println("Starting server on port 9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}

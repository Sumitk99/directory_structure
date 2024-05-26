package models

import (
	"fmt"
)

type Node struct {
	Name     string  `json:"name"`
	Parent   *Node   `json:"parent"`
	Children []*Node `json:"children"`
}

type Input struct {
	Path string `json:"path"`
}

type Output struct {
	Current string   `json:"current directory"`
	Showing string   `json:"showing"`
	Items   []string `json:"items"`
}

var CurrentDirectory *Node
var RootDirectory *Node

func init() {
	RootDirectory = &Node{Name: "root", Parent: nil}
	CurrentDirectory = RootDirectory
	CreatePath("/Android/Data")
	CreatePath("/whatsapp/doc/images")
}

func FindFolder(s string, root *Node) *Node {
	for _, it := range root.Children {
		if it.Name == s {
			return it
		}
	}
	fmt.Println("Child not found, adding...")
	return nil
}

func Navigate(root *Node, s string) *Node {
	start := 1
	var cur_folder string
	temp_directory := root

	for i := start; i < len(s); i++ {

		if s[i] == '/' || i == len(s)-1 {
			x := 0
			if i == len(s)-1 {
				x = 1
			}
			cur_folder = s[start : i+x]
			fmt.Printf("in navigate : %s\n", cur_folder)
			if cur_folder == ".." {
				temp_directory = temp_directory.Parent
			} else {
				temp_directory = FindFolder(cur_folder, temp_directory)
			}
			start = i + 1
		}
	}
	return temp_directory
}

func Insert(s string, root *Node) {
	fmt.Printf("Inserting %s\n root directory given : %s\n", s, root.Name)
	temp_directory := root

	start := 1
	for i := 1; i < len(s); i++ {
		if s[i] == '/' || i == len(s)-1 {
			x := 0
			if i == len(s)-1 {
				x = 1
			}
			temp_directory.Children = append(temp_directory.Children, &Node{Name: s[start : i+x], Parent: temp_directory})
			temp_directory = temp_directory.Children[len(temp_directory.Children)-1]
			start = i + 1
		}
	}
}

func ListItems(s string) Output {
	var Response Output
	Response.Current = CurrentDirectory.Name
	temp_directory := Navigate(CurrentDirectory, s)
	Response.Showing = temp_directory.Name
	for _, it := range temp_directory.Children {
		Response.Items = append(Response.Items, it.Name)
	}
	return Response
}

func CreatePath(s string) {
	fmt.Println("Creating Path")
	start := 1
	var cur_folder string
	temp_directory := CurrentDirectory
	for i := start; i < len(s); i++ {
		if s[i] == '/' || i == len(s)-1 {
			x := 0
			if i == len(s)-1 {
				x = 1
			}
			cur_folder = s[start : i+x]
			fmt.Printf("in create : %s\n", cur_folder)
			if cur_folder == ".." {
				temp_directory = temp_directory.Parent
			} else {
				temp := FindFolder(cur_folder, temp_directory)
				if temp != nil {
					temp_directory = temp
					fmt.Println("temp is not nill : %s\n", temp.Name)
				} else {
					Insert(s[start-1:], temp_directory)
					break
				}
			}
			start = i + 1
		}
	}
}

func UpdatePath(s string) {
	CurrentDirectory = Navigate(CurrentDirectory, s)
}

func DeletePath(s string) {
	fmt.Printf("Deleting Path : %s\n", s)
	var index int
	t := []rune(s)
	temp_directory := Navigate(CurrentDirectory, s)
	fmt.Printf("Into the folder : %s\n", temp_directory.Name)
	temp_directory = temp_directory.Parent
	fmt.Printf("It's Parent : %s\n", temp_directory.Name)
	for i := len(s) - 1; i >= 0; i-- {
		if '/' == t[i] {
			index = i
			break
		}
	}
	target_directory := s[index+1:]
	fmt.Printf("target : %s\n", target_directory)
	for i, it := range temp_directory.Children {
		if it.Name == target_directory {
			temp_directory.Children = append(temp_directory.Children[:i], temp_directory.Children[i+1:]...)
		}
	}

}

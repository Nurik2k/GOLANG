package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type TreeNode struct {
	File       os.FileInfo
	fileInFile []TreeNode
}

func (n TreeNode) Name() string {
	if !n.File.IsDir() {
		return fmt.Sprintf("%s (%v)", n.File.Name(), n.Size())
	}

	return n.File.Name()
}

func (n TreeNode) Size() string {
	if n.File.Size() != 0 {
		return fmt.Sprintf("%db", n.File.Size())
	}

	return "empty"
}

func GetNodes(path string, withFiles bool) ([]TreeNode, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var nodes []TreeNode
	for _, file := range files {
		if !withFiles && !file.IsDir() {
			continue
		}

		node := TreeNode{
			File: file,
		}
		if file.IsDir() {
			fileInFile, err := GetNodes(path+string(os.PathSeparator)+file.Name(), withFiles)
			if err != nil {
				return nil, err
			}

			node.fileInFile = fileInFile
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func PrintTree(out io.Writer, nodes []TreeNode, parentPrefix string) {
	var (
		lastIndex = len(nodes) - 1
		prefix    = "├───"
		_prefix   = "│\t"
	)

	for i, node := range nodes {
		if i == lastIndex {
			prefix = "└───"
			_prefix = "\t"
		}
		fmt.Fprint(out, parentPrefix, prefix, node.Name(), "\n")

		if node.File.IsDir() {
			PrintTree(out, node.fileInFile, parentPrefix+_prefix)
		}
	}
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	nodes, err := GetNodes(path, printFiles)
	if err != nil {
		return
	}

	PrintTree(out, nodes, "")

	return
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

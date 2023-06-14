package GOLANG

import (
	"fmt"
	"io/ioutil"
	"log"
)

var path = "S://HomeWorks_GO"

func main() {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}

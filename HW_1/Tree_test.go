package main

import (
	"bytes"
	"testing"
)

func TestDirTree(t *testing.T) {
	// Создаем буфер для записи вывода
	var buf bytes.Buffer

	// Запускаем функцию DirTree с выводом в буфер
	err := DirTree(&buf, ".", true)
	if err != nil {
		t.Fatalf("DirTree failed: %v", err)
	}

	// Ожидаемый вывод
	expectedOutput := `├───Tree.go (1797b)
        ├───Tree_test.go (2474b)
        ├───go.mod (23b)
        └───hw1_tree
                ├───dockerfile (75b)
                ├───hw1.md (4621b)
                ├───main.go (352b)
                ├───main_test.go (1865b)
                └───testdata
                        ├───project
                        |       ├───file.txt (19b)
                        |       └───gopher.png (70372b)
                        ├───static
                        |       ├───a_lorem
                        |       |       ├───dolor.txt (0)
                        |       |       ├───gopher.png (70372b)
                        |       |       └───ipsum
                        |       |               └───gopher.png (70372b)
                        |       ├───css
                        |       |       └───body.css (28b)
                        |       ├───empty.txt (0)
                        |       ├───html
                        |       |       └───index.html (57b)
                        |       ├───js
                        |       |       └───site.js (10b)
                        |       └───z_lorem
                        |               ├───dolor.txt (0)
                        |               ├───gopher.png (70372b)
                        |               └───ipsum
                        |                       └───gopher.png (70372b)
                        ├───zline
                        |       ├───empty.txt (0)
                        |       └───lorem
                        |               ├───dolor.txt (0)
                        |               ├───gopher.png (70372b)
                        |               └───ipsum
                        |                       └───gopher.png (70372b)
                        └───zzfile.txt (0)`

	// Проверяем, что полученный вывод совпадает с ожидаемым выводом
	if buf.String() != expectedOutput {
		t.Errorf("DirTree output doesn't match expected output:\n\nExpected:\n%s\n\nGot:\n%s", expectedOutput, buf.String())
	}
}

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
	expectedOutput := `
├───Tree.go (1797b)
├───Tree_test.go (759b)
└───go.mod (23b)
`

	// Проверяем, что полученный вывод совпадает с ожидаемым выводом
	if buf.String() != expectedOutput {
		t.Errorf("DirTree output doesn't match expected output:\n\nExpected:\n%s\n\nGot:\n%s", expectedOutput, buf.String())
	}
}

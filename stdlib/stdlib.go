package stdlib

import (
	_ "embed"
)

//go:embed heap.mr
var heapLib string

//go:embed string.mr
var stringLib string

//go:embed input.mr
var inputLib string

func LoadStdLibFiles() string {
	libNames := [3]string{
		heapLib,
		stringLib,
		inputLib,
	}

	libString := ""

	for _, fileContent := range libNames {
		libString += "\n" + string(fileContent)
	}

	return libString
}

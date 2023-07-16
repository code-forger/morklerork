package stdlib

import (
	"log"
	"os"
)

func LoadStdLibFiles() string {
	libNames := [1]string{"./stdlib/heap.mr"}

	libString := ""

	for _, name := range libNames {
		content, err := os.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}
		libString += "\n" + string(content)
	}

	return libString
}
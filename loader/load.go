package loader

import (
	"log"
	"os"
)

func Load() string {
	programName := os.Args[1]

	content, err := os.ReadFile(programName)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

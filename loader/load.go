package loader

import (
	"log"
	"morklerork/stdlib"
	"os"
)

func Load() string {
	programNames := os.Args[1:]

	programString := ""

	for _, name := range programNames {
		content, err := os.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}
		programString += "\n" + string(content)
	}

	return stdlib.LoadStdLibFiles() + programString
}

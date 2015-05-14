package main

import (
	"fmt"
	"github.com/roxtar/deps/importer"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("package name required as argument")
	}
	imports, err := importer.GetImportsPackage(os.Args[1])
	if err != nil {
		log.Fatalf("Error getting imports: %s", err.Error())
	}

	for _, pkg := range imports {
		fmt.Println(pkg)
	}
}

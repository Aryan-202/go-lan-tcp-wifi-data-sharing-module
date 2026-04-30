package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get cwd: %v", err)
	}
	fmt.Printf("The current working directory is: %s\n", cwd)
}

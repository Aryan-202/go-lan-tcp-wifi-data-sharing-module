package main

import (
	"fmt"
	
	"os"

	"github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module/helpers"
)

func main() {
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf("%-20s %s\n", "File Name", "Size")

	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()

			if err != nil {
				fmt.Printf("Could not get info %s: %v\n", entry.Name(), err)
				continue
			}

			fmt.Printf("%-20s %d\n", entry.Name(), info.Size())
		}
	}

	cleaned := helpers.CleanInput("  Hello   World  ")
	fmt.Println("Cleaned input:", cleaned)
}

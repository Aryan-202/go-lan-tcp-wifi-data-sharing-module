package helpers

import "strings"

func CleanInput(str string) []string {
	exit := strings.ToLower(str)
	words := strings.Fields(exit)
	return  words
}
package core

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Prompt(label string) string {
	fmt.Print(label + ": ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func Select(label string, items []string) int {
	fmt.Println(label)
	for i, item := range items {
		fmt.Printf("%d. %s\n", i+1, item)
	}
	
	for {
		choice := Prompt("Select")
		idx, err := strconv.Atoi(choice)
		if err == nil && idx > 0 && idx <= len(items) {
			return idx - 1
		}
		fmt.Println("Invalid selection")
	}
}

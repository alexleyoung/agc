package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Run() {
	fmt.Println("AGC session started. What would you like to do?")
	for true {
		reader := bufio.NewReader(os.Stdin)
		input, err := getInput(reader)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Echo: %s\n", input)
	}
}

func getInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	return input, nil
}

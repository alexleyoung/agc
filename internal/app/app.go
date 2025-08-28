package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/alexleyoung/agc/internal/ai"
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

		// ai.Chat(context.Background())
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

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func isValidId(id string) bool {
	if id == "" {
		return false
	}
	for _, c := range id {
		if !unicode.In(c, unicode.Latin) {
			return false
		}
	}
	return true
}

func handleAssignment(input string, env map[string]int) {
	parts := strings.Split(input, "=")
	if len(parts) != 2 {
		fmt.Println("Invalid assignment")
	} else {
		varId := strings.TrimSpace(parts[0])
		if !isValidId(varId) {
			fmt.Println("Invalid identifier")
			return
		}

		val := strings.TrimSpace(parts[1])
		intVal, err := strconv.Atoi(val)
		if err != nil {
			if isValidId(val) {
				var ok bool
				intVal, ok = env[val]
				if !ok {
					fmt.Println("Unknown variable")
					return
				}
			} else {
				fmt.Println("Invalid assignment")
				return
			}
		}

		env[varId] = intVal
	}
}

func handleExpression(input string, env map[string]int) {
	sum := 0
	operator := "+"
	invalid := false

	setInvalid := func(message string) {
		fmt.Println(message)
		invalid = true
	}

	for _, chunk := range strings.Split(input, " ") {
		if invalid {
			break
		}
		if operator == "" {
			if chunk[0] == '+' {
				operator = "+"
			} else if chunk[0] == '-' {
				if len(chunk)%2 == 0 {
					operator = "+"
				} else {
					operator = "-"
				}
			} else {
				setInvalid("Invalid expression")
			}
		} else {
			num, err := strconv.Atoi(chunk)
			if err != nil {
				if isValidId(chunk) {
					var ok bool
					num, ok = env[chunk]
					if !ok {
						setInvalid("Unknown variable")
					}
				} else {
					setInvalid("Invalid expression")
				}
			}
			if operator == "+" {
				sum += num
			} else {
				sum -= num
			}
			operator = ""
		}
	}
	if operator != "" {
		setInvalid("Invalid expression")
	}
	if !invalid {
		fmt.Println(sum)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	env := make(map[string]int, 0)
	for line, _ := reader.ReadString('\n'); line != "/exit\n"; line, _ = reader.ReadString('\n') {
		line = strings.TrimSpace(line)
		switch line {
		case "":
			continue
		case "/help":
			fmt.Println("The program calculates the sum of numbers")
		default:
			if strings.HasPrefix(line, "/") {
				fmt.Println("Unknown command")
			} else if strings.Contains(line, "=") {
				handleAssignment(line, env)
			} else {
				handleExpression(line, env)
			}
		}
	}
	fmt.Println("Bye!")

}

package main

import "fmt"

func main() {
	var input string
	fmt.Scanf("%s\n", &input)

	answer := 1
	min, max := 'A', 'Z'
	for _, ch := range input {
		if ch >= min && ch <= max {
			answer++
		}
	}

	fmt.Println(answer)
}

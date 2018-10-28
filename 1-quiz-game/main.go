package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "A CSV file in the format of 'question, answer'")
	timeLimit := flag.Int("limit", 30, "The time limit in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV"))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the given CSV file")
	}

	problems := parseLines(lines)
	correct := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, problem.question)
		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("You scored %d out of %d\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == problem.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	result := make([]problem, len(lines))

	for i, line := range lines {
		result[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return result
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

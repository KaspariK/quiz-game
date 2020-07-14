package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	csvName := flag.String("csv", "problems.csv", "csv file with problems and answers e.g., '5+5,10'")
	quizLength := flag.Int("time", 10, "number of seconds representing the length of the quiz")
	flag.Parse()

	file, err := os.Open(*csvName)
	if err != nil {
		log.Printf("Could not open file: %s", *csvName)
		os.Exit(1)
	}

	defer file.Close()

	r := csv.NewReader(file)

	lines, err := r.ReadAll()
	if err != nil {
		log.Printf("Could not read file: %s", *csvName)
	}

	problems := parseLines(lines)

	quizTimer := time.NewTimer(time.Duration(*quizLength) * time.Second)

	correct := 0

	for _, problem := range problems {
		fmt.Printf("%s = ", problem.question)
		answerCh := make(chan string)

		go func() {
			var userAnswer string

			fmt.Scan(&userAnswer)
			answerCh <- userAnswer
		}()

		select {
		case <-quizTimer.C:
			fmt.Printf("\nTime is up! You scored %d out of %d", correct, len(problems))
			return
		case userAnswer := <-answerCh:
			if userAnswer == problem.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))

	for i, line := range lines {
		problems[i] = problem{
			question: line[0],
			answer:   line[1],
		}
	}

	return problems
}

type problem struct {
	question string
	answer   string
}

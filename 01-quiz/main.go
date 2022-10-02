package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   int
}

func main() {
	csvFilename := flag.String("csv", "promlems.csv", "a csv file in the format of 'question'answer'")
	timeLimit := flag.Int("limit", 30, "the limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}
	problems := parseLines(lines)

	timer := time.NewTicker(time.Duration(*timeLimit) * time.Second)
	correct := 0

loop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)
		answerCh := make(chan int)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- parseAnswer(answer)
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break loop
		case answer := <-answerCh:
			fmt.Println(answer)
			if answer == p.answer {
				correct++
			}
		}

	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   parseAnswer(line[1]),
		}
	}
	return ret
}

func parseAnswer(answer string) int {
	trimed := strings.Trim(answer, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+~{}[]:\"';-=,.<>/\\|`")
	ret, _ := strconv.Atoi(trimed)
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

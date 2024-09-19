package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const timeToAnswer = 60 * time.Second

func main() {
	lines, err := readCSV()
	if err != nil {
		fmt.Println(err)
		return
	}

	problems := parseLines(lines)

	timer := time.NewTimer(timeToAnswer)
	correct := 0
	total := 0

	for i := 0; i < len(problems); i++ {
		p := randomQuestion(problems)
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scan(&answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nTime's up! You scored %d out of %d.\n", correct, total)
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
		total++
	}
	fmt.Printf("You scored %d out of %d.\n", correct, total)
}

func randomQuestion(p []problem) problem {
	random := rand.Intn(len(p))
	return p[random]
}

type problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{q: line[0], a: strings.TrimSpace(line[1])}
	}
	return ret
}

func readCSV() ([][]string, error) {
	csvFilename := flag.String("csv", "problems.csv", "a CSV file in the format of 'question,answer'")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		return nil, fmt.Errorf("Failed to open the CSV file: %s", *csvFilename)
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the provided CSV file.")
	}

	return lines, nil
}

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in format of 'question,answer'")
	flag.Parse()

	file, err := os.Open(*csvFileName)

	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file: %s\n", *csvFileName))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit(fmt.Sprintf("Failed to pares the provided CSV file."))
	}

	problems := parseLInes(lines)
	goodAnswers := 0

	for i, p := range problems {
		var answer string
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)

		_, err = fmt.Scanf("%s\n", &answer)

		if err != nil {
			exit("There was an error parsing your answer")
		}

		if answer == p.answer {
			goodAnswers++
		}
	}

	fmt.Printf("Result %d/%d \n", goodAnswers, len(problems))
}

func parseLInes(Lines [][]string) []problem {
	ret := make([]problem, len(Lines))
	for i, line := range Lines {
		ret[i] = problem{
			line[0],
			strings.TrimSpace(line[1]),
		}
	}

	return ret
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}

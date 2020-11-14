package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var goodAnswers = 0

func main() {
	flags := initFlags()
	lines := readFile(flags.csvFile)
	problems := parseLInes(lines)

	prepareToStart()

	startQuiz(problems, flags)
}

func readFile(filepath string) [][]string {
	file, err := os.Open(filepath)

	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file: %s\n", filepath))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit(fmt.Sprintf("Failed to pares the provided CSV file."))
	}

	return lines
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

func startQuiz(problems []problem, flags flags) {
	timer := startTimer(flags)
	answerChanel := make(chan string)

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)

		go getAnswer(answerChanel)

		procedeWithQuizAndStopWHenTimerEnds(timer.C, answerChanel, p)
	}

	fmt.Printf("Result %d/%d \n", goodAnswers, len(problems))
}

func procedeWithQuizAndStopWHenTimerEnds(timerChan <-chan time.Time, answerChan <-chan string, problem problem) {
	select {
	case <-timerChan:
		fmt.Println()
		return
	case answer := <-answerChan:
		if problem.validateAnswer(answer) {
			goodAnswers++
		}
	}
}

func prepareToStart() {
	var t string
	fmt.Println("Press enter to start")
	_, _ = fmt.Scanf("%s\n", &t)
}

func getAnswer(answerChan chan<- string) {
	var answer string

	_, err := fmt.Scanf("%s\n", &answer)

	if err != nil {
		exit("There was an error parsing your answer")
	}

	answerChan <- answer
}

func initFlags() flags {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in format of 'question,answer'")
	timeout := flag.Int("timer", 5, "Time after which quiz will end")

	flag.Parse()

	return flags{
		csvFile:   *csvFileName,
		timerTime: *timeout,
	}
}

type problem struct {
	question string
	answer   string
}

func (p problem) validateAnswer(answer string) bool {
	if answer != p.answer {
		return false
	}

	return true
}

func startTimer(flags flags) *time.Timer {
	return time.NewTimer(time.Duration(flags.timerTime) * time.Second)
}

type flags struct {
	csvFile   string
	timerTime int
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}

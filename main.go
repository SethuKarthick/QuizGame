package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problems struct {
	q string
	a string
}

func main() {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of question, answer")
	timeLimit := flag.Int("TimeLimit", 30, "The time Limit for the quiz in seconds")
	flag.Parse()
	_ = csvFileName

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open csv file - %s\n", *csvFileName))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided csv file")
	}
	probs := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
quizloop:
	for i, p := range probs {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		var answerCh = make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			// fmt.Printf("\nyou scored %d out of %d\n", correct, len(probs))
			fmt.Println()
			break quizloop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("you scored %d out of %d\n", correct, len(probs))

}

func parseLines(lines [][]string) []problems {
	ret := make([]problems, len(lines))
	for i, line := range lines {
		ret[i] = problems{q: line[0], a: strings.TrimSpace(line[1])}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

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
	csvFileName, timeLimit := getFlagValues()
	file := openFile(csvFileName)
	records := parseCSV(file)
	problems := parseRecords(records)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	score := 0
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)

		answerCh := make(chan string)
		scanAnswer(answerCh)

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d/%d\n", score, len(problems))

			return
		case answer := <-answerCh:
			if answer == problem.answer {
				score++
			}
		}
	}

	fmt.Printf("You scored %d/%d\n", score, len(problems))
}

func getFlagValues() (*string, *int) {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	return csvFileName, timeLimit
}

func openFile(csvFileName *string) *os.File {
	file, err := os.Open(*csvFileName)
	if err != nil {
		fmt.Printf("Failed to open the CSV file: %s\n", *csvFileName)
		os.Exit(1)
	}

	return file
}

func parseCSV(file *os.File) [][]string {
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Failed to read the provided CSV file")
		os.Exit(1)
	}

	return records
}

func scanAnswer(answerCh chan string) {
	go func() {
		var answer string
		fmt.Scanf("%s\n", &answer)

		answerCh <- answer
	}()
}

type problem struct {
	question string
	answer   string
}

func parseRecords(records [][]string) []problem {
	problems := make([]problem, len(records))

	for i, record := range records {
		problems[i] = problem{
			question: record[0],
			answer:   strings.TrimSpace(record[1]),
		}
	}

	return problems
}

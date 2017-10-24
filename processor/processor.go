package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("File processor")

	if len(os.Args) != 5 {
		log.Fatalln("Usage: processor originFile destinyFile errorFile unprocessedFile")
	}

	originFile, err := os.Open(os.Args[1])
	PrintFatalError(err, "Error opening the from file")
	defer originFile.Close()

	destinyFile, err := os.Create(os.Args[2])
	PrintFatalError(err, "Error creating the destiny file")
	defer destinyFile.Close()

	errorFile, err := os.Create(os.Args[3])
	PrintFatalError(err, "Error creating the error file")
	defer errorFile.Close()

	unprocessedFile, err := os.Create(os.Args[4])
	PrintFatalError(err, "Error creating the unprocessed file")
	defer unprocessedFile.Close()

	scanner := bufio.NewScanner(originFile)
	count := 0
	for scanner.Scan() {
		count++
		_, err := destinyFile.WriteString(scanner.Text() + "\n")
		if err != nil {
			_, err := errorFile.WriteString(err.Error())
			PrintFatalError(err, "Error writing the error of processing the line "+string(count))
		}
	}

	fmt.Println("END File processor")

}

//PrintFatalError Utility to print an error with a message with a log fatal (it exits the program)
func PrintFatalError(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}

}

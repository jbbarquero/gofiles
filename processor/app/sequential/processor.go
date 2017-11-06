package main

import (
	"fmt"
	"log"
	"os"

	proc "github.com/jbbarquero/gofiles/processor"
)

func main() {
	fmt.Println("File processor")

	if len(os.Args) != 5 {
		log.Fatalln("Usage: processor originFile destinyFile errorFile unprocessedFile")
	}

	originFile, err := os.Open(os.Args[1])
	proc.PrintFatalError(err, "Error opening the from file")
	defer originFile.Close()

	destinyFile, err := os.Create(os.Args[2])
	proc.PrintFatalError(err, "Error creating the destiny file")
	defer destinyFile.Close()

	errorFile, err := os.Create(os.Args[3])
	proc.PrintFatalError(err, "Error creating the error file")
	defer errorFile.Close()

	unprocessedFile, err := os.Create(os.Args[4])
	proc.PrintFatalError(err, "Error creating the unprocessed file")
	defer unprocessedFile.Close()

	proc.ProcessSequential(originFile, destinyFile, errorFile)

	fmt.Println("END File processor")

}

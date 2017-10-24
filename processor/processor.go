package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type fileInfo struct {
	id  string
	seq int
}

type fileInfoResult struct {
	id   string
	seqs []int
}

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

	fileInfos := processFile(originFile, errorFile)

	processResult(destinyFile, errorFile, fileInfos)

	fmt.Println("END File processor")

}

func processFile(originFile, errorFile *os.File) []fileInfo {
	fileInfos := []fileInfo{}
	r := csv.NewReader(originFile)
	r.Comma = ';'
	r.Comment = '#'
	count := 0
	//var previous fileInfo
	for {
		count++
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			_, err := errorFile.WriteString(fmt.Sprintf("Error at line %d: %v", count, err))
			PrintFatalError(err, fmt.Sprintf("Error at processing file. Error writing the error of processing the line %d", count))
		}
		seq, err := strconv.Atoi(record[1])
		if err != nil {
			_, err := errorFile.WriteString(fmt.Sprintf("Error at line %d: %v", count, err))
			PrintFatalError(err, fmt.Sprintf("Error at processing file. Error writing the error of processing the line %d", count))
		}
		fileInfo := fileInfo{record[0], seq}
		fileInfos = append(fileInfos, fileInfo)
	}
	fmt.Printf("Processed from file: %d\n", len(fileInfos))
	return fileInfos
}

func processResult(destinyFile, errorFile *os.File, fileInfos []fileInfo) {
	w := csv.NewWriter(destinyFile)

	count := 0
	for _, fileInfo := range fileInfos {
		count++
		record := []string{fileInfo.id, strconv.Itoa(fileInfo.seq)}
		if err := w.Write(record); err != nil {
			_, err := errorFile.WriteString(fmt.Sprintf("Error at line %d: %v", count, err))
			PrintFatalError(err, fmt.Sprintf("Error at processing resultfile. Error writing the error of processing the line %d", count))
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		_, err := errorFile.WriteString(fmt.Sprintf("Error with CSV %v", err))
		PrintFatalError(err, fmt.Sprintf("Error at processing resultfile. Error writing the error of processing the line %d", count))
	}

	fmt.Printf("Processed to file: %d\n", count)
}

//PrintFatalError Utility to print an error with a message with a log fatal (it exits the program)
func PrintFatalError(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}

}

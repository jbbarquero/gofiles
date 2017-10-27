package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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

	fileInfos := readFile(originFile, errorFile)
	fileInfoResults := processInfo(destinyFile, errorFile, fileInfos)
	writeResults(destinyFile, errorFile, fileInfoResults)

	fmt.Println("END File processor")

}

//readFile read the file and returns its contents in an array of fileInfo struct
func readFile(originFile, errorFile *os.File) []fileInfo {
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
		fmt.Printf("FileInfo read at line %d: %v\n", count, fileInfo)
		fileInfos = append(fileInfos, fileInfo)
	}
	fmt.Printf("Number of FileInfo processed from file: %d\n", len(fileInfos))
	return fileInfos
}

//processInfo processes the info provided as an array of fileInfo and return it
//in the form of an array of fileInfoResult
func processInfo(destinyFile, errorFile *os.File, fileInfos []fileInfo) []fileInfoResult {

	fileInfoResults := []fileInfoResult{}

	count := 0
	var result *fileInfoResult
	for _, fileInfo := range fileInfos {

		fmt.Printf("Processinf fileInfo %d: %v\n", count, fileInfo)

		if count == 0 {
			result = &fileInfoResult{fileInfo.id, []int{fileInfo.seq}}
			fmt.Printf("FileInfoResult created at line %d: %v\n", count, *result)
		} else {
			if fileInfo.id == fileInfos[count-1].id {
				result.seqs = append(result.seqs, fileInfo.seq)
			} else {
				fileInfoResults = append(fileInfoResults, *result)
				fmt.Printf("FileInfoResults update with a FileInfoResult: %v\n", result)
				result = &fileInfoResult{fileInfo.id, []int{fileInfo.seq}}
				fmt.Printf("FileInfoResult created at line %d: %v\n", count, result)
			}
		}

		count++
	}

	if result != nil {
		fileInfoResults = append(fileInfoResults, *result)
	}

	fmt.Printf("Processed to file: %d\n", count)
	return fileInfoResults
}

//writeResults write the results in the file
func writeResults(destinyFile, errorFile *os.File, fileInfoResults []fileInfoResult) {
	w := csv.NewWriter(destinyFile)
	w.Comma = ';'

	count := 0
	for _, fileInfoResult := range fileInfoResults {
		count++
		textSqs := []string{}
		for i := range fileInfoResult.seqs {
			textSqs = append(textSqs, strconv.Itoa(fileInfoResult.seqs[i]))
		}
		record := []string{fileInfoResult.id, strings.Join(textSqs, ",")}
		fmt.Printf("FileInfoResult written %d: %v\n", count, record)
		if err := w.Write(record); err != nil {
			_, err := errorFile.WriteString(fmt.Sprintf("Error at line %d: %v", count, err))
			PrintFatalError(err, fmt.Sprintf("Error at writting resultfile. Error writing the error of processing the line %d", count))
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		_, err := errorFile.WriteString(fmt.Sprintf("Error with CSV %v", err))
		PrintFatalError(err, fmt.Sprintf("Error at processing resultfile. Error writing the error of processing the line %d", count))
	}

}

//PrintFatalError Utility to print an error with a message with a log fatal (it exits the program)
func PrintFatalError(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}

}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func main() {
	fmt.Println("Go files: searcher begins")

	var args []string
	args = flag.Args()
	log.Printf("flag.Args: %v\n", args)
	args = os.Args
	log.Printf("os.Args: %v\n", args)

	if len(args) == 0 {
		log.Fatalf("Error, missing arguments (%v). Usage: main [options] [path]", len(args))
	}

	path := args[len(args)-1]
	log.Println("Path: ", path)

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.Parse()

	log.Println("verbose: ", verbose)

	fmt.Println("Go files: searcher ends")
}

package processor

import "log"

//PrintFatalError Utility to print an error with a message with a log fatal (it exits the program)
func PrintFatalError(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}

}

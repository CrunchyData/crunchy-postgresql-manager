package mylog

import (
	"log"
	"os"
)

type MyLogger struct {
	Info  *log.Logger
	Error *log.Logger
}

func Init(LOGFILE string) MyLogger {

	file, err := os.OpenFile(LOGFILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", os.Stderr, ":", err)
	}

	mylogger := MyLogger{}

	mylogger.Info = log.New(file,
		"INFO : ",
		log.Ldate|log.Ltime|log.Lshortfile)

	mylogger.Error = log.New(file,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	return mylogger
}

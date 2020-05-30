package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	ERROR = 0
	INFO  = 2
	DEBUG = 4
)

func New(path string) error {

	output, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	log.SetOutput(output)

	return nil
}

func Print(level int, format string, v ...interface{}) {

	switch level {
	case 0:
		log.Println(fmt.Sprintf("[error] "+format, v...))
	case 2:
		log.Println(fmt.Sprintf("[info] "+format, v...))
	case 4:
		log.Println(fmt.Sprintf("[debug] "+format, v...))
	}
}

func Println(level int, v ...interface{}) {

	switch level {
	case 0:
		log.Println("[error]", fmt.Sprintln(v...))
	case 2:
		log.Println("[info]", fmt.Sprintln(v...))
	case 4:
		log.Println("[debug]", fmt.Sprintln(v...))
	}
}

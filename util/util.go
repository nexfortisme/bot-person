package util

import (
	"fmt"
	"log"
)

func HandleErrors(err error) {
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
}

func HandleFatalErrors(err error, message string) {
	if err != nil {
		log.Fatalf(message)
	}
}

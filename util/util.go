package util

import "fmt"

func HandleErrors(err error) {
	if err != nil {
		fmt.Println("Error: " + err.Error());
	}
}